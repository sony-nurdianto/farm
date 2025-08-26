package models

type UpdateUsers struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}
