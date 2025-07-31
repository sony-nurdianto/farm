package entity

import "time"

type Users struct {
	Id        string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
