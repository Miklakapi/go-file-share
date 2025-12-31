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

type RoomsController struct {
	Deps *app.DependencyBag
}

func NewRoomsController(deps *app.DependencyBag) *RoomsController {
	return &RoomsController{Deps: deps}
}

func (h *RoomsController) Get(ctx *gin.Context) {
	rooms, err := h.Deps.FileShareService.Rooms(h.Deps.AppContext)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	result := make([]dto.Room, 0, len(rooms))
	for _, r := range rooms {
		result = append(result, dto.NewRoom(r))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *RoomsController) CheckAccess(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

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

func (h *RoomsController) GetByUUID(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)

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
		"data": dto.NewRoom(room),
	})
}

func (h *RoomsController) Create(ctx *gin.Context) {
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
		"data":  dto.NewRoom(room),
		"token": token,
	})
}

func (h *RoomsController) Delete(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

	if err := h.Deps.FileShareService.DeleteRoom(h.Deps.AppContext, roomId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cookiePath := strings.TrimSuffix(ctx.Request.URL.Path, "/")
	ctx.SetCookie("auth_token", "", -1, cookiePath, "", false, true)

	ctx.Status(http.StatusNoContent)
}
