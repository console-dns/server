package content

import "net/http"

var CONTENT = "global-content"

func GetContent(r *http.Request) *Content {
	return r.Context().Value(CONTENT).(*Content)
}
