package models

type Stock struct {
	StockID int64 `json:"stock"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Company string `json:"company"`
}