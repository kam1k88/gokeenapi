package config

import (
	"errors"
	"os"

	"github.com/noksa/gokeenapi/pkg/gokeenrestapimodels"
	"gopkg.in/yaml.v3"
)

var (
	Cfg = GokeenapiConfig{}
)

// Runtime holds runtime configuration that is not persisted to YAML
type Runtime struct {
	RouterInfo struct {
		Version gokeenrestapimodels.Version `yaml:"-"`
	} `yaml:"-"`
}

// GokeenapiConfig represents the main configuration structure for the application
type GokeenapiConfig struct {
	// Keenetic router connection settings
	Keenetic Keenetic `yaml:"keenetic"`
	// DataDir specifies custom data directory for storing application data (optional)
	DataDir string `yaml:"dataDir,omitempty"`
	// Routes contains list of routing configurations for different interfaces
	Routes []Route `yaml:"routes"`
	// DNS contains DNS records configuration
	DNS DNS `yaml:"dns"`
	// Logs contains logging configuration (optional)
	Logs Logs `yaml:"logs,omitempty"`
}

// Keenetic holds connection parameters for the Keenetic router
type Keenetic struct {
	// URL of the router (IP address or KeenDNS hostname with http/https)
	URL string `yaml:"url"`
	// Login for router admin access (can be overridden by GOKEENAPI_KEENETIC_LOGIN env var)
	Login string `yaml:"login"`
	// Password for router admin access (can be overridden by GOKEENAPI_KEENETIC_PASSWORD env var)
	Password string `yaml:"password"`
}

// Route defines routing configuration for a specific interface
type Route struct {
	// InterfaceID specifies the target interface (e.g., Wireguard0)
	InterfaceID string `yaml:"interfaceId"`
	// BatFile contains paths to local .bat files with route definitions
	BatFile []string `yaml:"bat-file"`
	// BatURL contains URLs to remote .bat files with route definitions
	BatURL []string `yaml:"bat-url"`
}

// DnsRecord represents a single DNS record with domain and IP addresses
type DnsRecord struct {
	// Domain name for the DNS record
	Domain string `yaml:"domain"`
	// IP addresses associated with the domain (supports multiple IPs)
	IP []string `yaml:"ip"`
}

// DNS contains DNS-related configuration
type DNS struct {
	// Records contains list of DNS records to manage
	Records []DnsRecord `yaml:"records"`
}

// Logs contains logging configuration options
type Logs struct {
	// Debug enables debug-level logging for troubleshooting
	Debug bool `yaml:"debug"`
}

func LoadConfig(configPath string) error {
	if configPath == "" {
		v, ok := os.LookupEnv("GOKEENAPI_CONFIG")
		if ok {
			configPath = v
		} else {
			return errors.New("config path is empty. Specify it via --config flag or GOKEENAPI_CONFIG environment variable")
		}
	}
	b, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &Cfg)
	if err != nil {
		return err
	}
	// read some sensitive variables and replace values if found
	v, ok := os.LookupEnv("GOKEENAPI_KEENETIC_LOGIN")
	if ok {
		Cfg.Keenetic.Login = v
	}
	v, ok = os.LookupEnv("GOKEENAPI_KEENETIC_PASSWORD")
	if ok {
		Cfg.Keenetic.Password = v
	}
	_, ok = os.LookupEnv("GOKEENAPI_INSIDE_DOCKER")
	if ok {
		Cfg.DataDir = "/etc/gokeenapi"
	}
	return nil
}
