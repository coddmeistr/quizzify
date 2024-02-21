package sort

import (
	"context"
	"net/http"
	"strings"

	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/api"
)

const (
	ASC               = "asc"
	DESC              = "desc"
	OptionsContextKey = "sort_options"
)

type Options struct {
	Field string
	Order string
}

func Middleware(defaultSortField string, defaultSortOrder string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sortBy := r.URL.Query().Get("sort_by")
			sortOrder := r.URL.Query().Get("sort_order")

			if sortBy == "" {
				sortBy = defaultSortField
			}

			if sortOrder == "" {
				sortOrder = defaultSortOrder
			} else {
				upperSortOrder := strings.ToUpper(sortOrder)
				if upperSortOrder != ASC && upperSortOrder != DESC {
					w.WriteHeader(http.StatusBadRequest)
					apiErr := api.ErrorResponse{
						Message: "collation must be asc or desc",
						Details: nil,
					}
					// TODO: better error handling
					_, err := w.Write(apiErr.Marshal())
					if err != nil {
						return
					}
					return
				}
			}

			options := Options{
				Field: sortBy,
				Order: sortOrder,
			}
			ctx := context.WithValue(r.Context(), OptionsContextKey, options)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
