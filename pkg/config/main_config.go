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

type Runtime struct {
	RouterInfo struct {
		Version gokeenrestapimodels.Version `yaml:"-"`
	} `yaml:"-"`
}

type GokeenapiConfig struct {
	Keenetic Keenetic `yaml:"keenetic"`
	DataDir  string   `yaml:"dataDir,omitempty"`
	Routes   []Route  `yaml:"routes"`
	DNS      DNS      `yaml:"dns"`
	Logs     Logs     `yaml:"logs,omitempty"`
}
type Keenetic struct {
	URL      string `yaml:"url"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}
type Route struct {
	InterfaceID string   `yaml:"interfaceId"`
	BatFile     []string `yaml:"bat-file"`
	BatURL      []string `yaml:"bat-url"`
}
type DnsRecord struct {
	Domain string   `yaml:"domain"`
	IP     []string `yaml:"ip"`
}
type DNS struct {
	Records []DnsRecord `yaml:"records"`
}

type Logs struct {
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
