package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Miklakapi/go-file-share/internal/app"
	"github.com/gin-gonic/gin"
)

type HtmlController struct {
	Deps *app.DependencyBag
}

func NewHtmlController(deps *app.DependencyBag) *HtmlController {
	return &HtmlController{Deps: deps}
}

func (h *HtmlController) Index(ctx *gin.Context) {
	h.serveIndex(ctx)
}

func (h *HtmlController) SPAFallback(ctx *gin.Context) {
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

func (h *HtmlController) Assets() string {
	return h.Deps.Config.PublicDir + "/assets"
}

func (h *HtmlController) Favicon() string {
	return h.Deps.Config.PublicDir + "/favicon.ico"
}

func (h *HtmlController) serveIndex(ctx *gin.Context) {
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
