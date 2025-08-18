package models

type RciIpRoute struct {
	Network   string `json:"network"`
	Host      string `json:"host"`
	Mask      string `json:"mask"`
	Interface string `json:"interface"`
	Auto      bool   `json:"auto"`
}

type RciShowInterface struct {
	Id          string                             `json:"id"`
	Type        string                             `json:"type"`
	Description string                             `json:"description"`
	Address     string                             `json:"address"`
	Wireguard   RciShowInterfaceWireguardInterface `json:"wireguard,omitempty"`
	Connected   string                             `json:"connected"`
	Link        string                             `json:"link"`
	State       string                             `json:"state"`
}

type RciShowInterfaceWireguardInterface struct {
	PublicKey  string                                   `json:"public-key"`
	ListenPort int                                      `json:"listen-port"`
	Peer       []RciShowInterfaceWireguardInterfacePeer `json:"peer"`
}

type RciShowInterfaceWireguardInterfacePeer struct {
	PublicKey             string `json:"public-key"`
	RemoteEndpointAddress string `json:"remote-endpoint-address"`
	LocalEndpointAddress  string `json:"local-endpoint-address"`
	RemotePort            int    `json:"remote-port"`
	LocalPort             int    `json:"local-port"`
}
