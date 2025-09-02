package models

import "time"

type Farm struct {
	ID          string    `redis:"id"`
	FarmerID    string    `redis:"farmer_id"`
	FarmName    string    `redis:"farm_name"`
	FarmType    string    `redis:"farm_type"`
	FarmSize    float64   `redis:"farm_size"`
	FarmStatus  string    `redis:"farm_status"`
	Description string    `redis:"description"`
	AddressesID string    `redis:"addresses_id"`
	CreatedAt   time.Time `redis:"created_at"`
	UpdatedAt   time.Time `redis:"updated_at"`
}

type FarmAddress struct {
	ID          string    `redis:"id"`
	Street      string    `redis:"street"`
	Village     string    `redis:"village"`
	SubDistrict string    `redis:"sub_district"`
	City        string    `redis:"city"`
	Province    string    `redis:"province"`
	PostalCode  string    `redis:"postal_code"`
	CreatedAt   time.Time `redis:"created_at"`
	UpdatedAt   time.Time `redis:"updated_at"`
}

type FarmWithAddress struct {
	Farm
	FarmAddress
}
