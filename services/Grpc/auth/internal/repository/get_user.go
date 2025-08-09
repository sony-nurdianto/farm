package repository

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/entity"
)

func (rp authRepo) GetUserByEmail(email string) (user entity.Users, _ error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*500,
	)

	defer cancel()

	row := rp.getUserByEmailStmt.QueryRowContext(ctx, email)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
