package cms

import "net/http"

type ShortcodeInterface interface {
	Alias() string
	Description() string
	Render(r *http.Request, s string, m map[string]string) string
}
