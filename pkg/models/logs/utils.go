package logs

import "net/http"

var CONTENT = "_LOGS_"

func GetContent(r *http.Request) LogPusher {
	return r.Context().Value(CONTENT).(LogPusher)
}
