//go:build prod
// +build prod

package resources

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/benbjohnson/hashfs"
)

var (
	//go:embed static
	StaticDirectory embed.FS
	StaticSys       = hashfs.NewFS(StaticDirectory)
)

func Handler() http.Handler {
	slog.Debug("static assets are embedded")
	return hashfs.FileServer(StaticSys)
}

func StaticPath(path string) string {
	return "/" + StaticSys.HashName("static/"+path)
}

func StaticHTMLContent(path string) templ.Component {
	content, err := StaticDirectory.ReadFile(StaticPath(path))
	if err != nil {
		return templ.Raw(fmt.Sprintf("<!-- Error reading %s: %v -->", path, err))
	}
	return templ.Raw(string(content))
}
