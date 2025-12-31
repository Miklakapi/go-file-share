package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/Miklakapi/go-file-share/internal/api/dto"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Deps *app.DependencyBag
}

func NewAuthController(deps *app.DependencyBag) *AuthController {
	return &AuthController{Deps: deps}
}

func (h *AuthController) Auth(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)

	requestData := dto.AuthRoomRequest{}
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data",
		})
		return
	}

	token, expiresAt, err := h.Deps.FileShareService.AuthRoom(h.Deps.AppContext, roomId, requestData.Password, time.Second*time.Duration(requestData.Lifespan))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	cookiePath := strings.TrimSuffix(strings.TrimSuffix(ctx.Request.URL.Path, "/"), "/auth")
	maxAge := max(int(time.Until(expiresAt).Seconds()), 0)

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie("auth_token", token, maxAge, cookiePath, "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"data": token,
	})
}

func (h *AuthController) Logout(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

	if err := h.Deps.FileShareService.LogoutRoom(h.Deps.AppContext, roomId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cookiePath := strings.TrimSuffix(strings.TrimSuffix(ctx.Request.URL.Path, "/"), "/logout")
	ctx.SetCookie("auth_token", "", -1, cookiePath, "", false, true)

	ctx.Status(http.StatusNoContent)
}
