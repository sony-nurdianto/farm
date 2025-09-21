package models

import "time"

type AnimalProduct struct {
	ID             string
	AnimalID       string
	CategoryID     string
	QualityGradeID string
	Quantity       float32
	Unit           string
	PricePerUnit   float32
	TotalValue     float32
	ProuduceAt     time.Time
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
