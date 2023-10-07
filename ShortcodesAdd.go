package cms

import (
	"net/http"
)

func (cms *Cms) ShortcodesAdd(shortcodes map[string]func(r *http.Request, s string, m map[string]string) string) {
	for key, fnRender := range shortcodes {
		cms.ShortcodeAdd(key, fnRender)
	}
}
