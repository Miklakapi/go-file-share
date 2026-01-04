package middleware

import (
	"errors"
	"net/http"

	apierrors "github.com/Miklakapi/go-file-share/internal/api/api-errors"
	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		last := c.Errors.Last().Err
		httpErr := MapErrors(last)

		c.AbortWithStatusJSON(httpErr.Status, gin.H{
			"code":    httpErr.Code,
			"message": httpErr.Message,
		})
	}
}

func MapErrors(err error) HTTPError {
	switch {

	// ======================
	// ROOM
	// ======================
	case errors.Is(err, ports.ErrRoomAlreadyExists):
		return HTTPError{Status: http.StatusConflict, Code: "ROOM_ALREADY_EXISTS", Message: "Room already exists"}

	case errors.Is(err, ports.ErrRoomNotFound),
		errors.Is(err, domain.ErrRoomNotFound):
		return HTTPError{Status: http.StatusNotFound, Code: "ROOM_NOT_FOUND", Message: "Room not found"}

	case errors.Is(err, domain.ErrInvalidRoomTTL),
		errors.Is(err, domain.ErrRoomLifespanTooLong):
		return HTTPError{Status: http.StatusBadRequest, Code: "INVALID_ROOM_LIFESPAN", Message: "Invalid room lifespan"}

	// ======================
	// PASSWORD
	// ======================
	case errors.Is(err, domain.ErrEmptyPassword):
		return HTTPError{Status: http.StatusBadRequest, Code: "PASSWORD_EMPTY", Message: "Password is required"}

	case errors.Is(err, domain.ErrEmptyPasswordHash):
		return HTTPError{Status: http.StatusInternalServerError, Code: "PASSWORD_HASH_FAILED", Message: "Internal server error"}

	case errors.Is(err, domain.ErrInvalidPassword):
		return HTTPError{Status: http.StatusUnauthorized, Code: "INVALID_PASSWORD", Message: "Invalid password"}

	// ======================
	// TOKEN
	// ======================
	case errors.Is(err, domain.ErrEmptyToken):
		return HTTPError{Status: http.StatusBadRequest, Code: "TOKEN_EMPTY", Message: "Token is required"}

	case errors.Is(err, domain.ErrTokenLifespanTooLong):
		return HTTPError{Status: http.StatusBadRequest, Code: "TOKEN_LIFESPAN_TOO_LONG", Message: "Token lifespan is too long"}

	case errors.Is(err, ports.ErrInvalidToken),
		errors.Is(err, ports.ErrTokenExpired),
		errors.Is(err, ports.ErrTokenParse),
		errors.Is(err, ports.ErrTokenRoomMismatch),
		errors.Is(err, ports.ErrTokenSignAlgo),
		errors.Is(err, domain.ErrTokenNotFound):
		return HTTPError{Status: http.StatusUnauthorized, Code: "TOKEN_INVALID", Message: "Invalid or expired token"}

	// ======================
	// FILE
	// ======================
	case errors.Is(err, domain.ErrFileNotFound):
		return HTTPError{Status: http.StatusNotFound, Code: "FILE_NOT_FOUND", Message: "File not found"}

	case errors.Is(err, domain.ErrInvalidFile):
		return HTTPError{Status: http.StatusBadRequest, Code: "INVALID_FILE", Message: "Invalid file"}

	case errors.Is(err, ports.ErrEmptyFilename):
		return HTTPError{Status: http.StatusBadRequest, Code: "FILENAME_EMPTY", Message: "Filename is required"}

	case errors.Is(err, ports.ErrNilReader):
		return HTTPError{Status: http.StatusInternalServerError, Code: "FILE_STREAM_MISSING", Message: "Internal server error"}

	// ======================
	// CONFIG / SERVER
	// ======================
	case errors.Is(err, ports.ErrEmptyUploadDir),
		errors.Is(err, ports.ErrInvalidUploadDir):
		return HTTPError{Status: http.StatusInternalServerError, Code: "SERVER_MISCONFIGURED", Message: "Server configuration error"}

	case errors.Is(err, ports.ErrPublishPanic):
		return HTTPError{Status: http.StatusInternalServerError, Code: "EVENTBUS_ERROR", Message: "Internal server error"}

	// ======================
	// API
	// ======================
	case errors.Is(err, apierrors.ErrInvalidRequest):
		return HTTPError{Status: http.StatusBadRequest, Code: "INVALID_REQUEST", Message: "Invalid request payload"}

	// ======================
	// FALLBACK
	// ======================
	default:
		return HTTPError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Internal server error"}
	}
}
