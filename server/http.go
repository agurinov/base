package server

import (
	"net"
	"net/http"

	"github.com/boomfunc/base/conf"
)

type HTTPWrapper struct {
	listener net.Listener
	router   *conf.Router
}

func newHTTPWrapper(listener net.Listener, router *conf.Router) *HTTPWrapper {
	return &HTTPWrapper{
		listener: listener,
		router:   router,
	}
}

func (wrp *HTTPWrapper) Serve() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		request, err := NewRequest(req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := handleRequest(request, wrp.router, w); err != nil {
			switch err {
			case conf.ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	})

	http.Serve(wrp.listener, nil)
}
