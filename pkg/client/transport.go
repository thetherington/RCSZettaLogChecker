package client

import "net/http"

type ZettaTransport struct {
	transport http.RoundTripper
	username  string
	password  string
	secret    string
}

func (z *ZettaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if z.username != "" && z.password != "" {
		req.SetBasicAuth(z.username, z.password)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if z.secret != "" {
		req.Header.Set("APIKEY", z.secret)
	}

	return z.transport.RoundTrip(req)
}
