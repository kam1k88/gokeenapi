package gokeenrestapimodels

type ParseRequest struct {
	Parse string `json:"parse"`
}

type ParseResponse struct {
	Parse Parse `json:"parse"`
}
type Status struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Ident   string `json:"ident"`
	Message string `json:"message"`
}
type Parse struct {
	Prompt string   `json:"prompt"`
	Status []Status `json:"status"`
}
