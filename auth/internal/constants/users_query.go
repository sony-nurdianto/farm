package constants

const (
	QUERY_CREATE_USERS string = `
		insert into %s 
			(email, password_hash)
		values
			($1,$2)
		returning id, email, password_hash, created_at, updated_at 
	`
)
