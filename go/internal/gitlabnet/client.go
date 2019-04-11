package gitlabnet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/gitlab-shell/go/internal/config"
)

const (
	internalApiPath  = "/api/v4/internal"
	secretHeaderName = "Gitlab-Shared-Secret"
)

type GitlabClient interface {
	Get(path string) (*http.Response, error)
	Post(path string, data interface{}) (*http.Response, error)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func GetClient(config *config.Config) (GitlabClient, error) {
	url := config.GitlabUrl
	if strings.HasPrefix(url, UnixSocketProtocol) {
		return buildSocketClient(config), nil
	}
	if strings.HasPrefix(url, HttpProtocol) {
		return buildHttpClient(config), nil
	}

	return nil, fmt.Errorf("Unsupported protocol")
}

func normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if !strings.HasPrefix(path, internalApiPath) {
		path = internalApiPath + path
	}
	return path
}

func newRequest(method, host, path string, data interface{}) (*http.Request, error) {
	path = normalizePath(path)

	var jsonReader io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		jsonReader = bytes.NewReader(jsonData)
	}

	request, err := http.NewRequest(method, host+path, jsonReader)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func parseError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	defer resp.Body.Close()
	parsedResponse := &ErrorResponse{}

	if err := json.NewDecoder(resp.Body).Decode(parsedResponse); err != nil {
		return fmt.Errorf("Internal API error (%v)", resp.StatusCode)
	} else {
		return fmt.Errorf(parsedResponse.Message)
	}

}

func doRequest(client *http.Client, config *config.Config, request *http.Request) (*http.Response, error) {
	encodedSecret := base64.StdEncoding.EncodeToString([]byte(config.Secret))
	request.Header.Set(secretHeaderName, encodedSecret)

	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Internal API unreachable")
	}

	if err := parseError(response); err != nil {
		return nil, err
	}

	return response, nil
}
