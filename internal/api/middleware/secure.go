package middleware

import "github.com/gin-gonic/gin"

func SecureHeaders(ctx *gin.Context) {
	ctx.Header("X-Frame-Options", "DENY")
	ctx.Header("X-Content-Type-Options", "nosniff")
	ctx.Header("Referrer-Policy", "strict-origin")
	ctx.Header(
		"Permissions-Policy",
		"geolocation=(),midi=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()",
	)
	ctx.Header(
		"Content-Security-Policy",
		"default-src 'self' data: blob:; "+
			"connect-src * data: blob:; "+
			"img-src * data: blob:; "+
			"style-src 'self' 'unsafe-inline'; "+
			"script-src 'self' 'unsafe-inline';",
	)

	ctx.Next()
}
