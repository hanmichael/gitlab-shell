package gitlabnet

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-shell/go/internal/config"
)

const (
	HttpProtocol = "http://"
)

type GitlabHttpClient struct {
	httpClient *http.Client
	config     *config.Config
}

func buildHttpClient(config *config.Config) *GitlabHttpClient {
	httpClient := &http.Client{
		Timeout: config.HttpSettings.ReadTimeout,
	}

	return &GitlabHttpClient{httpClient: httpClient, config: config}
}

func (c *GitlabHttpClient) Get(path string) (*http.Response, error) {
	return c.doRequest("GET", path, nil)
}

func (c *GitlabHttpClient) Post(path string, data interface{}) (*http.Response, error) {
	return c.doRequest("POST", path, data)
}

func (c *GitlabHttpClient) doRequest(method, path string, data interface{}) (*http.Response, error) {
	request, err := newRequest(method, c.config.GitlabUrl, path, data)
	if err != nil {
		return nil, err
	}

	user, password := c.config.HttpSettings.User, c.config.HttpSettings.Password
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}

	return doRequest(c.httpClient, c.config, request)
}
