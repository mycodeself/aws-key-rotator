package mock

import "net/http"

type HTTPClientDoMock struct {
	DoFn func(req *http.Request) (*http.Response, error)
}

func (m *HTTPClientDoMock) Do(req *http.Request) (*http.Response, error) {
	return m.DoFn(req)
}
