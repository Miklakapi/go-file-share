package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type PagesHandler struct {
	Deps *app.DependencyBag
}

func NewPagesHandler(deps *app.DependencyBag) *PagesHandler {
	return &PagesHandler{Deps: deps}
}

func (h *PagesHandler) Index(ctx *gin.Context) {
	h.serveIndex(ctx)
}

func (h *PagesHandler) SPAFallback(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	if strings.HasPrefix(path, "/api/") {
		ctx.Status(http.StatusNotFound)
		return
	}

	publicDir := h.Deps.Config.PublicDir

	fullPath := filepath.Join(publicDir, filepath.Clean(path))
	if fileExists(fullPath) {
		ctx.File(fullPath)
		return
	}

	h.serveIndex(ctx)
}

func (h *PagesHandler) serveIndex(ctx *gin.Context) {
	indexPath := filepath.Join(h.Deps.Config.PublicDir, "index.html")

	if !fileExists(indexPath) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Frontend not built (missing index.html)",
		})
		return
	}

	ctx.File(indexPath)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
