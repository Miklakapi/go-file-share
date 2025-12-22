package handlers

import (
	"net/http"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type RoomsHandler struct {
	Deps *app.DependencyBag
}

func NewRoomsHandler(deps *app.DependencyBag) *RoomsHandler {
	return &RoomsHandler{Deps: deps}
}

func (h *RoomsHandler) Get(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *RoomsHandler) GetByUUID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *RoomsHandler) Create(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *RoomsHandler) Delete(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
