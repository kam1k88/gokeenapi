package models

type SecurityLevel struct {
	Public bool `json:"public"`
}
type Address struct {
	Address string `json:"address"`
	Mask    string `json:"mask"`
}
type NameServer struct {
	NameServer string `json:"name-server"`
}
type AdjustMss struct {
	Pmtu bool `json:"pmtu"`
}
type TCP struct {
	AdjustMss AdjustMss `json:"adjust-mss"`
}
type IP struct {
	Address    Address      `json:"address"`
	Mtu        string       `json:"mtu"`
	NameServer []NameServer `json:"name-server"`
	TCP        TCP          `json:"tcp"`
}
type Asc struct {
	Jc   string `json:"jc"`
	Jmin string `json:"jmin"`
	Jmax string `json:"jmax"`
	S1   string `json:"s1"`
	S2   string `json:"s2"`
	H1   string `json:"h1"`
	H2   string `json:"h2"`
	H3   string `json:"h3"`
	H4   string `json:"h4"`
}
type Endpoint struct {
	Address string `json:"address"`
}
type KeepaliveInterval struct {
	Interval int `json:"interval"`
}
type AllowIps struct {
	Address string `json:"address"`
	Mask    string `json:"mask"`
}
type Peer struct {
	Key               string            `json:"key"`
	Comment           string            `json:"comment"`
	Endpoint          Endpoint          `json:"endpoint"`
	KeepaliveInterval KeepaliveInterval `json:"keepalive-interval"`
	PresharedKey      string            `json:"preshared-key"`
	AllowIps          []AllowIps        `json:"allow-ips"`
}
type Wireguard struct {
	Asc  Asc    `json:"asc"`
	Peer []Peer `json:"peer"`
}
type RciShowScInterface struct {
	Description   string        `json:"description"`
	SecurityLevel SecurityLevel `json:"security-level"`
	IP            IP            `json:"ip"`
	Wireguard     Wireguard     `json:"wireguard"`
}
