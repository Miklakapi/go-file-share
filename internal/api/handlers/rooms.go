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

type RoomsHandler struct {
	Deps *app.DependencyBag
}

func NewRoomsHandler(deps *app.DependencyBag) *RoomsHandler {
	return &RoomsHandler{Deps: deps}
}

func (h *RoomsHandler) Get(ctx *gin.Context) {
	rooms, err := h.Deps.FileShareService.Rooms(h.Deps.AppContext)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

func (h *RoomsHandler) CheckAccess(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ok, err := h.Deps.FileShareService.CheckRoomAccess(h.Deps.AppContext, roomId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "room not found",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (h *RoomsHandler) GetByUUID(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	room, ok, err := h.Deps.FileShareService.Room(h.Deps.AppContext, roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "room not found",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

func (h *RoomsHandler) Create(ctx *gin.Context) {
	requestData := dto.CreateRoomRequest{}

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data",
		})
		return
	}

	duration := time.Second * time.Duration(requestData.Lifespan)

	room, token, err := h.Deps.FileShareService.CreateRoom(h.Deps.AppContext, requestData.Password, duration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	basePath := strings.TrimSuffix(ctx.Request.URL.Path, "/")
	cookiePath := basePath + "/" + room.ID.String()

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie("auth_token", token, int(duration.Seconds()), cookiePath, "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

func (h *RoomsHandler) Delete(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.Deps.FileShareService.DeleteRoom(h.Deps.AppContext, roomId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cookiePath := strings.TrimSuffix(ctx.Request.URL.Path, "/")
	ctx.SetCookie("auth_token", "", -1, cookiePath, "", false, true)

	ctx.Status(http.StatusNoContent)
}
