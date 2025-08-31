package models

import (
	"strings"
	"time"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type SortOrder int

const (
	SortOrderUnknown SortOrder = iota
	SortOrderAsc
	SorOrderDesc
)

func (so SortOrder) String() string {
	switch so {
	case SortOrderAsc:
		return "ASC"
	case SorOrderDesc:
		return "DESC"
	default:
		return "Unknown"
	}
}

func (so SortOrder) ProtoSortOrder() pbgen.SortOrder {
	switch so {
	case SortOrderAsc:
		return pbgen.SortOrder_SortOrder_ASC
	case SorOrderDesc:
		return pbgen.SortOrder_SortOrder_DESC
	default:
		return pbgen.SortOrder_SortOrder_UKNOWN
	}
}

func (so SortOrder) IntToSortOrder(val int) SortOrder {
	switch val {
	case 1:
		return SortOrderAsc
	case 2:
		return SorOrderDesc
	default:
		return SortOrderUnknown
	}
}

func (so SortOrder) StringToSortOrder(val string) SortOrder {
	val = strings.ToLower(val)

	switch val {
	case "asc", "ascending":
		return SortOrderAsc
	case "desc", "descending":
		return SorOrderDesc
	default:
		return SortOrderUnknown
	}
}

func (so SortOrder) ProtoSortOrderFactory(val pbgen.SortOrder) SortOrder {
	switch val {
	case pbgen.SortOrder_SortOrder_ASC:
		return SortOrderAsc
	case pbgen.SortOrder_SortOrder_DESC:
		return SorOrderDesc
	default:
		return SortOrderUnknown
	}
}

type FarmAddress struct {
	ID          string `json:"id"`
	Street      string `json:"street"`
	Village     string `json:"village"`
	SubDistrict string `json:"sub_district"`
	City        string `json:"city"`
	Province    string `json:"province"`
	PostalCode  string `json:"postal_code"`
}

type Farm struct {
	ID          string  `json:"id"`
	FarmerID    string  `json:"farmer_id"`
	FarmName    string  `json:"farm_name"`
	FarmType    string  `json:"farm_type"`
	FarmSize    float64 `json:"farm_size"`
	FarmStatus  string  `json:"farm_status"`
	Description string  `json:"description"`
	Addresses   FarmAddress
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetFarmsRequest struct {
	SearchName string    `json:"search_name"`
	SortOrder  SortOrder `json:"sort_order"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

type GetFarmsResponse struct {
	Data  []Farm
	Total int `json:"total"`
}
