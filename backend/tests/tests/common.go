package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api"
	"github.com/coddmeistr/quizzify/backend/tests/tests/helpers"
	"io"
	"net/http"
)

// Helper function to use it in integration tests
// It helps to simply create new test instance to test other features
// TestCreateTest should not use this function and write manual test creation
func createNewTest(client *http.Client, creatorID int) (testID string, err error) {
	test := helpers.GenerateRandomTest("strict_test")

	userInfo := helpers.UserInfo{
		ID:          creatorID,
		Permissions: []int{},
	}
	userBytes, err := json.Marshal(userInfo)
	if err != nil {
		return "", err
	}
	test.CreatorID = &creatorID

	bts, err := json.Marshal(test)
	if err != nil {
		return "", err
	}
	reader := io.NopCloser(bytes.NewReader(bts))
	req, err := http.NewRequest(http.MethodPost, host+createUrl, reader)
	req.Header.Set("Auth-User-Info", string(userBytes))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	bts, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var respBody api.Response
	err = json.Unmarshal(bts, &respBody)
	if err != nil {
		return "", err
	}
	if respBody.Payload == nil {
		return "", fmt.Errorf("payload is nil")
	}
	var newTestID string
	err = json.Unmarshal(respBody.Payload, &newTestID)
	if err != nil {
		return "", err
	}

	return newTestID, nil
}
