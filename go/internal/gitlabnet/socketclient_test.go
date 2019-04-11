package gitlabnet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/gitlab-shell/go/internal/config"
)

func TestReadTimeoutSettingForSocketClient(t *testing.T) {
	expectedTimeout := time.Duration(500)

	config := &config.Config{HttpSettings: config.HttpSettingsConfig{ReadTimeout: expectedTimeout}}
	client := buildSocketClient(config)

	assert.Equal(t, expectedTimeout, client.httpClient.Timeout)
}
