package masterservice

import (
	"io"
	"net/http"
)

type HttpClient interface {
	Post(url string, body io.Reader) (*http.Response, error)
}

type httpClient struct{}

func NewHttpClient() HttpClient {
	return &httpClient{}
}

func (h *httpClient) Post(url string, body io.Reader) (*http.Response, error) {
	return http.Post(url, "application/json", body)
}
