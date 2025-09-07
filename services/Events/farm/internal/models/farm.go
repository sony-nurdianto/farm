package models

type Farm struct {
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

type FarmAddress struct {
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

type FarmWithAddress struct {
	Farm
	Street           string `avro:"street" redis:"street" json:"street"`
	Village          string `avro:"village" redis:"village" json:"village"`
	SubDistrict      string `avro:"sub_district" redis:"sub_district" json:"sub_district"`
	City             string `avro:"city" redis:"city" json:"city"`
	Province         string `avro:"province" redis:"province" json:"province"`
	PostalCode       string `avro:"postal_code" redis:"postal_code" json:"postal_code"`
	AddressCreatedAt string `avro:"address_created_at" redis:"address_created_at" json:"address_created_at"`
	AddressUpdatedAt string `avro:"address_updated_at" redis:"address_updated_at" json:"address_updated_at"`
}
