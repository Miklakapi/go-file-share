package main

import (
	"net/http"
	"strconv"

	"github.com/Miklakapi/go-file-share/internal/files"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	fileHub := files.NewFileHub()
	fileHub.Run()

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	engine.GET("/", func(ctx *gin.Context) {
		// HTML
	})
	engine.GET("/rooms", func(ctx *gin.Context) {
		// List of rooms
	})
	engine.GET("/rooms/:uuid", func(ctx *gin.Context) {
		// Room by uuid
	})
	engine.POST("/rooms", func(ctx *gin.Context) {
		password := ctx.PostForm("password")
		lifespanStr := ctx.PostForm("lifespan")

		if password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing field: 'password'"})
			return
		}
		lifespanSec := 3600
		if lifespanStr != "" {
			lifespanSec, err := strconv.Atoi(lifespanStr)
			if err != nil || lifespanSec <= 0 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'lifespan' (expected positive integer seconds)"})
				return
			}
		}

		newFileContainer := files.NewFileContainer(password, lifespanSec, fileHub)
		fileHub.Register <- newFileContainer
		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})
	engine.POST("/rooms/:uuid/auth", func(ctx *gin.Context) {
		// auth
	})
	engine.POST("/rooms/:uuid/auth/logout", func(ctx *gin.Context) {
		// logout
	})
	engine.DELETE("/rooms/:uuid", func(ctx *gin.Context) {
		// delete room
	})
	// File upload
	// POST /rooms/:uuid/files
	// Batch file upload
	// POST /rooms/:uuid/files/batch
	// List of files
	// GET /rooms/:uuid/files
	// Single file data
	// GET /rooms/:uuid/files/:fileId
	// Download file
	// GET /rooms/:uuid/files/:fileId/download
	// Remve file
	// DELETE /rooms/:uuid/files/:fileId

	// engine.POST("/upload", func(ctx *gin.Context) {
	// 	file, err := ctx.FormFile("file")
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{
	// 			"error": "Missing file: expected field 'file'",
	// 		})
	// 		return
	// 	}
	// 	newFileContainer := files.NewFileContainer(password, lifespanSec, fileHub)
	// 	fileHub.Register <- newFileContainer
	// 	ctx.JSON(http.StatusOK, gin.H{"ok": true})
	// })
}
