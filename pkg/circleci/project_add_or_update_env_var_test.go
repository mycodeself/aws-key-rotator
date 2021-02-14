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

func TestProjectAddOrUpdateEnvVar(t *testing.T) {
	projectSlug := "github/test/test"
	envVarName := "TEST_VAR_NAME"
	value := "test_value"

	expectedMethod := http.MethodPost
	expectedUrl := fmt.Sprintf("https://circleci.com/api/v2/project/%s/envvar", projectSlug)

	// create response
	response := fmt.Sprintf(`{"name": "%s","value": "%s"}`, envVarName, value)
	res := ioutil.NopCloser(bytes.NewReader([]byte(response)))

	c := &Client{
		ApiKey: "TOKEN",
	}

	c.HTTPClient = &mock.HTTPClientDoMock{
		DoFn: func(req *http.Request) (*http.Response, error) {
			reqJson := &struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}{}

			err := json.NewDecoder(req.Body).Decode(reqJson)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, expectedMethod, req.Method)
			assert.Equal(t, expectedUrl, req.URL.String())
			assert.Equal(t, value, reqJson.Value)
			assert.Equal(t, envVarName, reqJson.Name)

			return &http.Response{
				Status: http.StatusText(201),
				Body:   res,
			}, nil
		},
	}

	r, _ := c.ProjectAddOrUpdateEnvVar(projectSlug, envVarName, value)

	assert.Equal(t, envVarName, r.Name)
	assert.Equal(t, value, r.Value)
}
