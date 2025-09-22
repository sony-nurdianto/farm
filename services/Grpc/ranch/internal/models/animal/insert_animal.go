package animal

import "time"

type InsertAnimal struct {
	ID               string
	FarmID           string
	Name             string
	Species          string
	Breed            string
	BirthDate        *string
	OriginType       string
	OriginFrom       *string
	PurchasedPrice   *float32
	Gender           string
	WeightKg         float32
	HealthStatus     string
	RegistrationDate time.Time
	UpdatedAt        time.Time
	IsActive         bool
	Notes            string
}
