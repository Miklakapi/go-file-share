package controllers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/Miklakapi/go-file-share/internal/api/dto"
	"github.com/Miklakapi/go-file-share/internal/api/middleware"
	fileShare "github.com/Miklakapi/go-file-share/internal/file-share/application"
	"github.com/gin-gonic/gin"
)

type FilesController struct {
	fileShareService *fileShare.Service
}

func NewFilesController(fileShareService *fileShare.Service) *FilesController {
	return &FilesController{fileShareService: fileShareService}
}

func (fC *FilesController) Get(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	token := middleware.MustToken(ctx)

	files, err := fC.fileShareService.Files(ctx.Request.Context(), roomId, token)
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

func (fC *FilesController) GetByUUID(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	fileId := middleware.MustFileIDParam(ctx)
	token := middleware.MustToken(ctx)

	file, err := fC.fileShareService.File(ctx.Request.Context(), roomId, fileId, token)
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

func (fC *FilesController) Download(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	fileId := middleware.MustFileIDParam(ctx)
	token := middleware.MustToken(ctx)

	meta, rc, err := fC.fileShareService.DownloadFile(ctx.Request.Context(), roomId, fileId, token)
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

func (fC *FilesController) Upload(ctx *gin.Context) {
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

	file, err := fC.fileShareService.UploadFile(ctx.Request.Context(), roomId, token, fh.Filename, src)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": dto.NewFileRoomFile(file),
	})
}

func (fC *FilesController) Delete(ctx *gin.Context) {
	roomId := middleware.MustRoomIDParam(ctx)
	fileId := middleware.MustFileIDParam(ctx)
	token := middleware.MustToken(ctx)

	if err := fC.fileShareService.DeleteFile(ctx.Request.Context(), roomId, fileId, token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
