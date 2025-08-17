package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (rp authRepo) GetUserByEmail(ctx context.Context, email string) (user entity.Users, _ error) {
	tracer := otel.Tracer("auth-service")
	_, span := tracer.Start(ctx, "Repo:GetUserByEmail")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "get_user_by_email"),
		attribute.String("layer", "repository"),
		attribute.String("db.operation", "SELECT"),
		attribute.String("user.email", email),
	)

	dbctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	span.AddEvent("executing_database_query")
	row := rp.getUserByEmailStmt.QueryRowContext(dbctx, email)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		if err == sql.ErrNoRows {
			span.SetStatus(codes.Ok, "User not found")
		} else {
			span.SetStatus(codes.Error, "Database query failed")
		}
		return user, err
	}

	span.SetAttributes(attribute.String("user.id", user.Id))
	span.SetStatus(codes.Ok, "User retrieved successfully")
	return user, nil
}
