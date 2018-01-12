package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ascarter/requestid"
	"github.com/google/uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Auth 認証（dbはフェイク）
func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "admin" {
				respondError(w, http.StatusUnauthorized, fmt.Errorf("利用権限がありません"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
	return
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		t := t2.Sub(t1)
		reqID, ok := requestid.FromContext(r.Context())
		if !ok {
			reqID = uuid.New().String()
		}
		log.Infof(ctx, "request_id %s req_time %s req_time_nsec %v", reqID, t.String(), t.Nanoseconds())
	})
}
