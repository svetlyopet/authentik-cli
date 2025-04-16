package authentik

import (
	"fmt"
	"io"
	"time"

	"net/http"
	"net/url"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

type authentik struct {
	url    string
	token  string
	client http.Client
}

func New(url, token string) ak.AuthentikRepository {
	return &authentik{
		url:   url,
		token: token,
		client: http.Client{
			Timeout: time.Duration(20) * time.Second,
		},
	}
}

func (a *authentik) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	return a.doRequestWithQuery(method, url, body, nil)
}

func (a *authentik) doRequestWithQuery(method, url string, body io.Reader, values *url.Values) (*http.Response, error) {
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", fmt.Sprintf("%s; go-http-client/1.1", constants.CmdName))

	if values != nil {
		req.URL.RawQuery = values.Encode()
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *authentik) preflightCheck() error {
	if a.url == "" || a.token == "" {
		return fmt.Errorf("URL and token for an authentik target are not set.\n Hint: run the %s config command first", constants.CmdName)
	}

	return nil
}

func (a *authentik) GetAuthentikTargetUrl() string {
	return a.url
}
