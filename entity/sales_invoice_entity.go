package entity

import "database/sql"

type SalesInvoiceEntity struct {
	AbstractEntity
	CompanyID        sql.NullInt64
	InvoiceNumber    sql.NullString
	InvoiceDate      sql.NullTime
	CustomerID       sql.NullInt64
	OrderID          sql.NullInt64
	TotalGrossAmount sql.NullFloat64
	TotalNetAmount   sql.NullFloat64
}

func (SalesInvoiceEntity) TableName() string {
	return "sales_invoice"
}

type ReportSalesInvoiceEntity struct {
	CustomerID      sql.NullInt64
	ProductID       sql.NullInt64
	Qty             sql.NullInt64
	LineGrossAmount sql.NullFloat64
	LineNetAmount   sql.NullFloat64
}
