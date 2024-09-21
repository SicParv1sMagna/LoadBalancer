package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeTempConfigFile(content string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "config-*.yml")
	if err != nil {
		return "", err
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		return "", err
	}
	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func TestMustLoad_Success(t *testing.T) {
	configContent := `
env: "test"
healthCheckInterval: "10s"
servers:
  - "http://localhost:8081"
  - "http://localhost:8082"
listenPort: ":9090"
`
	configPath, err := writeTempConfigFile(configContent)
	assert.NoError(t, err)
	defer os.Remove(configPath)

	os.Setenv("CONFIG_PATH", configPath)
	defer os.Unsetenv("CONFIG_PATH")

	cfg := MustLoad()

	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "10s", cfg.HealthCheckInterval)
	assert.Equal(t, ":9090", cfg.ListenPort)
	assert.Equal(t, []string{"http://localhost:8081", "http://localhost:8082"}, cfg.Servers)
}

func TestMustLoad_NoConfigPathEnv(t *testing.T) {
	os.Unsetenv("CONFIG_PATH")

	assert.Panics(t, func() {
		MustLoad()
	}, "Expected MustLoad to panic when CONFIG_PATH is not set")
}

func TestMustLoad_ConfigFileNotFound(t *testing.T) {
	os.Setenv("CONFIG_PATH", "/path/to/non-existent-yml")
	defer os.Unsetenv("CONFIG_PATH")

	assert.Panics(t, func() {
		MustLoad()
	}, "Expected MustLoad to panic when config file is not found")
}

func TestMustLoad_InvalidConfigContent(t *testing.T) {
	invalidContent := `invalid yaml`

	configPath, err := writeTempConfigFile(invalidContent)
	assert.NoError(t, err)
	defer os.Remove(configPath)

	os.Setenv("CONFIG_PATH", configPath)
	defer os.Unsetenv("CONFIG_PATH")

	assert.Panics(t, func() {
		MustLoad()
	}, "Expected MustLoad to panic when the config content is invalid")
}
