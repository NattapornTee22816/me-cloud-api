package main

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/NattapornTee22816/me-cloud-api/pkg/cache"
	"github.com/NattapornTee22816/me-cloud-api/pkg/db"
	"github.com/NattapornTee22816/me-cloud-api/pkg/handler"
	"github.com/NattapornTee22816/me-cloud-api/pkg/middleware"
	"github.com/NattapornTee22816/me-cloud-api/pkg/repository"
	"github.com/NattapornTee22816/me-cloud-api/pkg/server"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, database, cacheSession, cacheResponse := initial()

	go func() {
		app.Run()
	}()

	app.Hooks().OnShutdown(func() error {
		fmt.Println("Running cleanup tasks...")

		_ = database.Close()
		_ = cacheSession.Close()
		_ = cacheResponse.Close()

		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// This blocks the main thread until an interrupt is received
	_ = <-c
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()
}

func initial() (*server.Server, clickhouse.Conn, *redis.Client, *redis.Client) {
	database, err := db.NewDatabase()
	if err != nil {
		panic(err)
	}

	cacheSession, err := cache.NewCache(cache.DBSession)
	if err != nil {
		panic(err)
	}

	cacheResponse, err := cache.NewCache(cache.DBCache)
	if err != nil {
		panic(err)
	}

	app := server.NewServer()

	// set middleware
	middleware.NewFiberMiddleware(app.App, database, cacheSession, cacheResponse).
		UseRecovery().
		UseCompress().
		UseHelmet().
		UseCors().
		UseRequestID().
		UseKeyAuth().
		UseLogging().
		UseHealthCheck().
		UseSwagger()

	setHandler(app.App, &database, cacheSession, cacheResponse)

	return app, database, cacheSession, cacheResponse
}

func setHandler(app *fiber.App, database *clickhouse.Conn, cacheSession, cacheResponse *redis.Client) {
	repo := repository.NewRepository(database)

	handler.NewUserHandler(app, repo, cacheSession)
	handler.NewAuthHandler(app, repo, cacheSession)

}
