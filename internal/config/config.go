package config

var (
	cfg = Configuration{}
)

type Configuration struct {
	Keenetic Keenetic `yaml:"keenetic"`
}

type Keenetic struct {
	Api      string `yaml:"api"`
	Login    string `yaml:"login,omitempty"`
	Password string `yaml:"password,omitempty"`
}
