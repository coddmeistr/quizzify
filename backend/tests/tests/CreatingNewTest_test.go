package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coddmeistr/quizzify-online-tests/backend/tests/pkg/api"
	"github.com/coddmeistr/quizzify-online-tests/backend/tests/pkg/numbers"
	p "github.com/coddmeistr/quizzify-online-tests/backend/tests/pkg/pointer"
	"github.com/coddmeistr/quizzify-online-tests/backend/tests/tests/helpers"
	"github.com/coddmeistr/quizzify-online-tests/backend/tests/tests/suits"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"reflect"
	"testing"
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
	fmt.Println(string(bts))
	require.NoError(t, err)
	reader := io.NopCloser(bytes.NewReader(bts))
	req, err := http.NewRequest(http.MethodPost, host+createUrl, reader)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bts, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	fmt.Println(string(bts))
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NotNil(t, respBody.Payload)
	var newTestID string
	err = json.Unmarshal(respBody.Payload, &newTestID)
	require.NoError(t, err)

	getUrlWithID := getUrl + fmt.Sprint(newTestID)
	req, err = http.NewRequest(http.MethodGet, host+getUrlWithID+"?withAnswers=true", nil)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err = s.Client.Do(req)
	require.NoError(t, err)
	bts, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(bts, &respBody)
	require.NoError(t, err)
	fmt.Println(string(bts))
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, respBody.Payload)
	payload := respBody.Payload
	fmt.Println(string(payload))
	var dtest helpers.GetTestResponse
	err = json.Unmarshal(payload, &dtest)
	require.NoError(t, err)
	// all field checks
	assert.Equal(t, newTestID, dtest.ID)
	assert.Equal(t, test.CreatorID, dtest.CreatorID)
	assert.Equal(t, test.Type, dtest.Type)
	assert.Equal(t, test.Title, dtest.Title)
	assert.Equal(t, test.ShortText, dtest.ShortText)
	assert.Equal(t, test.LongText, dtest.LongText)
	assert.Equal(t, test.MainImage, dtest.MainImage)
	assert.Equal(t, test.Tags, dtest.Tags)
	assert.True(t, reflect.DeepEqual(test.Questions, dtest.Questions))

}

func TestCreateTest_FailCases(t *testing.T) {
	s := suits.NewDefault(t)

	authUser := helpers.UserInfo{
		ID:          3,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(authUser)
	require.NoError(t, err)
	validTest := helpers.GenerateRandomTest("strict_test")
	validTest.CreatorID = &authUser.ID

	tc := []struct {
		name       string
		test       helpers.Test
		mutator    func(test *helpers.Test)
		wantStatus int
		wantError  bool
	}{
		{
			name:       "no required field",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
			test:       validTest,
			mutator: func(test *helpers.Test) {
				test.ShortText = nil
			},
		},
		{
			name:       "no rights",
			wantStatus: http.StatusForbidden,
			wantError:  true,
			test:       validTest,
			mutator: func(test *helpers.Test) {
				test.CreatorID = p.Int(*test.CreatorID + 1)
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mutator != nil {
				tt.mutator(&tt.test)
			}
			bts, err := json.Marshal(tt.test)
			require.NoError(t, err)
			reader := io.NopCloser(bytes.NewReader(bts))
			req, err := http.NewRequest(http.MethodPost, host+createUrl, reader)
			req.Header.Set("Auth-User-Info", string(userBytes))
			resp, err := s.Client.Do(req)
			require.NoError(t, err)
			bts, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			var respBody api.Response
			err = json.Unmarshal(bts, &respBody)
			require.NoError(t, err)
			require.Equal(t, tt.wantStatus, resp.StatusCode)
			if tt.wantError {
				require.NotNil(t, respBody.Error)
			} else {
				require.Nil(t, respBody.Error)
			}
		})
	}

}
