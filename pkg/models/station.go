package models

type Station struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	CallLetters string `json:"callLetters,omitempty"`
	Role        any    `json:"role,omitempty"`
	InternalId  int    `json:"internalId,omitempty"`
}
