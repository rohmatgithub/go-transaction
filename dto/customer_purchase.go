package dto

import "time"

type CustomerPurchaseRequest struct {
	CustomerID  int64  `json:"customer_id"`
	FromDateStr string `json:"from_date"`
	ThruDateStr string `json:"thru_date"`
	FromDate    time.Time
	ThurDate    time.Time
}

type CustomerPurchaseResponse struct {
	CustomerID      int64   `json:"customer_id"`
	ProductID       int64   `json:"product_id"`
	Qty             int64   `json:"qty"`
	LineGrossAmount float64 `json:"line_gross_amount"`
	LineNetAmount   float64 `json:"line_net_amount"`
}
