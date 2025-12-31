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
	"github.com/Miklakapi/go-file-share/internal/api/controllers"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/Miklakapi/go-file-share/internal/config"
	filestore "github.com/Miklakapi/go-file-share/internal/file-share/adapters/file-store"
	roomrepository "github.com/Miklakapi/go-file-share/internal/file-share/adapters/room-repository"
	"github.com/Miklakapi/go-file-share/internal/file-share/adapters/security"
	"github.com/gin-gonic/gin"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	roomRepo := roomrepository.NewMemoryRepo()
	fileStore := filestore.DiskStore{}
	hasher := security.BcryptHasher{Cost: 12}
	tokenService := security.NewJwtService(config.JWTSecret)

	deps := app.NewDependencyBag(appCtx, config, roomRepo, fileStore, hasher, tokenService)

	if err := fileStore.ClearAll(appCtx, config.UploadDir); err != nil {
		log.Fatalf("file error: %v", err)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	api.RegisterRoutes(engine, &api.ControllerBag{
		HealthController: controllers.NewHealthController(deps),
		HtmlController:   controllers.NewHtmlController(deps),
		AuthController:   controllers.NewAuthController(deps),
		RoomsController:  controllers.NewRoomsController(deps),
		FilesController:  controllers.NewFilesController(deps),
		AuthMiddleware:   middleware.AuthMiddleware(deps),
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
