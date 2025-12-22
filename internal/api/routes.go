package api

import (
	"github.com/Miklakapi/go-file-share/internal/api/handlers"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	HealthHandler *handlers.HealthHandler
	PagesHandler  *handlers.PagesHandler
	RoomsHandler  *handlers.RoomsHandler
	AuthHandler   *handlers.AuthHandler
	FilesHandler  *handlers.FilesHandler
}

func RegisterRoutes(router *gin.Engine, deps *app.DependencyBag, handlers *Handlers) {
	router.GET("/", handlers.PagesHandler.Index)
	router.Static("/assets", deps.Config.PublicDir+"/assets")
	router.StaticFile("/favicon.ico", deps.Config.PublicDir+"/favicon.ico")

	api := router.Group("/api/v1")
	api.GET("/ping", handlers.HealthHandler.Ping)
	api.GET("/health", handlers.HealthHandler.Health)

	apiRooms := api.Group("/rooms")
	apiRooms.GET("", handlers.RoomsHandler.Get)
	apiRooms.GET("/:uuid", handlers.RoomsHandler.GetByUUID)
	apiRooms.POST("", handlers.RoomsHandler.Create)
	apiRooms.POST("/:uuid/auth", handlers.AuthHandler.Auth)

	secured := api.Group("", middleware.AuthMiddleware(deps))

	securedRooms := secured.Group("/rooms/:uuid")
	securedRooms.DELETE("", handlers.RoomsHandler.Delete)
	securedRooms.POST("/logout", handlers.AuthHandler.Logout)
	securedRooms.GET("/files", handlers.FilesHandler.Get)
	securedRooms.GET("/files/:fUuid", handlers.FilesHandler.GetByUUID)
	securedRooms.GET("/files/:fUuid/download", handlers.FilesHandler.Download)
	securedRooms.POST("/files", handlers.FilesHandler.Upload)
	securedRooms.DELETE("/files/:fUuid", handlers.FilesHandler.Delete)

	// SPA Fallback
	router.NoRoute(handlers.PagesHandler.SPAFallback)
}
