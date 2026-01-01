package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type HtmlController struct {
	publicDir string
}

func NewHtmlController(publicDir string) *HtmlController {
	return &HtmlController{publicDir: publicDir}
}

func (hC *HtmlController) Index(ctx *gin.Context) {
	hC.serveIndex(ctx)
}

func (hC *HtmlController) SPAFallback(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	if strings.HasPrefix(path, "/api/") {
		ctx.Status(http.StatusNotFound)
		return
	}

	publicDir := hC.publicDir
	fullPath := filepath.Join(publicDir, filepath.Clean(path))
	if fileExists(fullPath) {
		ctx.File(fullPath)
		return
	}

	hC.serveIndex(ctx)
}

func (hC *HtmlController) Assets() string {
	return hC.publicDir + "/assets"
}

func (hC *HtmlController) Favicon() string {
	return hC.publicDir + "/favicon.ico"
}

func (hC *HtmlController) serveIndex(ctx *gin.Context) {
	indexPath := filepath.Join(hC.publicDir, "index.html")

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
