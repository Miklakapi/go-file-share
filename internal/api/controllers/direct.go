package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	apierrors "github.com/Miklakapi/go-file-share/internal/api/api-errors"
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
		_ = ctx.Error(apierrors.ErrInvalidRequest)
		return
	}

	transfer, err := dC.directTransfer.Receive(ctx.Request.Context(), code)
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
		_ = ctx.Error(apierrors.ErrInvalidRequest)
		return
	}

	fh, err := ctx.FormFile("file")
	if err != nil {
		_ = ctx.Error(apierrors.ErrInvalidFile)
		return
	}

	src, err := fh.Open()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	defer func() { _ = src.Close() }()

	if err := dC.directTransfer.Send(ctx.Request.Context(), code, fh.Filename, src); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
