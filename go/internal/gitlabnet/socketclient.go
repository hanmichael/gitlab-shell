package gitlabnet

import (
	"context"
	"net"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/gitlab-shell/go/internal/config"
)

const (
	// We need to set the base URL to something starting with HTTP, the host
	// itself is ignored as we're talking over a socket.
	socketBaseUrl      = "http://unix"
	UnixSocketProtocol = "http+unix://"
)

type GitlabSocketClient struct {
	httpClient *http.Client
	config     *config.Config
}

func buildSocketClient(config *config.Config) *GitlabSocketClient {
	path := strings.TrimPrefix(config.GitlabUrl, UnixSocketProtocol)
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		},
		Timeout: config.HttpSettings.ReadTimeout,
	}

	return &GitlabSocketClient{httpClient: httpClient, config: config}
}

func (c *GitlabSocketClient) Get(path string) (*http.Response, error) {
	return c.doRequest("GET", path, nil)
}

func (c *GitlabSocketClient) Post(path string, data interface{}) (*http.Response, error) {
	return c.doRequest("POST", path, data)
}

func (c *GitlabSocketClient) doRequest(method, path string, data interface{}) (*http.Response, error) {
	request, err := newRequest(method, socketBaseUrl, path, data)
	if err != nil {
		return nil, err
	}

	return doRequest(c.httpClient, c.config, request)
}
