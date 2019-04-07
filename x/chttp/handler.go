package chttp

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
)

// Handler is an http.Handler that renders a chttp.Response
type Handler func(*http.Request) Response

// HandlerFunc wraps a `chttp.Handler` in a standard `http` handler
func (h Handler) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := h(r)
		bs, err := res.Marshal()

		if err != nil {
			fmt.Println("INTERNAL DECODING ERROR: ", err, string(res.Data()))
			panic(err)
		}

		// adds the CSRF token to the requests here
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		w.WriteHeader(res.HTTPCode())
		_, err = w.Write(bs)

		if err != nil {
			panic(err)
		}
	}
}
