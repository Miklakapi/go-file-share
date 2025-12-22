package handlers

import (
	"net/http"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Deps *app.DependencyBag
}

func NewAuthHandler(deps *app.DependencyBag) *AuthHandler {
	return &AuthHandler{Deps: deps}
}

func (h *AuthHandler) Auth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
