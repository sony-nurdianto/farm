package animal

import (
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/constants"
	"github.com/sony-nurdianto/farm/services/Grpc/ranch/internal/pbgen"
)

type Animal struct {
	ID               string
	FarmID           string
	Name             string
	Species          string
	Breed            string
	BirthDate        *string
	OriginType       pbgen.OriginType
	OriginFrom       *string
	PurchasedPrice   *float32
	Gender           pbgen.GENDER
	WeightKg         float32
	HealthStatus     pbgen.HealthStatus
	RegistrationDate time.Time
	UpdatedAt        time.Time
	IsActive         bool
	Notes            string
}

func (a Animal) originTypeString() string {
	switch a.OriginType {
	case pbgen.OriginType_ORIGIN_TYPE_BORN:
		return "born"
	case pbgen.OriginType_ORIGIN_TYPE_HATCHED:
		return "hatched"
	case pbgen.OriginType_ORIGIN_TYPE_PURCHASED:
		return "purchased"
	default:
		return "unspecified"
	}
}

func (a Animal) genderString() string {
	switch a.Gender {
	case pbgen.GENDER_GENDER_MALE:
		return "M"
	case pbgen.GENDER_GENDER_FEMALE:
		return "F"
	default:
		return constants.UNSPECIFIED
	}
}

func (a Animal) healthStatusString() string {
	switch a.HealthStatus {
	case pbgen.HealthStatus_HEALTH_STATUS_HEALTHY:
		return "healthy"
	case pbgen.HealthStatus_HEALTH_STATUS_SICK:
		return "sick"
	case pbgen.HealthStatus_HEALTH_STATUS_RECOVERED:
		return "recovered"
	case pbgen.HealthStatus_HEALTH_STATUS_DEAD:
		return "dead"
	case pbgen.HealthStatus_HEALTH_STATUS_SOLD:
		return "sold"
	default:
		return constants.UNSPECIFIED
	}
}

func (a Animal) FactoryInsertAnimal() InsertAnimal {
	return InsertAnimal{
		ID:               a.ID,
		FarmID:           a.FarmID,
		Name:             a.Name,
		Species:          a.Species,
		Breed:            a.Breed,
		BirthDate:        a.BirthDate,
		OriginType:       a.originTypeString(),
		OriginFrom:       a.OriginFrom,
		PurchasedPrice:   a.PurchasedPrice,
		Gender:           a.genderString(),
		WeightKg:         a.WeightKg,
		HealthStatus:     a.healthStatusString(),
		RegistrationDate: a.RegistrationDate,
		UpdatedAt:        a.UpdatedAt,
		IsActive:         a.IsActive,
		Notes:            a.Notes,
	}
}
