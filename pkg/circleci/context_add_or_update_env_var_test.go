package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mycodeself/aws-key-rotator/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestContextAddOrUpdateEnvVar(t *testing.T) {
	contextId := "b4580531-971d-4890-8d5d-4a75da866611"
	envVarName := "TEST_VAR_NAME"
	value := "test_value"

	expectedMethod := http.MethodPut
	expectedUrl := fmt.Sprintf("https://circleci.com/api/v2/context/%s/environment-variable/%s", contextId, envVarName)

	// create response
	response := fmt.Sprintf(`{"variable": "%s","created_at": "2015-09-21T17:29:21.042Z","context_id": "%s"}`, envVarName, contextId)
	res := ioutil.NopCloser(bytes.NewReader([]byte(response)))

	c := &Client{
		ApiKey: "TOKEN",
	}

	c.HTTPClient = &mock.HTTPClientDoMock{
		DoFn: func(req *http.Request) (*http.Response, error) {
			reqJson := &struct {
				Value string `json:"value"`
			}{}

			err := json.NewDecoder(req.Body).Decode(reqJson)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, expectedMethod, req.Method)
			assert.Equal(t, expectedUrl, req.URL.String())
			assert.Equal(t, value, reqJson.Value)

			return &http.Response{
				Status: http.StatusText(200),
				Body:   res,
			}, nil
		},
	}

	r, _ := c.ContextAddOrUpdateEnvVar(contextId, envVarName, value)

	assert.Equal(t, contextId, r.ContextId)
	assert.Equal(t, envVarName, r.Variable)
}
