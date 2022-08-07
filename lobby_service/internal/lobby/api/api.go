package api

import (
	"io"
	"net/http"
	"net/url"
)

func MakeRequest(requestType string, u string, body io.ReadCloser) (*http.Response, error) {
	var client http.Client

	_url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	request := http.Request{
		Method: requestType,
		URL:    _url,
		Body:   body,
	}
	return client.Do(&request)

}
