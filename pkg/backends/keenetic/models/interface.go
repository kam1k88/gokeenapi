package models

type RciShowInterface struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Connected   string `json:"connected"`
	Link        string `json:"link"`
	State       string `json:"state"`
	DefaultGw   bool   `json:"defaultgw,omitempty"`
}

type Import struct {
	Import   string `json:"import"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

type CreatedInterface struct {
	Created string `json:"created"`
	Status  []struct {
		Status  string `json:"status"`
		Code    string `json:"code"`
		Ident   string `json:"ident"`
		Message string `json:"message"`
	} `json:"status"`
}
