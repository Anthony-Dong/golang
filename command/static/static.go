package static

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed string_quote.html
//go:embed favicon.ico
var fs embed.FS

func registerStatic(r *gin.Engine) {
	r.StaticFS("/", http.FS(fs))
}
