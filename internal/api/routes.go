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
	router.NoRoute(handlers.PagesHandler.SPAFallback)

	api := router.Group("/api/v1")
	{
		api.GET("/ping", handlers.HealthHandler.Ping)
		api.GET("/health", handlers.HealthHandler.Health)

		rooms := api.Group("/rooms")
		{
			rooms.GET("", handlers.RoomsHandler.Get)
			rooms.POST("", handlers.RoomsHandler.Create)

			room := rooms.Group("/:roomID", middleware.RoomIDParam())
			{
				room.GET("", handlers.RoomsHandler.GetByUUID)
				room.POST("/auth", handlers.AuthHandler.Auth)
			}
		}

		securedRooms := api.Group("/rooms/:roomID", middleware.RoomIDParam(), middleware.AuthMiddleware(deps))
		{
			securedRooms.DELETE("", handlers.RoomsHandler.Delete)
			securedRooms.GET("/access", handlers.RoomsHandler.CheckAccess)
			securedRooms.POST("/logout", handlers.AuthHandler.Logout)

			files := securedRooms.Group("/files")
			{
				files.GET("", handlers.FilesHandler.Get)
				files.POST("", handlers.FilesHandler.Upload)

				file := files.Group("/:fileID", middleware.FileIDParam())
				{
					file.GET("", handlers.FilesHandler.GetByUUID)
					file.GET("/download", handlers.FilesHandler.Download)
					file.DELETE("", handlers.FilesHandler.Delete)
				}
			}
		}
	}
}
