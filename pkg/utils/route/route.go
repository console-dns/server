package route

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/console-dns/server/pkg/content"
	self_errors "github.com/console-dns/server/pkg/errors"
	auth_model "github.com/console-dns/server/pkg/models/auth"
	logs_model "github.com/console-dns/server/pkg/models/logs"
	"github.com/pkg/errors"
)

func MVCBind(dh WebHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = dh(&WebRequest{
			Request:   r,
			W:         w,
			Content:   content.GetContent(r),
			LoginMeta: auth_model.GetLoginMeta(r),
			LogPusher: logs_model.GetContent(r),
			Hacker:    r.Context().Value("hacker").(auth_model.Hacker),
		})
		if err != nil {
			pushError(w, r, err)
			return
		}
	}
}

func ApiBind(ah ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := ah(&ApiRequest{
			Request:   r,
			Client:    auth_model.GetApiMeta(r).User().Name,
			Content:   content.GetContent(r),
			LogPusher: logs_model.GetContent(r),
			Body:      string(body),
		})
		if err != nil {
			pushError(w, r, err)
			return
		}
		if result == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		if d, ok := result.(string); ok {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(d))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		marshal, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
		w.Header().Set("ETag", fmt.Sprintf("\"%x\"", sha256.Sum256(marshal)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(marshal)
	}
}

func pushError(w http.ResponseWriter, r *http.Request, err error) {
	var status self_errors.HttpStatus
	if errors.As(err, &status) {
		status.Writer()(w, r)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
