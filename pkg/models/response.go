package models

type APIResponse struct {
	ResponseType string `json:"responseType"`
	SyncCounter  int    `json:"syncCounter"`
}
