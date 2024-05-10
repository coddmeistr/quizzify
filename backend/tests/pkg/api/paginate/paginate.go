package paginate

import (
	"context"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api"
	"net/http"
	"strconv"
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
					api.WriteErrorMessage(w, http.StatusBadRequest, "INVALID_PAGE_NUMBER", "invalid page number")
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
					api.WriteErrorMessage(w, http.StatusBadRequest, "INVALID_PER_PAGE", "invalid per page number")
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
