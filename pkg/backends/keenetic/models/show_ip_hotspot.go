package models

type RciShowIpHotspot struct {
	Host   []Host `json:"host,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}
type Host struct {
	Mac        string `json:"mac,omitempty"`
	Via        string `json:"via,omitempty"`
	IP         string `json:"ip,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	Name       string `json:"name,omitempty"`
	Registered bool   `json:"registered,omitempty"`
	Link       string `json:"link,omitempty"`
}
