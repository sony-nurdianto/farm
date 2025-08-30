package models

type UpdateFarmAddr struct {
	ID          string `json:"id"`
	Street      string `json:"street"`
	Village     string `json:"village"`
	SubDistrict string `json:"sub_district"`
	City        string `json:"city"`
	Province    string `json:"province"`
	PostalCode  string `json:"postal_code"`
}

type UpdateFarm struct {
	ID          string  `json:"id"`
	FarmName    string  `json:"farm_name"`
	FarmType    string  `json:"farm_type"`
	FarmSize    float64 `json:"farm_size"`
	FarmStatus  string  `json:"farm_status"`
	Description string  `json:"description"`
}

type UpdateFarmWithAddr struct {
	Farm    *UpdateFarm
	Address *UpdateFarmAddr
}
