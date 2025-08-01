package constants

const (
	QUERY_CREATE_USERS string = `
		insert into %s 
			(id,email, password_hash)
		values
			($1,$2,$3)
		returning id, email, password_hash, created_at, updated_at 
	`

	QUERY_GET_USER_BY_EMAIL string = `
		select * from %s
		where email = $1
	`
)
