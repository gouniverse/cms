package cms

import (
	"net/http"

	"github.com/gouniverse/shortcode"
)

// ContentRenderShortcodes renders the shortcodes in a string
func (cms *Cms) ContentRenderShortcodes(req *http.Request, content string) (string, error) {
	sh, err := shortcode.NewShortcode(shortcode.WithBrackets("<", ">"))

	if err != nil {
		return "", err
	}

	for _, shortcode := range cms.shortcodes {
		content = sh.RenderWithRequest(req, content, shortcode.Alias(), shortcode.Render())
	}

	return content, nil
}
