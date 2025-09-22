package repo

import (
	"context"

	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/models/animal"
)

func (rp ranchRepo) InsertAnimal(ctx context.Context, data animal.InsertAnimal) (res animal.InsertAnimal, _ error) {
	row := rp.ranchDB.insertAnimalStmt.QueryRowContext(
		ctx,
		data.ID,
		data.FarmID,
		data.Name,
		data.Species,
		data.Breed,
		data.BirthDate,
		data.OriginType,
		data.OriginFrom,
		data.PurchasedPrice,
		data.Gender,
		data.WeightKg,
		data.HealthStatus,
		data.RegistrationDate,
		data.UpdatedAt,
		data.IsActive,
		data.Notes,
	)

	if err := row.Scan(&res); err != nil {
		return res, err
	}

	return res, nil
}
