package controllers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/Miklakapi/go-file-share/internal/api/dto"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type FilesController struct {
	Deps *app.DependencyBag
}

func NewFilesController(deps *app.DependencyBag) *FilesController {
	return &FilesController{Deps: deps}
}

func (h *FilesController) Get(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

	files, err := h.Deps.FileShareService.Files(h.Deps.AppContext, roomId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := make([]dto.RoomFile, 0, len(files))
	for _, f := range files {
		result = append(result, dto.NewFileRoomFile(f))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *FilesController) GetByUUID(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	fileId := middleware.MustFileIDParam(ctx)
	token := middleware.MustToken(ctx)

	file, err := h.Deps.FileShareService.File(h.Deps.AppContext, roomId, fileId, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": dto.NewFileRoomFile(file),
	})
}

func (h *FilesController) Download(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	fileId := middleware.MustFileIDParam(ctx)
	token := middleware.MustToken(ctx)

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

func (h *FilesController) Upload(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": dto.NewFileRoomFile(file),
	})
}

func (h *FilesController) Delete(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	fileId := middleware.MustFileIDParam(ctx)
	token := middleware.MustToken(ctx)

	if err := h.Deps.FileShareService.DeleteFile(h.Deps.AppContext, roomId, fileId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
