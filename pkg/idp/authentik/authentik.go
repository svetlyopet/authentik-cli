package authentik

import (
	"fmt"
	"io"
	"time"

	"net/http"
	"net/url"

	"github.com/svetlyopet/authentik-cli/pkg/idp"
)

type authentik struct {
	url    string
	token  string
	client http.Client
}

func New(url, token string) idp.AuthentikRepository {
	return &authentik{
		url:   url,
		token: token,
		client: http.Client{
			Timeout: time.Duration(20) * time.Second,
		},
	}
}

func (a authentik) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	return a.doRequestWithQuery(method, url, body, nil)
}

func (a authentik) doRequestWithQuery(method, url string, body io.Reader, values *url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "authentik-cli; go-http-client/1.1")

	if values != nil {
		req.URL.RawQuery = values.Encode()
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
