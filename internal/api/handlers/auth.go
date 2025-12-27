package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/Miklakapi/go-file-share/internal/api/dto"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Deps *app.DependencyBag
}

func NewAuthHandler(deps *app.DependencyBag) *AuthHandler {
	return &AuthHandler{Deps: deps}
}

func (h *AuthHandler) Auth(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

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
	maxAge := int(time.Until(expiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie("auth_token", token, maxAge, cookiePath, "", false, true)

	ctx.Status(http.StatusNoContent)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.Deps.FileShareService.LogoutRoom(h.Deps.AppContext, roomId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
