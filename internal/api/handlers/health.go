package handlers

import (
	"net/http"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Deps *app.DependencyBag
}

func NewHealthHandler(deps *app.DependencyBag) *HealthHandler {
	return &HealthHandler{Deps: deps}
}

func (h *HealthHandler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *HealthHandler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
