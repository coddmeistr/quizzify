package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/maxik12233/quizzify-online-tests/backend/protos/gen/go/sso"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/tests/suits"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID     = 1
	appSecret = "test-secret"

	passDefaultLen = 10
	deltaSeconds   = 1
)

func TestRegisterLogin_OK(t *testing.T) {
	ctx, st := suits.NewDefault(t)

	login := gofakeit.Username()
	email := gofakeit.Email()
	pass := randomPassword()
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	loginTime := time.Now()
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(respLogin.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, login, claims["login"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), int(claims["exp"].(float64)), deltaSeconds)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
