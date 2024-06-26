package appjwt

import (
	"time"

	"github.com/coddmeistr/quizzify/backend/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, perms []int, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID
	claims["permissions"] = perms

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
