package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transfer struct {
	r           io.ReadCloser
	filename    string
	contentType string
}

type DirectController struct {
	testChan chan *transfer
}

func NewDirectController() *DirectController {
	return &DirectController{
		testChan: make(chan *transfer, 1),
	}
}

func (dC *DirectController) DownloadStream(ctx *gin.Context) {
	select {
	case tr := <-dC.testChan:
		defer func() { _ = tr.r.Close() }()

		ct := tr.contentType
		if ct == "" {
			ct = "application/octet-stream"
		}

		ctx.Header("Cache-Control", "no-store")
		ctx.Header("Content-Type", ct)
		ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, tr.filename))

		_, _ = io.Copy(ctx.Writer, tr.r)
		return

	case <-ctx.Request.Context().Done():
		return
	}
}

func (dC *DirectController) UploadStream(ctx *gin.Context) {
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

	pr, pw := io.Pipe()

	select {
	case dC.testChan <- &transfer{
		r:           pr,
		filename:    fh.Filename,
		contentType: fh.Header.Get("Content-Type"),
	}:

	case <-ctx.Request.Context().Done():
		_ = pr.Close()
		_ = pw.Close()
		return
	}

	_, copyErr := io.Copy(pw, src)
	_ = src.Close()

	if copyErr != nil {
		_ = pw.CloseWithError(copyErr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Upload interrupted",
		})
		return
	}
	_ = pw.Close()

	ctx.Status(http.StatusNoContent)
}
