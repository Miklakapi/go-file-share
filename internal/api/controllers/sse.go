package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/gin-gonic/gin"
)

type SSEController struct {
	appCtx          context.Context
	eventSubscriber ports.EventSubscriber
}

func NewSSEController(appCtx context.Context, eventSubscriber ports.EventSubscriber) *SSEController {
	return &SSEController{appCtx: appCtx, eventSubscriber: eventSubscriber}
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

	ctx.Writer.WriteHeaderNow()
	flusher.Flush()

	pingTicker := time.NewTicker(60 * time.Second)
	defer pingTicker.Stop()

	reqCtx := ctx.Request.Context()

	createCh, unsubscribe1, err := sC.eventSubscriber.Subscribe(ports.EventRoomCreate)
	if err != nil {
		return
	}
	defer unsubscribe1()

	deleteCh, unsubscribe2, err := sC.eventSubscriber.Subscribe(ports.EventRoomDelete)
	if err != nil {
		return
	}
	defer unsubscribe2()

	for {
		select {
		case <-createCh:
			if !sC.sendEvent(ctx, flusher, "RoomsChange", time.Now().Format(time.RFC3339)) {
				return
			}

		case <-deleteCh:
			if !sC.sendEvent(ctx, flusher, "RoomsChange", time.Now().Format(time.RFC3339)) {
				return
			}

		case <-pingTicker.C:
			if !sC.sendEvent(ctx, flusher, "Ping", time.Now().Format(time.RFC3339)) {
				return
			}

		case <-reqCtx.Done():
			return

		case <-sC.appCtx.Done():
			return
		}
	}
}

func (sC *SSEController) sendEvent(ctx *gin.Context, flusher http.Flusher, name, data string) bool {
	if _, err := fmt.Fprintf(ctx.Writer, "event: %s\ndata: %s\n\n", name, data); err != nil {
		return false
	}
	flusher.Flush()
	return true
}
