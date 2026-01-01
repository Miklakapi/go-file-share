package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SSEController struct {
}

func NewSSEController() *SSEController {
	return &SSEController{}
}

func (sC *SSEController) SSE(ctx *gin.Context) {
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.String(http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	messageTicker := time.NewTicker(15 * time.Second)
	defer messageTicker.Stop()

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case t := <-messageTicker.C:
			fmt.Fprintf(ctx.Writer, "event: time\n")
			fmt.Fprintf(ctx.Writer, "data: Current time %s\n\n", t.Format(time.RFC3339))
			flusher.Flush()

		case <-pingTicker.C:
			fmt.Fprintf(ctx.Writer, "data: Ping\n\n")
			flusher.Flush()

		case <-ctx.Request.Context().Done():
			return
		}
	}
}
