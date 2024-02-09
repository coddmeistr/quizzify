package apperror

import "net/http"

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("", "")

		//err := h(w, r)
	}
}
