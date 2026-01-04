package controllers

import (
	"net/http"
	"strings"
	"time"

	apierrors "github.com/Miklakapi/go-file-share/internal/api/api-errors"
	"github.com/Miklakapi/go-file-share/internal/api/dto"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	fileShare "github.com/Miklakapi/go-file-share/internal/file-share/application"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	fileShareService *fileShare.Service
}

func NewAuthController(fileShareService *fileShare.Service) *AuthController {
	return &AuthController{fileShareService: fileShareService}
}

func (aC *AuthController) Auth(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)

	requestData := dto.AuthRoomRequest{}
	if err := ctx.ShouldBind(&requestData); err != nil {
		_ = ctx.Error(apierrors.ErrInvalidRequest)
		return
	}

	token, expiresAt, err := aC.fileShareService.AuthRoom(ctx.Request.Context(), roomId, requestData.Password, time.Second*time.Duration(requestData.Lifespan))
	if err != nil {
		_ = ctx.Error(err)
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

func (aC *AuthController) Logout(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

	if err := aC.fileShareService.LogoutRoom(ctx.Request.Context(), roomId, token); err != nil {
		_ = ctx.Error(err)
		return
	}

	cookiePath := strings.TrimSuffix(strings.TrimSuffix(ctx.Request.URL.Path, "/"), "/logout")
	ctx.SetCookie("auth_token", "", -1, cookiePath, "", false, true)

	ctx.Status(http.StatusNoContent)
}
