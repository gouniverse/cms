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

	for k, v := range cms.shortcodes {
		content = sh.RenderWithRequest(req, content, k, v)
	}

	return content, nil
}
