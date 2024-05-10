package tests

import (
	"encoding/json"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/numbers"
	"github.com/coddmeistr/quizzify/backend/tests/tests/helpers"
	"github.com/coddmeistr/quizzify/backend/tests/tests/suits"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestGettingTest_Ok(t *testing.T) {
	s := suits.NewDefault(t)

	userID := numbers.RandomInt(1, 100)
	userInfo := helpers.UserInfo{
		ID:          userID,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(userInfo)
	require.NoError(t, err)

	testID, err := createNewTest(s.Client, userInfo.ID)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, host+getUrl+testID, nil)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bts, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, respBody.Payload)
	var gotTest helpers.GetTestResponse
	err = json.Unmarshal(respBody.Payload, &gotTest)
	require.NoError(t, err)

	require.Equal(t, testID, gotTest.ID)
	require.NotNil(t, gotTest.Questions)

	for _, q := range *gotTest.Questions {
		assert.Nil(t, q.Answer)
	}

}

func TestGettingTest_OkWithAnswers(t *testing.T) {
	s := suits.NewDefault(t)

	userID := numbers.RandomInt(1, 100)
	userInfo := helpers.UserInfo{
		ID:          userID,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(userInfo)
	require.NoError(t, err)

	testID, err := createNewTest(s.Client, userInfo.ID)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, host+getUrl+testID+"?withAnswers=true", nil)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bts, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, respBody.Payload)
	var gotTest helpers.GetTestResponse
	err = json.Unmarshal(respBody.Payload, &gotTest)
	require.NoError(t, err)

	require.Equal(t, testID, gotTest.ID)
	require.NotNil(t, gotTest.Questions)

	for _, q := range *gotTest.Questions {
		assert.NotNil(t, q.Answer)
	}

}

func TestGettingTest_OkWithoutAnswersOnDifferentUser(t *testing.T) {
	s := suits.NewDefault(t)

	userID := numbers.RandomInt(1, 100)
	userInfo := helpers.UserInfo{
		ID:          userID,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(userInfo)
	require.NoError(t, err)

	testID, err := createNewTest(s.Client, userInfo.ID+1)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, host+getUrl+testID, nil)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bts, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, respBody.Payload)
	var gotTest helpers.GetTestResponse
	err = json.Unmarshal(respBody.Payload, &gotTest)
	require.NoError(t, err)

	require.Equal(t, testID, gotTest.ID)
	require.NotNil(t, gotTest.Questions)

	for _, q := range *gotTest.Questions {
		assert.Nil(t, q.Answer)
	}

}

func TestGettingTest_FailNoRightsToGetWithAnswers(t *testing.T) {
	s := suits.NewDefault(t)

	userID := numbers.RandomInt(1, 100)
	userInfo := helpers.UserInfo{
		ID:          userID,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(userInfo)
	require.NoError(t, err)

	testID, err := createNewTest(s.Client, userInfo.ID+1)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, host+getUrl+testID+"?withAnswers=true", nil)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bts, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
	require.NotNil(t, respBody.Error)
}

func TestGettingTest_FailNotFound(t *testing.T) {
	s := suits.NewDefault(t)

	userID := numbers.RandomInt(1, 100)
	userInfo := helpers.UserInfo{
		ID:          userID,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(userInfo)
	require.NoError(t, err)

	_, err = createNewTest(s.Client, userInfo.ID)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, host+getUrl+"NOT_EXISTING_ID"+"?withAnswers=true", nil)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bts, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	require.NotNil(t, respBody.Error)
}
