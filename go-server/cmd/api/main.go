package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/direwen/go-server/internal/config"
	"github.com/direwen/go-server/internal/dashboard"
	custommw "github.com/direwen/go-server/internal/middleware"
	"github.com/direwen/go-server/internal/platform/llm"
	"github.com/direwen/go-server/internal/response"
	"github.com/direwen/go-server/internal/scenario"
	"github.com/direwen/go-server/internal/session"
	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/direwen/go-server/internal/template"
	"github.com/direwen/go-server/internal/util"
	"github.com/direwen/go-server/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lpernett/godotenv"
	"strconv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found.")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	if os.Getenv("EXPERIMENT_TARGET_COUNT") == "" {
		log.Fatal("Required to set EXPERIMENT_TARGET_COUNT")
	}

	config.ConnectDB()
	db := config.GetDB()
	txManager := database.NewTransactionManager(db)

	// Init LLM Client Pool
	pool := llm.NewClientPool()
	pool.Register(domain.TaskScenario, "GROQ_API_KEY")
	pool.Register(domain.TaskFeedback, "OPENROUTER_API_KEY")

	// Template
	templateRepo := template.NewRepository(db)
	templateService := template.NewService(templateRepo)

	// Seed Templates (must happen before loading into cache)
	if err := template.SeedContextTemplates(templateRepo); err != nil {
		log.Fatal("Failed to seed templates: ", err)
	}

	log.Println("Loading Map Templates into Memory ....")
	if err := templateService.LoadAllTemplates(context.Background()); err != nil {
		log.Fatal("Failed to load templates: ", err)
	}
	log.Println("Templates Loaded")


	// Session
	experimentTargetCount, err := strconv.Atoi(os.Getenv("EXPERIMENT_TARGET_COUNT"))
	if err != nil {
		log.Fatal("Failed to convert EXPERIMENT_TARGET_COUNT to int: ", err)
	}
	sessionRepo := session.NewRepository(db)
	sessionService := session.NewService(sessionRepo, pool, experimentTargetCount)
	sessionHandler := session.NewHandler(sessionService)

	// Scenario
	scenarioRepo := scenario.NewRepository(db)
	scenarioService := scenario.NewService(
		scenarioRepo,
		sessionService,
		templateService,
		pool,
	)
	scenarioHandler := scenario.NewHandler(scenarioService)

	// Response
	responseRepo := response.NewRepository(db)
	responseService := response.NewService(responseRepo, sessionService, scenarioService, txManager)
	responseHandler := response.NewHandler(responseService)

	dashboardRepo := dashboard.NewRepository(db)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	// Init Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = util.CustomEchoErrorHandler

	e.POST("api/v1/sessions", sessionHandler.Create)
	e.GET("api/v1/dashboard", dashboardHandler.GetDashboard)

	protected := e.Group("/api/v1")
	protected.Use(custommw.JWTMiddleware())
	{
		protected.GET("/sessions/feedback", sessionHandler.GetSessionFeedback)
		protected.GET("/scenarios/next", scenarioHandler.GetNext)
		protected.POST("/scenarios/:scenario_id/responses", responseHandler.Create)
	}

	if os.Getenv("LOCAL_FRONTEND_PORT") == "" {
		log.Fatal("LOCAL_FRONTEND_PORT is not set")
	}

	origins := []string{
		fmt.Sprintf("http://localhost:%s", os.Getenv("LOCAL_FRONTEND_PORT")),
	}
	if clientURL := os.Getenv("CLIENT_URL"); clientURL != "" {
		origins = append(origins, clientURL)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	if os.Getenv("SERVER_PORT") == "" {
		log.Fatal("SERVER_PORT is not set")
	}

	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
