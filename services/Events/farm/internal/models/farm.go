package models

import "time"

type Farm struct {
	ID          string    `avro:"id" redis:"id"`
	FarmerID    string    `avro:"farmer_id" redis:"farmer_id"`
	FarmName    string    `avro:"farm_name" redis:"farm_name"`
	FarmType    string    `avro:"farm_type" redis:"farm_type"`
	FarmSize    float64   `avro:"farm_size" redis:"farm_size"`
	FarmStatus  string    `avro:"farm_status" redis:"farm_status"`
	Description string    `avro:"description" redis:"description"`
	AddressesID string    `avro:"addresses_id" redis:"addresses_id"`
	CreatedAt   time.Time `avro:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `avro:"updated_at" redis:"updated_at"`
}

type FarmAddress struct {
	ID          string    `avro:"id" redis:"id"`
	Street      string    `avro:"street" redis:"street"`
	Village     string    `avro:"village" redis:"village"`
	SubDistrict string    `avro:"sub_district" redis:"sub_district"`
	City        string    `avro:"city" redis:"city"`
	Province    string    `avro:"province" redis:"province"`
	PostalCode  string    `avro:"postal_code" redis:"postal_code"`
	CreatedAt   time.Time `avro:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `avro:"updated_at" redis:"updated_at"`
}
