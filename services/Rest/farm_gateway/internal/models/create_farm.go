package models

type CreateFarmAddress struct {
	Street      string `json:"street"`
	Village     string `json:"village"`
	SubDistrict string `json:"sub_district"`
	City        string `json:"city"`
	Province    string `json:"province"`
	PostalCode  string `json:"postal_code"`
}

type CreateFarm struct {
	FarmerID    string  `json:"farmer_id"`
	FarmName    string  `json:"farm_name"`
	FarmType    string  `json:"farm_type"`
	FarmSize    float64 `json:"farm_size"`
	FarmStatus  string  `json:"farm_status"`
	Description string  `json:"description"`
	Address     CreateFarmAddress
}
