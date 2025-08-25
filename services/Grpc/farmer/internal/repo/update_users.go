package repo

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/models"
)

func (fr farmerRepo) UpdateUser(ctx context.Context, users *models.UpdateUsers) (models.Users, error) {
	var res models.Users
	var registeredAt time.Time
	var updatedAt time.Time

	row := fr.farmerDB.updateUserStmt.QueryRowContext(
		ctx,
		users.FullName,
		users.Email,
		users.Phone,
		time.Now().UTC(),
		users.ID,
	)

	if err := row.Scan(
		&res.ID,
		&res.FullName,
		&res.Email,
		&res.Phone,
		&registeredAt,
		&res.Verified,
		&updatedAt,
	); err != nil {
		return res, err
	}

	res.RegisteredAt = registeredAt.Format(time.RFC3339)
	res.UpdatedAt = updatedAt.Format(time.RFC3339)

	return res, nil
}
