package api

import (
	"context"
	"io"
	"net/http"
)

func MakeRequestWithContext(ctx context.Context, requestType string, u string, body io.ReadCloser) (*http.Response, error) {
	var client http.Client

	request, err := http.NewRequestWithContext(ctx, requestType, u, body)
	if err != nil {
		return nil, err
	}
	return client.Do(request)
}
