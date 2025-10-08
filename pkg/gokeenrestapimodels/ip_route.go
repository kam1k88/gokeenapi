package gokeenrestapimodels

type RciIpRoute struct {
	Network   string `json:"network"`
	Host      string `json:"host"`
	Mask      string `json:"mask"`
	Interface string `json:"interface"`
	Auto      bool   `json:"auto"`
}

type RciShowIpRoute struct {
	Destination string `json:"destination,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	Interface   string `json:"interface,omitempty"`
	Metric      int    `json:"metric,omitempty"`
	Flags       string `json:"flags,omitempty"`
	Rejecting   bool   `json:"rejecting,omitempty"`
	Proto       string `json:"proto,omitempty"`
	Floating    bool   `json:"floating,omitempty"`
	Static      bool   `json:"static,omitempty"`
}
