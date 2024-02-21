package paginate

import (
	"context"
	"net/http"
	"strconv"

	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/api"
)

const (
	OptionsContextKey = "paginate_options"
)

type Options struct {
	Page    int
	PerPage int
}

func Middleware(defaultPage int, defaultPerPage int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pageInQuery := r.URL.Query().Get("page")
			perPageInQuery := r.URL.Query().Get("per_page")

			var page int
			if pageInQuery == "" {
				page = defaultPage
			} else {
				var err error
				page, err = strconv.Atoi(pageInQuery)
				if err != nil {
					errResponse := api.ErrorResponse{
						Message: "bad pagination values",
					}
					w.WriteHeader(http.StatusBadRequest)
					// TODO: better error handling
					_, err := w.Write(errResponse.Marshal())
					if err != nil {
						return
					}
					return
				}
			}

			var perPage int
			if perPageInQuery == "" {
				perPage = defaultPerPage
			} else {
				var err error
				perPage, err = strconv.Atoi(perPageInQuery)
				if err != nil {
					errResponse := api.ErrorResponse{
						Message: "per_page value is too large",
					}
					w.WriteHeader(http.StatusBadRequest)
					// TODO: better error handling
					_, err := w.Write(errResponse.Marshal())
					if err != nil {
						return
					}
					return
				}
			}

			options := Options{
				Page:    page,
				PerPage: perPage,
			}

			ctx := context.WithValue(r.Context(), OptionsContextKey, options)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
