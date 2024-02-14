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

func TestAddPermission_OK(t *testing.T) {
	ctx, st := suits.NewDefault(t)

	// Register new test user
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

	// Login this user, check if login works
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

	// Test token claims, with checking permissions slice
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, login, claims["login"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), int(claims["exp"].(float64)), deltaSeconds)
	permsInterface := claims["permissions"].([]interface{})
	assert.Equal(t, 0, len(permsInterface))

	// Adding permissions to this user
	respPerms, err := st.PermsClient.AddPermission(ctx, &ssov1.AddPermissionRequest{
		UserId:       respReg.GetUserId(),
		PermissionId: 2,
	})
	require.NoError(t, err)
	assert.True(t, respPerms.GetGranted())
	respPerms, err = st.PermsClient.AddPermission(ctx, &ssov1.AddPermissionRequest{
		UserId:       respReg.GetUserId(),
		PermissionId: 3,
	})
	require.NoError(t, err)
	assert.True(t, respPerms.GetGranted())

	// Make a new login request to check if user got his permissions
	respLogin, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	token = respLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err = jwt.Parse(respLogin.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok = tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, []interface{}{2.0, 3.0}, claims["permissions"].([]interface{})) // Checking if user got his permissions
}

func TestRemovePermission_OK(t *testing.T) {
	ctx, st := suits.NewDefault(t)

	// Register new test user
	login := gofakeit.Username()
	email := gofakeit.Email()
	pass := randomPassword()
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    login,
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)

	// Login this user, check if login works
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

	permsInterface := claims["permissions"].([]interface{}) // Check only permission in map, cause we only need check perms
	assert.Equal(t, 0, len(permsInterface))                 // New user dont have any permissions

	// Add new permissions to this user
	respPerms, err := st.PermsClient.AddPermission(ctx, &ssov1.AddPermissionRequest{
		UserId:       respReg.GetUserId(),
		PermissionId: 2,
	})
	require.NoError(t, err)
	assert.True(t, respPerms.GetGranted())
	respPerms, err = st.PermsClient.AddPermission(ctx, &ssov1.AddPermissionRequest{
		UserId:       respReg.GetUserId(),
		PermissionId: 3,
	})
	require.NoError(t, err)
	assert.True(t, respPerms.GetGranted())

	// Make a new login request to check if user will get his permissions
	respLogin, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	token = respLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err = jwt.Parse(respLogin.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok = tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, []interface{}{2.0, 3.0}, claims["permissions"].([]interface{})) // Checking right permissions

	// Removing permissions from this user
	respRemPerms, err := st.PermsClient.RemovePermission(ctx, &ssov1.RemovePermissionRequest{
		UserId:       respReg.GetUserId(),
		PermissionId: 3,
	})
	require.NoError(t, err)
	assert.True(t, respRemPerms.GetRemoved())
	respRemPerms, err = st.PermsClient.RemovePermission(ctx, &ssov1.RemovePermissionRequest{
		UserId:       respReg.GetUserId(),
		PermissionId: 2,
	})
	require.NoError(t, err)
	assert.True(t, respRemPerms.GetRemoved())

	// Make a new login request to check if permission were deleted
	respLogin, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	token = respLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err = jwt.Parse(respLogin.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok = tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, []interface{}{}, claims["permissions"].([]interface{})) // Check for empty permissions slice
}
