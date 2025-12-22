package handlers

import (
	"net/http"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type FilesHandler struct {
	Deps *app.DependencyBag
}

func NewFilesHandler(deps *app.DependencyBag) *FilesHandler {
	return &FilesHandler{Deps: deps}
}

func (h *FilesHandler) Get(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *FilesHandler) GetByUUID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *FilesHandler) Download(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *FilesHandler) Upload(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *FilesHandler) Delete(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
