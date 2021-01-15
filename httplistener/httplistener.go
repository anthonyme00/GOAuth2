package httplistener

import (
	"fmt"
	"net/http"

	"github.com/anthonyme00/GOAuth2/httplistener/util"
)

type httphandler struct {
	handler func(w http.ResponseWriter, req *http.Request)
}

func (handler httphandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.handler(w, req)
}

type HTTPListener struct {
	port   int
	server *http.Server
}

type ListenerHandle chan *http.Request

func (listener *HTTPListener) OpenListener(handler func(w http.ResponseWriter, req *http.Request)) ListenerHandle {
	listener.port = util.GetRandomOpenPort()

	handle := make(chan *http.Request)

	serverhandler := httphandler{
		handler: func(w http.ResponseWriter, req *http.Request) {
			handler(w, req)
			handle <- req
		},
	}

	listener.server = &http.Server{Addr: fmt.Sprintf(":%d", listener.port), Handler: serverhandler}
	go func() {
		listener.server.ListenAndServe()
	}()

	return handle
}

func (listener *HTTPListener) GetResponse(handle ListenerHandle) *http.Request {
	req := <-handle
	listener.server.Close()
	return req
}

var OpenInBrowser = util.OpenInBrowser

func (listener *HTTPListener) GetUrl() string {
	return fmt.Sprintf("http://127.0.0.1:%d/", listener.port)
}

func RedirectToURL(url string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "http://google.com/", 300)
	}
}
