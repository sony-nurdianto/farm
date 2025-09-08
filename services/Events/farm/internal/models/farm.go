package models

import (
	"strconv"
	"time"
)

type AvrFarm struct {
	ID          string  `avro:"id" redis:"id" json:"id"`
	FarmerID    string  `avro:"farmer_id" redis:"farmer_id" json:"farmer_id"`
	FarmName    string  `avro:"farm_name" redis:"farm_name" json:"farm_name"`
	FarmType    string  `avro:"farm_type" redis:"farm_type" json:"farm_type"`
	FarmSize    float64 `avro:"farm_size" redis:"farm_size" json:"farm_size"`
	PhotoURL    *string `avro:"photo_url" redis:"photo_url" json:"photo_url"`
	FarmStatus  string  `avro:"farm_status" redis:"farm_status" json:"farm_status"`
	Description *string `avro:"description" redis:"description" json:"description"`
	CreatedAt   string  `avro:"created_at" redis:"created_at" json:"created_at"`
	UpdatedAt   string  `avro:"updated_at" redis:"updated_at" json:"updated_at"`
	AddressID   string  `avro:"address_id" redis:"address_id" json:"address_id"`
}

type AvrFarmAddress struct {
	ID          string `avro:"id" redis:"address_id" json:"id"`
	Street      string `avro:"street" redis:"street" json:"street"`
	Village     string `avro:"village" redis:"village" json:"village"`
	SubDistrict string `avro:"sub_district" redis:"sub_district" json:"sub_district"`
	City        string `avro:"city" redis:"city" json:"city"`
	Province    string `avro:"province" redis:"province" json:"province"`
	PostalCode  string `avro:"postal_code" redis:"postal_code" json:"postal_code"`
	CreatedAt   string `avro:"created_at" redis:"address_created_at" json:"created_at"`
	UpdatedAt   string `avro:"updated_at" redis:"address_updated_at" json:"updated_at"`
}

func stringToInt64(str string) int64 {
	if str == "" {
		return 0
	}

	if val, err := strconv.ParseInt(str, 10, 64); err == nil {
		return val
	}

	if t, err := time.Parse(time.RFC3339Nano, str); err == nil {
		return t.Unix()
	}

	return time.Now().Unix()
}

func (avrFarm *AvrFarm) ToFarm() Farm {
	updatedAt := stringToInt64(avrFarm.UpdatedAt)
	createdAt := stringToInt64(avrFarm.AddressID)

	return Farm{
		ID:          avrFarm.ID,
		FarmerID:    avrFarm.FarmerID,
		FarmName:    avrFarm.FarmName,
		FarmType:    avrFarm.FarmType,
		FarmSize:    avrFarm.FarmSize,
		PhotoURL:    avrFarm.PhotoURL,
		FarmStatus:  avrFarm.FarmStatus,
		Description: avrFarm.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		AddressID:   avrFarm.AddressID,
	}
}

func (avrAddr *AvrFarmAddress) ToFarmAddress() FarmAddress {
	createdAt := stringToInt64(avrAddr.CreatedAt)
	updatedAt := stringToInt64(avrAddr.UpdatedAt)

	return FarmAddress{
		ID:          avrAddr.ID,
		Street:      avrAddr.Street,
		Village:     avrAddr.Village,
		SubDistrict: avrAddr.SubDistrict,
		City:        avrAddr.City,
		Province:    avrAddr.Province,
		PostalCode:  avrAddr.PostalCode,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type Farm struct {
	ID          string  `redis:"id" json:"id"`
	FarmerID    string  `redis:"farmer_id" json:"farmer_id"`
	FarmName    string  `redis:"farm_name" json:"farm_name"`
	FarmType    string  `redis:"farm_type" json:"farm_type"`
	FarmSize    float64 `redis:"farm_size" json:"farm_size"`
	PhotoURL    *string `redis:"photo_url" json:"photo_url"`
	FarmStatus  string  `redis:"farm_status" json:"farm_status"`
	Description *string `redis:"description" json:"description"`
	CreatedAt   int64   `redis:"created_at" json:"created_at"`
	UpdatedAt   int64   `redis:"updated_at" json:"updated_at"`
	AddressID   string  `redis:"address_id" json:"address_id"`
}

type FarmAddress struct {
	ID          string `redis:"address_id" json:"id"`
	Street      string `redis:"street" json:"street"`
	Village     string `redis:"village" json:"village"`
	SubDistrict string `redis:"sub_district" json:"sub_district"`
	City        string `redis:"city" json:"city"`
	Province    string `redis:"province" json:"province"`
	PostalCode  string `redis:"postal_code" json:"postal_code"`
	CreatedAt   int64  `redis:"address_created_at" json:"created_at"`
	UpdatedAt   int64  `redis:"address_updated_at" json:"updated_at"`
}

type InsertFarm struct {
	ID               string  `avro:"id" redis:"id" json:"id"`
	FarmerID         string  `avro:"farmer_id" redis:"farmer_id" json:"farmer_id"`
	FarmName         string  `avro:"farm_name" redis:"farm_name" json:"farm_name"`
	FarmType         string  `avro:"farm_type" redis:"farm_type" json:"farm_type"`
	FarmSize         float64 `avro:"farm_size" redis:"farm_size" json:"farm_size"`
	PhotoURL         *string `avro:"photo_url" redis:"photo_url" json:"photo_url"`
	FarmStatus       string  `avro:"farm_status" redis:"farm_status" json:"farm_status"`
	Description      *string `avro:"description" redis:"description" json:"description"`
	CreatedAt        int64   `avro:"created_at" redis:"created_at" json:"created_at"`
	UpdatedAt        int64   `avro:"updated_at" redis:"updated_at" json:"updated_at"`
	AddressID        string  `avro:"address_id" redis:"address_id" json:"address_id"`
	Street           string  `avro:"street" redis:"street" json:"street"`
	Village          string  `avro:"village" redis:"village" json:"village"`
	SubDistrict      string  `avro:"sub_district" redis:"sub_district" json:"sub_district"`
	City             string  `avro:"city" redis:"city" json:"city"`
	Province         string  `avro:"province" redis:"province" json:"province"`
	PostalCode       string  `avro:"postal_code" redis:"postal_code" json:"postal_code"`
	AddressCreatedAt int64   `avro:"address_created_at" redis:"address_created_at" json:"address_created_at"`
	AddressUpdatedAt int64   `avro:"address_updated_at" redis:"address_updated_at" json:"address_updated_at"`
}
