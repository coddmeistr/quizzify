package user

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	InfoKey = "UserInfo"
)

const (
	userInfoHeader = "App-User-Info"
)

type Info struct {
	ID          int   `json:"id"`
	Permissions []int `json:"permissions"`
}

// List of permissions and their ids
// Ids are coming from SSO service
// So they should be synchronized with SSO service permission ids
const (
	Creator   = 1
	Moderator = 2
	Admin     = 3
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var (
				userInfo Info = Info{}
				wasError      = false
			)
			userJson := r.Header.Get(userInfoHeader)
			if err := json.Unmarshal([]byte(userJson), &userInfo); err != nil {
				wasError = true
			}

			if userInfo.ID == 0 {
				wasError = true
			}

			if !wasError {
				ctx = context.WithValue(ctx, InfoKey, userInfo)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserInfoFromContext(ctx context.Context) (Info, bool) {
	userInfo, ok := ctx.Value(InfoKey).(Info)
	return userInfo, ok
}
