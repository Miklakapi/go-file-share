package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Miklakapi/go-file-share/internal/api"
	"github.com/Miklakapi/go-file-share/internal/api/handlers"
	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/Miklakapi/go-file-share/internal/config"
	"github.com/Miklakapi/go-file-share/internal/files"
	"github.com/gin-gonic/gin"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fileHub := files.NewFileHub()
	go fileHub.Run(appCtx)

	deps := app.NewDependencyBag(config, fileHub, appCtx)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	api.RegisterRoutes(engine, deps, &api.Handlers{
		HealthHandler: handlers.NewHealthHandler(deps),
		PagesHandler:  handlers.NewPagesHandler(deps),
	})

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: engine,
	}

	go func() {
		log.Println("HTTP server started on :" + config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	<-appCtx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped gracefully")
}
