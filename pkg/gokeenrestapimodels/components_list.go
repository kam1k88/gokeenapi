package gokeenrestapimodels

// RciComponentsList represents the response from /rci/components/list endpoint
type RciComponentsList struct {
	Sandbox   string               `json:"sandbox"`
	Component map[string]Component `json:"component"`
}

// Component represents a single component available on the Keenetic router
type Component struct {
	Group     string `json:"group,omitempty"`
	Installed string `json:"installed,omitempty"`
	Libndm    string `json:"libndm,omitempty"`
	Version   string `json:"version,omitempty"`
}
