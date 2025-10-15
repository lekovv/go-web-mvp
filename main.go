package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobManager := scheduler.CreateJobs(appContainer.AuthService)
	go jobManager.StartAll(ctx)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS(&env))
	app.Use(middleware.RateLimiter())

	http.RegisterRoutes(app, appContainer, &env)

	gracefulShutdown(app, cancel, jobManager)

	log.Fatal(app.Listen(":" + env.ServerPort))
}

func gracefulShutdown(
	app *fiber.App,
	cancel context.CancelFunc,
	jobManager *scheduler.JobManager,
) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		cancel()

		log.Println("Waiting for jobs to finish...")
		jobManager.Wait()

		log.Println("Shutting down HTTP server...")
		err := app.Shutdown()
		if err != nil {
			return
		}
	}()
}
