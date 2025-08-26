package models

import "time"

type Farm struct {
	ID          string
	FarmerID    string
	FarmName    string
	FarmType    string
	FarmSize    float32
	FarmStatus  string
	Description string
	AddressesID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FarmAddress struct {
	ID          string
	Street      string
	Village     string
	SubDistrict string
	City        string
	Province    string
	PostalCode  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FarmWithAddress struct {
	Farm
	FarmAddress
}
