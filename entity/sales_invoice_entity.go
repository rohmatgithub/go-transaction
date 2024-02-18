package entity

import "database/sql"

type SalesInvoice struct {
	AbstractEntity
	CompanyID        sql.NullInt64
	InvoiceNumber    sql.NullString
	InvoiceDate      sql.NullTime
	CustomerID       sql.NullInt64
	OrderID          sql.NullInt64
	TotalGrossAmount sql.NullFloat64
	TotalNetAmount   sql.NullFloat64
}

func (SalesInvoice) TableName() string {
	return "sales_invoice"
}
