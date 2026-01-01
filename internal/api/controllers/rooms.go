package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/Miklakapi/go-file-share/internal/api/dto"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	fileShare "github.com/Miklakapi/go-file-share/internal/file-share/application"
	"github.com/gin-gonic/gin"
)

type RoomsController struct {
	fileShareService *fileShare.Service
}

func NewRoomsController(fileShareService *fileShare.Service) *RoomsController {
	return &RoomsController{fileShareService: fileShareService}
}

func (rC *RoomsController) Get(ctx *gin.Context) {
	rooms, err := rC.fileShareService.Rooms(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := make([]dto.Room, 0, len(rooms))
	for _, r := range rooms {
		result = append(result, dto.NewRoom(r))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (rC *RoomsController) CheckAccess(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

	ok, err := rC.fileShareService.CheckRoomAccess(ctx.Request.Context(), roomId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "room not found",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rC *RoomsController) GetByUUID(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)

	room, ok, err := rC.fileShareService.Room(ctx.Request.Context(), roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "room not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": dto.NewRoom(room),
	})
}

func (rC *RoomsController) Create(ctx *gin.Context) {
	requestData := dto.CreateRoomRequest{}

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data",
		})
		return
	}

	duration := time.Second * time.Duration(requestData.Lifespan)

	room, token, err := rC.fileShareService.CreateRoom(ctx.Request.Context(), requestData.Password, duration)
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

func (rC *RoomsController) Delete(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

	if err := rC.fileShareService.DeleteRoom(ctx.Request.Context(), roomId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cookiePath := strings.TrimSuffix(ctx.Request.URL.Path, "/")
	ctx.SetCookie("auth_token", "", -1, cookiePath, "", false, true)

	ctx.Status(http.StatusNoContent)
}
