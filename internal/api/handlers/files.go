package handlers

import (
	"io"
	"net/http"
	"strconv"

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

	tokenAny, ok := ctx.Get(middleware.CtxTokenKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Missing auth token",
		})
		return
	}
	token := tokenAny.(string)

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

	tokenAny, ok := ctx.Get(middleware.CtxTokenKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Missing auth token",
		})
		return
	}
	token := tokenAny.(string)

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

	tokenAny, ok := ctx.Get(middleware.CtxTokenKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Missing auth token",
		})
		return
	}
	token := tokenAny.(string)

	meta, rc, err := h.Deps.FileShareService.DownloadFile(h.Deps.AppContext, roomId, fileId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer func() { _ = rc.Close() }()

	ctx.Header("Content-Disposition", `attachment; filename="`+meta.Name+`"`)
	ctx.Header("Content-Type", "application/octet-stream")
	if meta.Size > 0 {
		ctx.Header("Content-Length", strconv.FormatInt(meta.Size, 10))
	}

	_, copyErr := io.Copy(ctx.Writer, rc)
	if copyErr != nil {
		return
	}
}

func (h *FilesHandler) Upload(ctx *gin.Context) {
	roomIdAny, _ := ctx.Get(middleware.CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)

	tokenAny, ok := ctx.Get(middleware.CtxTokenKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Missing auth token",
		})
		return
	}
	token := tokenAny.(string)

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

	file, err := h.Deps.FileShareService.UploadFile(h.Deps.AppContext, roomId, token, fh.Filename, src)
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

	tokenAny, ok := ctx.Get(middleware.CtxTokenKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Missing auth token",
		})
		return
	}
	token := tokenAny.(string)

	if err := h.Deps.FileShareService.DeleteFile(h.Deps.AppContext, roomId, fileId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
