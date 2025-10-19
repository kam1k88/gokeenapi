package models

type RciIpRoute struct {
	Network   string `json:"network"`
	Host      string `json:"host"`
	Mask      string `json:"mask"`
	Interface string `json:"interface"`
	Auto      bool   `json:"auto"`
}

type RciShowIpRoute struct {
	Destination string `json:"destination,omitempty"`
	Interface   string `json:"interface,omitempty"`
}
