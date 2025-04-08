package client_wrapper

import "net/http"

type HeaderRoundTripper struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (h *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}
	return h.Transport.RoundTrip(req)
}
