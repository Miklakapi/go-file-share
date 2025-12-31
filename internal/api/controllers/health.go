package controllers

import (
	"net/http"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	Deps *app.DependencyBag
}

func NewHealthController(deps *app.DependencyBag) *HealthController {
	return &HealthController{Deps: deps}
}

func (h *HealthController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *HealthController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
