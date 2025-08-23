package gokeenrestapimodels

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

type Import struct {
	Import   string `json:"import"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

type CreatedInterface struct {
	Intersects string `json:"intersects"`
	Created    string `json:"created"`
	Status     []struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Ident   string `json:"ident"`
		Message string `json:"message"`
	} `json:"status"`
}
