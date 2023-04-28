package models

type InitializeRequest struct {
	StoreName string `json:"storeName"`
	Threshold int    `json:"threshold"`
	Total     int    `json:"total"`
}

type UnlockRequest struct {
	StoreName string   `json:"storeName"`
	Shares    []string `json:"shares"`
}
