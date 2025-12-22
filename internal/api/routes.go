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
}

func RegisterRoutes(router *gin.Engine, deps *app.DependencyBag, handlers *Handlers) {
	router.GET("/", handlers.PagesHandler.Index)
	router.Static("/assets", deps.Config.PublicDir+"/assets")
	router.StaticFile("/favicon.ico", deps.Config.PublicDir+"/favicon.ico")

	apiRoutes := router.Group("/api/v1")

	apiRoutes.GET("/ping", handlers.HealthHandler.Ping)
	apiRoutes.GET("/health", handlers.HealthHandler.Health)

	apiRoutes.Use(middleware.AuthMiddleware(deps))
	// Secured endpoints

	// SPA Fallback
	router.NoRoute(handlers.PagesHandler.SPAFallback)
}
