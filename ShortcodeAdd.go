package cms

import (
	"net/http"
)

func (cms *Cms) ShortcodeAdd(key string, fnRender func(r *http.Request, s string, m map[string]string) string) {
	cms.shortcodes[key] = fnRender
}
