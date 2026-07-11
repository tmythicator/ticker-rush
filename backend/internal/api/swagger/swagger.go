// Package swagger handles serving OpenAPI/Swagger specifications.
package swagger

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed docs
var embedFS embed.FS

//go:embed index.html
var swaggerUIHTML []byte

// RegisterRoutes registers the swagger UI routes.
func RegisterRoutes(rg *gin.RouterGroup) {
	// Strip "docs" prefix so files like exchange/v1/exchange.swagger.json are at the root
	swaggerFS, err := fs.Sub(embedFS, "docs")
	if err != nil {
		panic("failed to create swagger sub-filesystem: " + err.Error())
	}

	// Serve the raw swagger JSON files
	rg.GET("/swagger/docs/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Param("filepath"), http.FS(swaggerFS))
	})

	// Serve the Swagger UI HTML page
	rg.GET("/swagger", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", swaggerUIHTML)
	})
}
