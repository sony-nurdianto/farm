package models

type UpdateUsers struct {
	ID       string  `avro:"id" json:"id"`
	FullName *string `avro:"full_name" json:"full_name"`
	Email    *string `avro:"email" json:"email"`
	Phone    *string `avro:"phone" json:"avro"`
}

type Users struct {
	ID           string `avro:"id" json:"id" redis:"id"`
	FullName     string `avro:"full_name" json:"full_name" redis:"full_name"`
	Email        string `avro:"email" json:"email" redis:"email"`
	Phone        string `avro:"phone" json:"phone" redis:"phone"`
	RegisteredAt string `avro:"registered_at" json:"registered_at" redis:"registered_at"`
	Verified     bool   `avro:"verified" json:"verified" redis:"verified"`
	UpdatedAt    string `avro:"updated_at" json:"updated_at" redis:"updated_at"`
}
