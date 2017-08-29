package httphack

import (
	"net/http"
)

// HackTransport reverse proxy with Scheme => http, Host => HackHost
type HackTransport struct {
	HackHost string
}

func NewHackTransport(hackHost string) *HackTransport {
	return &HackTransport{HackHost: hackHost}
}

func (t *HackTransport) RoundTrip(originReg *http.Request) (*http.Response, error) {
	u := *originReg.URL
	u.Scheme = "http"
	u.Host = t.HackHost
	req, err := http.NewRequest(originReg.Method, u.String(), originReg.Body)
	if err != nil {
		return nil, err
	}
	return http.DefaultTransport.RoundTrip(req)
}

func NewHackClient(hackHost string) *http.Client {
	return &http.Client{Transport: NewHackTransport(hackHost)}
}
