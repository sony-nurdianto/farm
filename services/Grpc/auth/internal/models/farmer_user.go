package models

type InsertFarmerUser struct {
	Id       string `avro:"id" json:"id"`
	FullName string `avro:"full_name" json:"full_name"`
	Email    string `avro:"email" json:"email"`
	Phone    string `avro:"phone" json:"phone"`
}

func (InsertFarmerUser) Schema() string {
	return `
		{
		  "type": "record",
		  "name": "InsertFarmerUser",
		  "fields": [
		    {
		      "name": "id",
		      "type": "string",
		      "default": ""
		    },
		    {
		      "name": "full_name",
		      "type": "string",
		      "default": ""
		    },
		    {
		      "name": "email",
		      "type": "string",
		      "default": ""
		    },
		    {
		      "name": "phone",
		      "type": "string",
		      "default": ""
		    }
		  ]
		}	
	`
}
