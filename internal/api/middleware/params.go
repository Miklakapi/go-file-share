package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CtxRoomIDKey = "roomID"
	CtxFileIDKey = "fileID"
)

func RoomIDParam() gin.HandlerFunc {
	return UUIDParam(CtxRoomIDKey, CtxRoomIDKey)
}

func FileIDParam() gin.HandlerFunc {
	return UUIDParam(CtxFileIDKey, CtxFileIDKey)
}

func UUIDParam(paramName string, ctxKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw := strings.TrimSpace(ctx.Param(paramName))
		if raw == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}

		id, err := uuid.Parse(raw)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}

		ctx.Set(ctxKey, id)
		ctx.Next()
	}
}
