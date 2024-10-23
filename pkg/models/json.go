package models

type Entity struct {
	Name    string   `json:"name"`
	Group   string   `json:"group"`
	Domains []string `json:"domains"`
	DNS     []string `json:"dns"`
	Timeout int      `json:"timeout"`
	IP4     []string `json:"ip4"`
	IP6     []string `json:"ip6"`
	Cidr4   []string `json:"cidr4"`
	Cidr6   []string `json:"cidr6"`
}
