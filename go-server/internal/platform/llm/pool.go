package llm

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/direwen/go-server/internal/shared/domain"
)

type Rotator struct {
	clients []domain.Client
	mu      sync.Mutex
	index   int
}

type pool struct {
	pool map[domain.LLMTask]*Rotator
	mu   sync.RWMutex
}

func NewClientPool() domain.LLMPool {
	return &pool{
		pool: make(map[domain.LLMTask]*Rotator),
	}
}

func (c *pool) Execute(task domain.LLMTask, cb func(client domain.Client) (any, error)) (any, error) {
	client, size, err := c.getClient(task)
	if err != nil {
		return nil, err
	}

	// Loop through all clients if one fails
	for i := 0; i < size; i++ {
		// Execute the callback function with the current client
		res, err := cb(client)
		if err == nil {
			return res, nil
		}
		// Get next client for retry
		client, _, _ = c.getClient(task)
	}
	return nil, errors.New("all clients exhausted")
}

func (c *pool) Register(task domain.LLMTask, prefix string) {
	// Collect API keys from environment
	apiKeys := collectEnvKeys(prefix)
	if len(apiKeys) == 0 {
		panic(fmt.Sprintf("no API keys found for prefix: %s", prefix))
	}

	// Create clients for each key
	clients := make([]domain.Client, 0, len(apiKeys))
	for _, key := range apiKeys {
		client, err := NewClient(task, key)
		if err != nil {
			panic(fmt.Sprintf("failed to create client for task %s: %v", task, err))
		}
		clients = append(clients, client)
	}

	// Save the rotator to the pool
	c.mu.Lock()
	c.pool[task] = &Rotator{
		clients: clients,
		index:   0,
	}
	c.mu.Unlock()
}

// To find all environment variables matching the prefix
func collectEnvKeys(prefix string) []string {
	var keys []string
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, prefix) {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && parts[1] != "" {
				keys = append(keys, parts[1])
			}
		}
	}
	return keys
}

func (c *pool) getClient(task domain.LLMTask) (domain.Client, int, error) {
	c.mu.RLock()
	rotator, exists := c.pool[task]
	c.mu.RUnlock()

	if !exists {
		return nil, 0, errors.New("task not registered")
	}

	rotator.mu.Lock()
	defer rotator.mu.Unlock()

	if len(rotator.clients) == 0 {
		return nil, 0, errors.New("no clients available")
	}

	client := rotator.clients[rotator.index]
	rotator.index = (rotator.index + 1) % len(rotator.clients)

	return client, len(rotator.clients), nil
}
