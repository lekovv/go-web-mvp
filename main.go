package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/db"
	"github.com/lekovv/go-web-mvp/http"
	"github.com/lekovv/go-web-mvp/layers"
	"github.com/lekovv/go-web-mvp/middleware"
	"github.com/lekovv/go-web-mvp/scheduler"
)

func main() {
	env, err := config.LoadEnv(".")
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	database := db.ConnectDB(&env)

	appContainer := layers.NewAppContainer(database.DB, &env)

	jobsCtx, jobsCancel := context.WithCancel(context.Background())
	defer jobsCancel()

	jobManager := scheduler.CreateJobs(appContainer.AuthService)
	go jobManager.StartAll(jobsCtx)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			_ = middleware.ErrorHandler()(c)
			return nil
		},
	})

	app.Use(logger.New())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS(&env))
	app.Use(middleware.RateLimiter())
	app.Use(middleware.ErrorHandler())

	http.RegisterRoutes(app, appContainer, &env)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Server starting on port %s", env.ServerPort)
		serverErrors <- app.Listen(":" + env.ServerPort)
	}()

	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v", sig)
		log.Println("Starting graceful shutdown...")

		log.Println("Stopping background jobs...")
		jobsCancel()
		jobManager.Wait()
		log.Println("Background jobs stopped")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		log.Println("Stopping HTTP server...")
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("Error during HTTP server shutdown: %v", err)
		} else {
			log.Println("HTTP server stopped")
		}

		log.Println("Closing database connections...")
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connections: %v", err)
		} else {
			log.Println("Database connections closed")
		}

		log.Println("Shutdown completed")
	}
}
