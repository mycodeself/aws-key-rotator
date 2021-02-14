package circleci

import (
	"fmt"
	"net/http"
	"time"
)

type ContextAddOrUpdateEnvVarResponse struct {
	Variable  string    `json:"variable"`
	CreatedAt time.Time `json:"created_at"`
	ContextId string    `json:"context_id"`
}

// Create or update an environment variable within a context
func (c *Client) ContextAddOrUpdateEnvVar(contextId, envVarName, value string) (*ContextAddOrUpdateEnvVarResponse, error) {
	payload := &struct {
		Value string `json:"value"`
	}{value}

	path := fmt.Sprintf("context/%s/environment-variable/%s", contextId, envVarName)

	var response ContextAddOrUpdateEnvVarResponse

	err := c.request(http.MethodPut, path, nil, payload, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
