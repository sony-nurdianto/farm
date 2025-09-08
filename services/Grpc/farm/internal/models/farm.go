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

type AvrFarm struct {
	ID               string    `avro:"id" redis:"id" json:"id"`
	FarmerID         string    `avro:"farmer_id" redis:"farmer_id" json:"farmer_id"`
	FarmName         string    `avro:"farm_name" redis:"farm_name" json:"farm_name"`
	FarmType         string    `avro:"farm_type" redis:"farm_type" json:"farm_type"`
	FarmSize         float64   `avro:"farm_size" redis:"farm_size" json:"farm_size"`
	PhotoURL         *string   `avro:"photo_url" redis:"photo_url" json:"photo_url"`
	FarmStatus       string    `avro:"farm_status" redis:"farm_status" json:"farm_status"`
	Description      *string   `avro:"description" redis:"description" json:"description"`
	CreatedAt        time.Time `avro:"created_at" redis:"created_at" json:"created_at"`
	UpdatedAt        time.Time `avro:"updated_at" redis:"updated_at" json:"updated_at"`
	AddressID        string    `avro:"address_id" redis:"address_id" json:"address_id"`
	Street           string    `avro:"street" redis:"street" json:"street"`
	Village          string    `avro:"village" redis:"village" json:"village"`
	SubDistrict      string    `avro:"sub_district" redis:"sub_district" json:"sub_district"`
	City             string    `avro:"city" redis:"city" json:"city"`
	Province         string    `avro:"province" redis:"province" json:"province"`
	PostalCode       string    `avro:"postal_code" redis:"postal_code" json:"postal_code"`
	AddressCreatedAt time.Time `avro:"address_created_at" redis:"address_created_at" json:"address_created_at"`
	AddressUpdatedAt time.Time `avro:"address_updated_at" redis:"address_updated_at" json:"address_updated_at"`
}

type RedisFarm struct {
	ID               string  `redis:"id" json:"id"`
	FarmerID         string  `redis:"farmer_id" json:"farmer_id"`
	FarmName         string  `redis:"farm_name" json:"farm_name"`
	FarmType         string  `redis:"farm_type" json:"farm_type"`
	FarmSize         float64 `redis:"farm_size" json:"farm_size"`
	PhotoURL         *string `redis:"photo_url" json:"photo_url"`
	FarmStatus       string  `redis:"farm_status" json:"farm_status"`
	Description      *string `redis:"description" json:"description"`
	CreatedAt        int64   `redis:"created_at" json:"created_at"`
	UpdatedAt        int64   `redis:"updated_at" json:"updated_at"`
	AddressID        string  `redis:"address_id" json:"address_id"`
	Street           string  `redis:"street" json:"street"`
	Village          string  `redis:"village" json:"village"`
	SubDistrict      string  `redis:"sub_district" json:"sub_district"`
	City             string  `redis:"city" json:"city"`
	Province         string  `redis:"province" json:"province"`
	PostalCode       string  `redis:"postal_code" json:"postal_code"`
	AddressCreatedAt int64   `redis:"address_created_at" json:"address_created_at"`
	AddressUpdatedAt int64   `redis:"address_updated_at" json:"address_updated_at"`
}

func (rf RedisFarm) RedisFarmToFarmWithAddress() FarmWithAddress {
	var description string
	if rf.Description != nil {
		description = *rf.Description
	}

	return FarmWithAddress{
		Farm: Farm{
			ID:          rf.ID,
			FarmerID:    rf.FarmerID,
			FarmName:    rf.FarmName,
			FarmType:    rf.FarmType,
			FarmSize:    rf.FarmSize,
			FarmStatus:  rf.FarmStatus,
			Description: description,
			AddressesID: rf.AddressID,
			CreatedAt:   time.Unix(rf.CreatedAt, 0),
			UpdatedAt:   time.Unix(rf.UpdatedAt, 0),
		},

		FarmAddress: FarmAddress{
			ID:          rf.AddressID,
			Street:      rf.Street,
			Village:     rf.Village,
			SubDistrict: rf.SubDistrict,
			City:        rf.City,
			Province:    rf.Province,
			PostalCode:  rf.PostalCode,
			CreatedAt:   time.Unix(rf.AddressCreatedAt, 0),
			UpdatedAt:   time.Unix(rf.AddressUpdatedAt, 0),
		},
	}
}
