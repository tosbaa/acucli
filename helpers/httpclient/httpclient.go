// httpclient.go

package httpclient

import (
	"crypto/tls"
	"net/http"
)

const (
	BASE_URL = "https://194.16.0.72:3443/api/v1"
)

// MyHTTPClient is a custom HTTP client with default headers and insecure skip verification.
var MyHTTPClient http.Client

// headerTransport is a custom transport that sets default headers for each request.
type headerTransport struct {
	headers   map[string]string
	Transport http.RoundTripper
}

func CreateHttpClient(api_key string) {
	MyHTTPClient = http.Client{
		Transport: &headerTransport{
			headers: map[string]string{
				"Accept": "application/json",
				"X-Auth": api_key, // You can set a default value or provide it dynamically.
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// RoundTrip sets default headers for each request.
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Set default headers
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}

	// Use the underlying transport to perform the request
	return t.Transport.RoundTrip(req)
}
