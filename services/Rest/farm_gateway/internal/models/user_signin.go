package models

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
