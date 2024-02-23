package user

import (
	"context"
	"encoding/json"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/slice"
	"net/http"
)

const (
	AuthInfoKey    = "AuthUserInfo"
	SubjectInfoKey = "SubjectUserInfo" // Deprecated
)

const (
	authUserInfoHeader    = "Auth-User-Info"
	subjectUserInfoHeader = "Subject-User-Info" // Deprecated
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

// AuthMiddleware TODO: Remove this method, no longer needed
func AuthMiddleware(lowestRoleForIndependentAccess int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var (
				authUserInfo = Info{}
				subjUserInfo = Info{}
			)
			authUserJson := r.Header.Get(authUserInfoHeader)
			if err := json.Unmarshal([]byte(authUserJson), &authUserInfo); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if authUserInfo.ID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			subjUserJson := r.Header.Get(subjectUserInfoHeader)
			if err := json.Unmarshal([]byte(subjUserJson), &subjUserInfo); err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if authUserInfo.ID != subjUserInfo.ID && slice.MaxInt(authUserInfo.Permissions) < lowestRoleForIndependentAccess {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = writeAuthUserInfo(r)

			next.ServeHTTP(w, r)
		})
	}
}

func writeAuthUserInfo(r *http.Request) *http.Request {
	ctx := r.Context()

	var (
		userInfo = Info{}
		wasError = false
	)
	userJson := r.Header.Get(authUserInfoHeader)
	if err := json.Unmarshal([]byte(userJson), &userInfo); err != nil {
		wasError = true
	}

	if userInfo.ID == 0 {
		wasError = true
	}

	if !wasError {
		ctx = context.WithValue(ctx, AuthInfoKey, userInfo)
		r = r.WithContext(ctx)
	}

	return r
}

// writeSubjectUserInfo TODO: Remove this method, no longer needed
func writeSubjectUserInfo(r *http.Request) *http.Request {
	ctx := r.Context()

	var (
		userInfo = Info{}
		wasError = false
	)
	userJson := r.Header.Get(subjectUserInfoHeader)
	if err := json.Unmarshal([]byte(userJson), &userInfo); err != nil {
		wasError = true
	}

	if userInfo.ID == 0 {
		wasError = true
	}

	if !wasError {
		ctx = context.WithValue(ctx, SubjectInfoKey, userInfo)
		r = r.WithContext(ctx)
	}

	return r
}

func AuthUserFromContext(ctx context.Context) (Info, bool) {
	userInfo, ok := ctx.Value(AuthInfoKey).(Info)
	return userInfo, ok
}

// SubjectUserFromContext TODO: Remove this method, no longer needed
func SubjectUserFromContext(ctx context.Context) (Info, bool) {
	userInfo, ok := ctx.Value(SubjectInfoKey).(Info)
	return userInfo, ok
}
