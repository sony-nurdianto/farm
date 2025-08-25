package models

type Users struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Verified     bool   `json:"verified"`
	RegisteredAt string `json:"registered_at"`
	UpdatedAt    string `json:"updated_at"`
}
