package api

import (
	"github.com/Miklakapi/go-file-share/internal/api/controllers"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type ControllerBag struct {
	HealthController *controllers.HealthController
	HtmlController   *controllers.HtmlController
	AuthController   *controllers.AuthController
	RoomsController  *controllers.RoomsController
	FilesController  *controllers.FilesController
	SSEController    *controllers.SSEController
	DirectController *controllers.DirectController
	AuthMiddleware   gin.HandlerFunc
	ErrorMiddleware  gin.HandlerFunc
}

func RegisterRoutes(router *gin.Engine, cB *ControllerBag) {
	securedRouter := router.Group("", middleware.SecureHeaders)

	securedRouter.GET("/", cB.HtmlController.Index)
	securedRouter.Static("/assets", cB.HtmlController.Assets())
	securedRouter.StaticFile("/favicon.ico", cB.HtmlController.Favicon())
	router.NoRoute(middleware.SecureHeaders, cB.HtmlController.SPAFallback)

	api := securedRouter.Group("/api/v1", cB.ErrorMiddleware)
	api.GET("/ping", cB.HealthController.Ping)
	api.GET("/health", cB.HealthController.Health)
	api.GET("/sse", cB.SSEController.SSE)

	direct := api.Group("/direct/:code")
	direct.GET("/download", cB.DirectController.DownloadStream)
	direct.POST("/upload", cB.DirectController.UploadStream)

	rooms := api.Group("/rooms")
	rooms.GET("", cB.RoomsController.Get)
	rooms.POST("", cB.RoomsController.Create)

	room := rooms.Group("/:roomID", middleware.SetRoomIDParam())
	room.GET("", cB.RoomsController.GetByUUID)
	room.POST("/auth", cB.AuthController.Auth)

	securedRooms := api.Group("/rooms/:roomID", middleware.SetRoomIDParam(), cB.AuthMiddleware)
	securedRooms.DELETE("", cB.RoomsController.Delete)
	securedRooms.GET("/access", cB.RoomsController.CheckAccess)
	securedRooms.POST("/logout", cB.AuthController.Logout)

	files := securedRooms.Group("/files")
	files.GET("", cB.FilesController.Get)
	files.POST("", cB.FilesController.Upload)

	file := files.Group("/:fileID", middleware.SetFileIDParam())
	file.GET("", cB.FilesController.GetByUUID)
	file.GET("/download", cB.FilesController.Download)
	file.DELETE("", cB.FilesController.Delete)
}
