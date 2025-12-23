package middleware

import (
	"net/http"
	"strings"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(deps *app.DependencyBag) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw, source := extractAuthToken(ctx)

		if raw == "" {
			abortUnauthorized(ctx, `Bearer realm="api"`, "Unauthorized: missing token")
			return
		}

		token, ok := parseBearerToken(raw, source)
		if !ok {
			abortUnauthorized(ctx, `Bearer error="invalid_request", error_description="Expected Bearer token"`, "Unauthorized: invalid token format")
			return
		}

		if err := deps.TokenService.Validate(deps.AppContext, token); err != nil {
			www, msg := mapJWTError(err)
			abortUnauthorized(ctx, www, msg)
			return
		}

		ctx.Next()
	}
}

func extractAuthToken(ctx *gin.Context) (raw string, source string) {
	raw = strings.TrimSpace(ctx.GetHeader("Authorization"))
	if raw != "" {
		return raw, "header"
	}

	c, err := ctx.Cookie("auth_token")
	if err == nil {
		c = strings.TrimSpace(c)
		if c != "" {
			return c, "cookie"
		}
	}

	return "", ""
}

func parseBearerToken(raw string, source string) (token string, ok bool) {
	parts := strings.Fields(raw)

	if source == "cookie" && len(parts) == 1 {
		return parts[0], true
	}

	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != "" {
		return parts[1], true
	}

	return "", false
}

func abortUnauthorized(ctx *gin.Context, wwwAuth string, msg string) {
	if wwwAuth != "" {
		ctx.Header("WWW-Authenticate", wwwAuth)
	}
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": msg,
	})
}

func mapJWTError(err error) (wwwAuth string, msg string) {
	msg = "Unauthorized: invalid or expired token"

	switch {
	// case errors.Is(err, utils.ErrTokenExpired):
	// 	wwwAuth = `Bearer error="invalid_token", error_description="Token expired"`
	// case errors.Is(err, utils.ErrTokenSignAlgo):
	// 	wwwAuth = `Bearer error="invalid_token", error_description="Unexpected signing method"`
	// case errors.Is(err, utils.ErrTokenParse):
	// 	wwwAuth = `Bearer error="invalid_token", error_description="Token could not be parsed"`
	// case errors.Is(err, utils.ErrTokenInvalid):
	// 	wwwAuth = `Bearer error="invalid_token", error_description="Token is invalid"`
	default:
		wwwAuth = `Bearer error="invalid_token", error_description="Token is invalid or expired"`
	}

	return wwwAuth, msg
}
