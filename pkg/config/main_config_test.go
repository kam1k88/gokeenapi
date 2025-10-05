package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_ValidFile(t *testing.T) {
	configContent := `keenetic:
  url: "http://192.168.1.1"
  login: "admin"
  password: "password"
routes:
  - interfaceId: "Wireguard0"
    bat-file: ["routes.bat"]
dns:
  records:
    - domain: "test.local"
      ip: ["192.168.1.100"]
logs:
  debug: true`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	err = LoadConfig(configPath)
	assert.NoError(t, err)

	assert.Equal(t, "http://192.168.1.1", Cfg.Keenetic.URL)
	assert.Equal(t, "admin", Cfg.Keenetic.Login)
	assert.Equal(t, "password", Cfg.Keenetic.Password)
	assert.Len(t, Cfg.Routes, 1)
	assert.Equal(t, "Wireguard0", Cfg.Routes[0].InterfaceID)
	assert.Len(t, Cfg.DNS.Records, 1)
	assert.Equal(t, "test.local", Cfg.DNS.Records[0].Domain)
	assert.True(t, Cfg.Logs.Debug)
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	invalidContent := `keenetic:
  url: "http://192.168.1.1"
  login: "admin"
  password: "password"
routes:
  - interfaceId: "Wireguard0"
    bat-file: ["routes.bat"
dns:`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	err = LoadConfig(configPath)
	assert.Error(t, err)
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	// Clear environment variable
	_ = os.Unsetenv("GOKEENAPI_CONFIG")

	err := LoadConfig("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config path is empty")
}

func TestLoadConfig_FromEnvironment(t *testing.T) {
	configContent := `keenetic:
  url: "http://192.168.1.1"
  login: "admin"
  password: "password"`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "env_config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set environment variable
	_ = os.Setenv("GOKEENAPI_CONFIG", configPath)
	defer func() { _ = os.Unsetenv("GOKEENAPI_CONFIG") }()

	err = LoadConfig("")
	assert.NoError(t, err)
	assert.Equal(t, "http://192.168.1.1", Cfg.Keenetic.URL)
}

func TestLoadConfig_EnvironmentOverrides(t *testing.T) {
	configContent := `keenetic:
  url: "http://192.168.1.1"
  login: "admin"
  password: "password"`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set environment overrides
	_ = os.Setenv("GOKEENAPI_KEENETIC_LOGIN", "env_admin")
	_ = os.Setenv("GOKEENAPI_KEENETIC_PASSWORD", "env_password")
	defer func() {
		_ = os.Unsetenv("GOKEENAPI_KEENETIC_LOGIN")
		_ = os.Unsetenv("GOKEENAPI_KEENETIC_PASSWORD")
	}()

	err = LoadConfig(configPath)
	assert.NoError(t, err)

	assert.Equal(t, "http://192.168.1.1", Cfg.Keenetic.URL)
	assert.Equal(t, "env_admin", Cfg.Keenetic.Login)
	assert.Equal(t, "env_password", Cfg.Keenetic.Password)
}

func TestLoadConfig_DockerEnvironment(t *testing.T) {
	configContent := `keenetic:
  url: "http://192.168.1.1"
  login: "admin"
  password: "password"`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set Docker environment
	_ = os.Setenv("GOKEENAPI_INSIDE_DOCKER", "true")
	defer func() { _ = os.Unsetenv("GOKEENAPI_INSIDE_DOCKER") }()

	err = LoadConfig(configPath)
	assert.NoError(t, err)

	assert.Equal(t, "/etc/gokeenapi", Cfg.DataDir)
}
