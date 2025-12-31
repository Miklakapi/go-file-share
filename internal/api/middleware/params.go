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

func SetRoomIDParam() gin.HandlerFunc {
	return UUIDParam(CtxRoomIDKey, CtxRoomIDKey)
}

func MustRoomIDParam(ctx *gin.Context) uuid.UUID {
	roomIdAny := ctx.MustGet(CtxRoomIDKey)
	roomId := roomIdAny.(uuid.UUID)
	return roomId
}

func SetFileIDParam() gin.HandlerFunc {
	return UUIDParam(CtxFileIDKey, CtxFileIDKey)
}

func MustFileIDParam(ctx *gin.Context) uuid.UUID {
	fileIdAny := ctx.MustGet(CtxFileIDKey)
	fileId := fileIdAny.(uuid.UUID)
	return fileId
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
