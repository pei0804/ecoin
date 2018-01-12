package app

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/pei0804/topicoin/httputil"
	"github.com/pei0804/topicoin/view"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type apiHandler func(context.Context, http.ResponseWriter, *http.Request) (int, interface{}, error)

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer func() {
		if rv := recover(); rv != nil {
			debug.PrintStack()
			log.Errorf(ctx, "panic: %s", rv)
			http.Error(w, http.StatusText(
				http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	status, res, err := h(ctx, w, r)
	if err != nil {
		log.Errorf(ctx, "error: %s", err)
		respondError(w, status, err)
		return
	}
	respondJSON(w, status, res)
	return
}

// respondJSON レスポンスとして返すjsonを生成して、writerに書き込む
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError レスポンスとして返すエラーを生成する
func respondError(w http.ResponseWriter, code int, err error) {
	if e, ok := err.(*httputil.HTTPError); ok {
		respondJSON(w, e.Code, e)
	} else if err != nil {
		he := httputil.HTTPError{
			Code:    code,
			Message: err,
		}
		respondJSON(w, code, he)
	}
}

type viewHandler func(context.Context, http.ResponseWriter, *http.Request) error

func (h viewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer func() {
		if rv := recover(); rv != nil {
			debug.PrintStack()
			log.Errorf(ctx, "panic: %s", rv)
			http.Error(w, http.StatusText(
				http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	r.ParseForm()
	var buf httputil.ResponseBuffer
	err := h(ctx, &buf, r)
	if err == nil {
		buf.WriteTo(w)
	} else if e, ok := err.(*httputil.HTTPError); ok {
		if e.Status >= 500 {
			// logError(r, err, nil)
		}
		// errorHandler(w, r, e.Status, e.Message)
		errorHandler(w, r, e.Status, e.Message)
	} else {
		// logError(r, err, nil)
		errorHandler(w, r, http.StatusInternalServerError, err)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	// errorとステータスをいい感じにする
	// デバックの時は、
	view.HTML(w, status, "404.html", map[string]interface{}{})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusNotFound, "404.html", map[string]interface{}{})
}
