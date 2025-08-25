package models

type UpdateUsers struct {
	ID       string  `json:"id"`
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}
