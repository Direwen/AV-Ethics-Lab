package main

import (
	"fmt"
	"log"
	"os"

	"github.com/direwen/go-server/internal/config"
	custommw "github.com/direwen/go-server/internal/middleware"
	"github.com/direwen/go-server/internal/scenario"
	"github.com/direwen/go-server/internal/session"
	"github.com/direwen/go-server/internal/template"
	"github.com/direwen/go-server/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lpernett/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found.")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	config.ConnectDB()
	db := config.GetDB()

	// Session
	sessionRepo := session.NewRepository(db)
	sessionService := session.NewService(sessionRepo)
	sessionHandler := session.NewHandler(sessionService)

	// Scenario
	scenarioRepo := scenario.NewRepository(db)
	scenarioService := scenario.NewService(scenarioRepo)
	scenarioHandler := scenario.NewHandler(scenarioService)

	// Template
	templateRepo := template.NewRepository(db)

	// Seed Templates
	if err := template.SeedContextTemplates(templateRepo); err != nil {
		log.Fatal("Failed to seed templates: ", err)
	}

	// Init Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = util.CustomEchoErrorHandler

	e.POST("api/v1/sessions", sessionHandler.Create)

	protected := e.Group("/api/v1")
	protected.Use(custommw.JWTMiddleware())
	{
		protected.GET("/scenarios/next", scenarioHandler.GetNext)
	}

	if os.Getenv("LOCAL_FRONTEND_PORT") == "" {
		log.Fatal("LOCAL_FRONTEND_PORT is not set")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			fmt.Sprintf("http://localhost:%s", os.Getenv("LOCAL_FRONTEND_PORT")),
		},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	if os.Getenv("SERVER_PORT") == "" {
		log.Fatal("SERVER_PORT is not set")
	}

	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
