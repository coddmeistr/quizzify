package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/api"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/numbers"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/tests/helpers"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/tests/suits"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

const (
	host      = "http://localhost:8080"
	createUrl = "/api/tests"
)

func TestCreatingNewTest_Ok(t *testing.T) {
	s := suits.NewDefault(t)

	test := helpers.GenerateRandomTest("strict_test")

	userID := numbers.RandomInt(1, 100)
	userInfo := helpers.UserInfo{
		ID:          userID,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(userInfo)
	require.NoError(t, err)
	test.CreatorID = &userID

	bts, err := json.Marshal(test)
	require.NoError(t, err)
	reader := io.NopCloser(bytes.NewReader(bts))
	req, err := http.NewRequest(http.MethodPost, host+createUrl, reader)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bts, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	fmt.Println(string(bts))
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NotNil(t, respBody.Payload)
	payload := *respBody.Payload
	assert.Equal(t, "test was created", payload.(string))

}
