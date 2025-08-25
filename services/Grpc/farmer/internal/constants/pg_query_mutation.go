package constants

const (
	UserQueryUpdate = `
		UPDATE users
		SET 
			full_name = COALESCE($1,full_name),
			email = COALESCE($2,email),
			phone = COALESCE($3,phone),	
			updated_at = $4
		WHERE id = $5
	 	RETURNING id, full_name, email,phone, registered_at,verified, updated_at; 
	`

	AccountUpdateQuery = `
		UPDATE accounts
		SET
			email = $1
		WHERE id = $2
	`
)
