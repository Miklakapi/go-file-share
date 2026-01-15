package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Miklakapi/go-file-share/internal/api"
	"github.com/Miklakapi/go-file-share/internal/api/controllers"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/Miklakapi/go-file-share/internal/config"
	"github.com/Miklakapi/go-file-share/internal/file-share/adapters/db"
	directtransfer "github.com/Miklakapi/go-file-share/internal/file-share/adapters/direct-transfer"
	eventbus "github.com/Miklakapi/go-file-share/internal/file-share/adapters/event-bus"
	filestore "github.com/Miklakapi/go-file-share/internal/file-share/adapters/file-store"
	sqliterepository "github.com/Miklakapi/go-file-share/internal/file-share/adapters/room-repository/sqlite-repository"
	"github.com/Miklakapi/go-file-share/internal/file-share/adapters/security"
	fileShare "github.com/Miklakapi/go-file-share/internal/file-share/application"
	fileShareDomain "github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/jobs"
	"github.com/gin-gonic/gin"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("cannot open log file: %v", err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gin.DefaultWriter = mw
	gin.DefaultErrorWriter = mw

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sqliteDb, err := db.NewSqlite(config.SqlitePath)
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDb.Conn.Close()

	err = sqliterepository.MakeMigrations(appCtx, sqliteDb.Conn, "./sqlite-migrations")
	if err != nil {
		log.Fatal(err)
	}

	roomRepo := sqliterepository.New(appCtx, sqliteDb.Conn)
	eventBus := eventbus.New()
	directTransfer := directtransfer.New()
	fileStore := filestore.DiskStore{}
	hasher := security.BcryptHasher{Cost: 12}
	tokenService := security.NewJwtService(config.JWTSecret)
	fileShareSettings := fileShareDomain.NewPolicy(
		config.DefaultRoomTTL,
		config.TokenTTL,
		config.MaxFiles,
		config.MaxRoomBytes,
		config.MaxRoomLifespan,
		config.MaxTokenLifespan,
		config.UploadDir,
	)
	fileShareService := fileShare.NewService(roomRepo, fileStore, hasher, tokenService, fileShareSettings)
	roomCleanupJob := jobs.New(fileShareService, eventBus, config.CleanupInterval)

	if err := roomRepo.WipeAll(appCtx); err != nil {
		log.Fatalf("file error: %v", err)
	}
	if err := fileStore.ClearAll(appCtx, config.UploadDir); err != nil {
		log.Fatalf("file error: %v", err)
	}

	closeJob, err := roomCleanupJob.Run(appCtx)
	if err != nil {
		log.Fatalf("file error: %v", err)
	}
	defer closeJob()

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	api.RegisterRoutes(engine, &api.ControllerBag{
		HealthController: controllers.NewHealthController(),
		HtmlController:   controllers.NewHtmlController(config.PublicDir),
		AuthController:   controllers.NewAuthController(fileShareService),
		RoomsController:  controllers.NewRoomsController(fileShareService, eventBus),
		FilesController:  controllers.NewFilesController(fileShareService),
		SSEController:    controllers.NewSSEController(appCtx, eventBus),
		DirectController: controllers.NewDirectController(directTransfer),
		AuthMiddleware:   middleware.AuthMiddleware(tokenService),
		ErrorMiddleware:  middleware.ErrorMiddleware(),
	})

	srv := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           engine,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
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
