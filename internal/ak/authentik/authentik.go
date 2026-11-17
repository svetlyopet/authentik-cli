package authentik

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	api "goauthentik.io/api/v3"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
)

const (
	apiBasePath    = "/api/v3"
	requestTimeout = 20 * time.Second
)

// authentik implements ak.AuthentikRepository on top of the official
// goauthentik.io/api/v3 client.
type authentik struct {
	url    string
	token  string
	client *api.APIClient
}

// New returns an ak.AuthentikRepository talking to the authentik instance at
// targetURL, authenticating with the given API token. A missing URL or token is
// reported on first use so commands such as `config` keep working without one.
func New(targetURL, token string) ak.AuthentikRepository {
	return &authentik{
		url:    targetURL,
		token:  token,
		client: newAPIClient(targetURL, token),
	}
}

func newAPIClient(targetURL, token string) *api.APIClient {
	cfg := api.NewConfiguration()
	// The server URL carries scheme, host and any path prefix the instance is
	// served under, so leave Configuration.Scheme/Host empty.
	cfg.Servers = api.ServerConfigurations{
		{URL: strings.TrimSuffix(targetURL, "/") + apiBasePath},
	}
	cfg.UserAgent = fmt.Sprintf("%s; %s", constants.CmdName, cfg.UserAgent)
	cfg.AddDefaultHeader("Authorization", "Bearer "+token)
	cfg.HTTPClient = &http.Client{Timeout: requestTimeout}

	return api.NewAPIClient(cfg)
}

func (a *authentik) preflightCheck() error {
	if a.url == "" || a.token == "" {
		return fmt.Errorf("URL and token for an authentik target are not set.\n Hint: run the %s config command first", constants.CmdName)
	}

	u, err := url.Parse(a.url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid authentik URL %q: expected the form https://authentik.example.com", a.url)
	}

	return nil
}

func (a *authentik) GetAuthentikTargetUrl() string {
	return a.url
}

// apiError converts an error returned by the API client into the CLI's error
// types: a 404 becomes NotExists, any other non-2xx response becomes
// UnexpectedResult carrying the response body, and transport or decoding
// errors are wrapped as-is.
func apiError(op string, resp *http.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return customErrors.NewNotExists(fmt.Sprintf("%s: not found", op))
	}

	var openAPIErr *api.GenericOpenAPIError
	if resp != nil && resp.StatusCode < http.StatusMultipleChoices && errors.As(err, &openAPIErr) {
		// A successful response the client could not decode, typically because
		// the server's API version does not match the client's.
		return customErrors.NewUnexpectedResult(fmt.Sprintf("%s: decoding response: %s", op, openAPIErr.Error()))
	}
	if errors.As(err, &openAPIErr) {
		detail := string(openAPIErr.Body())
		if detail == "" {
			detail = openAPIErr.Error()
		}
		return customErrors.NewUnexpectedResult(fmt.Sprintf("%s: %s", op, detail))
	}

	return fmt.Errorf("%s: %w", op, err)
}

// toInt32 narrows a primary key to the int32 the API client expects.
func toInt32(op string, v int) (int32, error) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, fmt.Errorf("%s: id %d is out of range", op, v)
	}
	return int32(v), nil
}
