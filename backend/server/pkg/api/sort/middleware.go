package sort

import (
	"context"
	"net/http"
	"strings"

	"github.com/maxik12233/quizzify-online-tests/server/pkg/api"
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

func Middleware(nextHandler http.HandlerFunc, defaultSortField, defaultSortOrder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				err := api.ErrorResponse{
					Message: "collation must be asc or desc",
					Details: nil,
				}
				w.Write(err.Marshal())
				return
			}
		}

		options := Options{
			Field: sortBy,
			Order: sortOrder,
		}
		ctx := context.WithValue(r.Context(), OptionsContextKey, options)
		r = r.WithContext(ctx)

		nextHandler(w, r)
	}
}
