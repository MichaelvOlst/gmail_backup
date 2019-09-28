package models

// Message holds the ID for comparing
type Message struct {
	ID        string `json:"id" storm:"unique"`
	AccountID int    `json:"account_id"`
}
