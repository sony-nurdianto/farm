package models

type InsertUser struct {
	Id       string `avro:"id" json:"id"`
	Email    string `avro:"email" json:"email"`
	Password string `avro:"password_hash" json:"password_hash"`
}

func (InsertUser) Schema() string {
	return `{
		"type": "record",
		"name": "InsertUser",
		"fields": [
		{ "name": "id", "type": "string", "default": ""},
			{ "name": "email", "type": "string", "default": ""},
			{ "name": "password_hash", "type": "string","default": ""}
		]
	}`
}
