package handlers

import (
	"net/http"

	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FilesHandler struct {
	Deps *app.DependencyBag
}

func NewFilesHandler(deps *app.DependencyBag) *FilesHandler {
	return &FilesHandler{Deps: deps}
}

func (h *FilesHandler) Get(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	files, err := h.Deps.FileShareService.Files(h.Deps.AppContext, roomId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": files,
	})
}

func (h *FilesHandler) GetByUUID(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	fileIdAny, _ := ctx.Get(middleware.CtxFileIDKey)
	fileId := fileIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	file, err := h.Deps.FileShareService.File(h.Deps.AppContext, roomId, fileId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": file,
	})
}

func (h *FilesHandler) Download(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	fileIdAny, _ := ctx.Get(middleware.CtxFileIDKey)
	fileId := fileIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err = h.Deps.FileShareService.DownloadFile(h.Deps.AppContext, roomId, fileId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *FilesHandler) Upload(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	fh, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	src, err := fh.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer func() { _ = src.Close() }()

	file, err := h.Deps.FileShareService.UploadFile(h.Deps.AppContext, roomId, token, "", src)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": file,
	})
}

func (h *FilesHandler) Delete(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	fileIdAny, _ := ctx.Get(middleware.CtxFileIDKey)
	fileId := fileIdAny.(uuid.UUID)

	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := h.Deps.FileShareService.DeleteFile(h.Deps.AppContext, roomId, fileId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
