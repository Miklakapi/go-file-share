package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/gin-gonic/gin"
)

type DirectController struct {
	directTransfer ports.DirectTransfer
}

func NewDirectController(directTransfer ports.DirectTransfer) *DirectController {
	return &DirectController{
		directTransfer: directTransfer,
	}
}

func (dC *DirectController) DownloadStream(ctx *gin.Context) {
	code := strings.TrimSpace(ctx.Param("code"))
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}

	transfer, err := dC.directTransfer.Receive(ctx, code)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	defer dC.directTransfer.Cancel(code)

	ctx.Header("Cache-Control", "no-store")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, transfer.Filename))

	_, _ = io.Copy(ctx.Writer, transfer.Reader)
}

func (dC *DirectController) UploadStream(ctx *gin.Context) {
	code := strings.TrimSpace(ctx.Param("code"))
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}

	fh, err := ctx.FormFile("file")
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	src, err := fh.Open()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	defer func() { _ = src.Close() }()

	if err := dC.directTransfer.Send(code, fh.Filename, src); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "TODO",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
