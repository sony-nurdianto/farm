package models

import (
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/pbgen"
)

type Animal struct {
	ID               string
	FarmID           string
	Name             string
	Species          string
	Breed            string
	BirthDate        string
	OriginType       pbgen.OriginType
	OriginFrom       string
	PurchasedPrice   float32
	Gender           pbgen.GENDER
	WeightKg         float32
	HealthStatus     pbgen.HealthStatus
	RegistrationDate time.Time
	IsActive         bool
	Notes            string
}
