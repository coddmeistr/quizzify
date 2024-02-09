package paginate

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maxik12233/quizzify-online-tests/server/pkg/api"
)

const (
	OptionsContextKey = "paginate_options"
)

type Options struct {
	Page    int
	PerPage int
}

func Middleware(defaultPage int, defaultPerPage int) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageInQuery := c.Request.URL.Query().Get("page")
		perPageInQuery := c.Request.URL.Query().Get("per_page")

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
				c.Writer.WriteHeader(http.StatusBadRequest)
				c.Writer.Write(errResponse.Marshal())
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
				c.Writer.WriteHeader(http.StatusBadRequest)
				c.Writer.Write(errResponse.Marshal())
				return
			}
		}

		options := Options{
			Page:    page,
			PerPage: perPage,
		}

		ctx := context.WithValue(c.Request.Context(), OptionsContextKey, options)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
