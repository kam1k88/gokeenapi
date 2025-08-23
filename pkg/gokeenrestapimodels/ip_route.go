package gokeenrestapimodels

type RciIpRoute struct {
	Network   string `json:"network"`
	Host      string `json:"host"`
	Mask      string `json:"mask"`
	Interface string `json:"interface"`
	Auto      bool   `json:"auto"`
}
