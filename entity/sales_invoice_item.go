package entity

import "database/sql"

type SalesInvoiceItemEntity struct {
	AbstractEntity
	InvoiceID       sql.NullInt64
	ProductID       sql.NullInt64
	Qty             sql.NullInt64
	SellingPrice    sql.NullFloat64
	LineGrossAmount sql.NullFloat64
	LineNetAmount   sql.NullFloat64
}

func (SalesInvoiceItemEntity) TableName() string {
	return "sales_invoice_item"
}
