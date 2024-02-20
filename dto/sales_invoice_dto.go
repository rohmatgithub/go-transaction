package dto

type ListSalesInvoiceResponse struct {
	ID               int64   `json:"id"`
	InvoiceNumber    string  `json:"invoice_number"`
	InvoiceDate      string  `json:"invoice_date"`
	TotalGrossAmount float64 `json:"total_gross_amount"`
	TotalNetAmount   float64 `json:"total_net_amount"`
	CustomerID       int64   `json:"customer_id"`
	CustomerCode     string  `json:"customer_code"`
	CustomerName     string  `json:"customer_name"`
}

type DetailSalesInvoice struct {
	ListSalesInvoiceResponse
	ListItem []ListSalesInvoiceItemResponse `json:"list_items"`
}
type ListSalesInvoiceItemResponse struct {
	ProductID       int64   `json:"product_id"`
	ProductCode     string  `json:"product_code"`
	ProductName     string  `json:"product_name"`
	Qty             int64   `json:"qty"`
	SellingPrice    float64 `json:"selling_price"`
	LineGrossAmount float64 `json:"line_gross_amount"`
	LineNetAmount   float64 `json:"line_net_amount"`
}
