package models

type InsertUser struct {
	Id       string `avro:"id"`
	Email    string `avro:"email"`
	Password string `avro:"password"`
}

func (InsertUser) Schema() string {
	return `{
		"type": "record",
		"name": "InsertUser",
		"namespace": "com.yourcompany.auth",
		"fields": [
			{ "name": "id", "type": "string" },
			{ "name": "email", "type": "string" },
			{ "name": "password", "type": "string" }
		]
	}`
}
