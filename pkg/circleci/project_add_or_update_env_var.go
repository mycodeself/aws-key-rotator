package circleci

import (
	"fmt"
	"net/http"
)

type ProjectAddOrUpdateEnvVarResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Creates a new environment variable. If the variable already exists, it updates it
func (c *Client) ProjectAddOrUpdateEnvVar(projectSlug, envVarName, value string) (*ProjectAddOrUpdateEnvVarResponse, error) {
	payload := &struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{envVarName, value}

	path := fmt.Sprintf("project/%s/envvar", projectSlug)

	var response ProjectAddOrUpdateEnvVarResponse

	err := c.request(http.MethodPost, path, nil, payload, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
