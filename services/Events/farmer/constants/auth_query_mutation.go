package constants

const (
	AuthUpdateEmailAccount string = `
		UPDATE accounts
		SET
			email = $1
		WHERE id = $2
	`
)
