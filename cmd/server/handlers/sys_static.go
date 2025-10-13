package handlers

import (
	"path/filepath"

	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/utils"
)

func ServeFavicon(c *gin.Context, a utils.App) {
	c.Header("Cache-Control", "public, max-age=3600, immutable")
	c.File("static/img/icons/favicon.ico")
}

func ServeStatic(c *gin.Context, a utils.App) {
	c.Header("Cross-Origin-Opener-Policy", "same-origin")
	c.Header("Cross-Origin-Embedder-Policy", "require-corp")
	c.Header("Cache-Control", "public, max-age=0, no-cache, no-store, must-revalidate")

	// force /static/ path to avoid serving other files
	filePath := filepath.Join("static", filepath.Clean(c.Param("path")))

	c.File(filePath)
}
