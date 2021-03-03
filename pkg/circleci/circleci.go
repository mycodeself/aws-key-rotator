package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
)

var (
	defaultURL = &url.URL{Host: "circleci.com", Scheme: "https", Path: "/api/v2/"}
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL    *url.URL
	ApiKey     string
	HTTPClient HTTPClient
}

type APIError struct {
	StatusCode int
	Message    string
}

func NewClient(apiKey string) *Client {
	c := &Client{
		ApiKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}

	return c
}

func NewClientFromEnv() *Client {
	apiKey := os.Getenv("CIRCLECI_TOKEN")
	if len(apiKey) <= 0 {
		log.Info().Msg("No CircleCI token provided in CIRCLECI_TOKEN env var")
		return nil
	}

	return NewClient(apiKey)
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Status Code %d: %s", e.StatusCode, e.Message)
}

func (c *Client) baseURL() *url.URL {
	if c.BaseURL == nil {
		return defaultURL
	}

	return c.BaseURL
}

func (c *Client) request(method, path string, params url.Values, payload interface{}, response interface{}) error {
	var requestBody io.Reader = nil

	url := c.baseURL().ResolveReference(&url.URL{Path: path, RawQuery: params.Encode()})

	if payload != nil {
		p, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		requestBody = bytes.NewReader(p)
	}

	req, err := http.NewRequest(method, url.String(), requestBody)

	if err != nil {
		return err
	}

	c.addRequestHeaders(req)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		errBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &APIError{
				StatusCode: res.StatusCode,
				Message:    fmt.Sprintf("Invalid response: %s", err),
			}
		}

		if len(errBody) > 0 {
			message := struct {
				Message string `json:"message"`
			}{}
			err = json.Unmarshal(errBody, &message)
			if err != nil {
				return &APIError{
					StatusCode: res.StatusCode,
					Message:    fmt.Sprintf("Unable to parse API error response: %s", errBody),
				}
			}
			return &APIError{
				StatusCode: res.StatusCode,
				Message:    message.Message,
			}
		}
	}

	if response != nil {
		err := json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) addRequestHeaders(req *http.Request) {
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Circle-Token", c.ApiKey)
}
