package models

type UpdateFarm struct {
	ID          string
	FarmName    string
	FarmType    string
	FarmSize    float64
	FarmStatus  string
	Description string
}

type UpdateFarmAddress struct {
	ID          string
	Street      string
	Village     string
	SubDistrict string
	City        string
	Province    string
	PostalCode  string
}
