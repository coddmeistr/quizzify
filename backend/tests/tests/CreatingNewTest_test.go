package tests

import (
	"fmt"
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

	req, err := http.NewRequest(http.MethodPost, host+createUrl, nil)
	resp, err := s.Client.Do(req)
	require.NoError(t, err)
	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	fmt.Println(string(bytes))
	require.Equal(t, http.StatusOK, resp.StatusCode)

}
