//go:build !prod
// +build !prod

package resources

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
)

func Handler() http.Handler {
	slog.Info("static assets are being served directly", "path", StaticDirectoryPath)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.StripPrefix("/static/", http.FileServerFS(os.DirFS(StaticDirectoryPath))).ServeHTTP(w, r)
	})
}

func StaticPath(path string) string {
	return "/static/" + path
}

func StaticHTMLContent(path string) templ.Component {
	content, err := os.ReadFile(filepath.Join(StaticDirectoryPath, path))
	if err != nil {
		return templ.Raw(fmt.Sprintf("<!-- Error reading %s: %v -->", path, err))
	}
	return templ.Raw(string(content))
}
