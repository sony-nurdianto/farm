package models

type InsertAccount struct {
	Id       string `avro:"id" json:"id"`
	Email    string `avro:"email" json:"email"`
	Password string `avro:"password_hash" json:"password_hash"`
}

func (InsertAccount) Schema() string {
	return `
		{
		  "type": "record",
		  "name": "InsertAuthAccount",
		  "fields": [
		    {
		      "name": "id",
		      "type": "string",
		      "default": ""
		    },
		    {
		      "name": "email",
		      "type": "string",
		      "default": ""
		    },
		    {
		      "name": "password",
		      "type": "string",
		      "default": ""
		    }
		  ]
		}
	`
}
