package models

type RciIpRoute struct {
	Network   string `json:"network"`
	Host      string `json:"host"`
	Mask      string `json:"mask"`
	Interface string `json:"interface"`
	Auto      bool   `json:"auto"`
}

type RciShowInterface struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Address     string `json:"address"`
}
