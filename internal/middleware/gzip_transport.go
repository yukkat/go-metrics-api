package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

// GzipTransport сжимает тело запроса перед отправкой на сервер.
type GzipTransport struct {
	Transport http.RoundTripper
}

func (t *GzipTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	if req.Body == nil {
		return transport.RoundTrip(req)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()

	if len(body) == 0 {
		return transport.RoundTrip(req)
	}

	var compressed bytes.Buffer
	zw := gzip.NewWriter(&compressed)
	if _, err := zw.Write(body); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}

	newReq := req.Clone(req.Context())
	newReq.Body = io.NopCloser(bytes.NewReader(compressed.Bytes()))
	newReq.ContentLength = int64(compressed.Len())
	newReq.Header.Set("Content-Encoding", "gzip")

	return transport.RoundTrip(newReq)
}
