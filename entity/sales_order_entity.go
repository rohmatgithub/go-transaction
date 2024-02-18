package entity

import "database/sql"

type SalesOrderEntity struct {
	AbstractEntity
	CompanyID        sql.NullInt64
	OrderNumber      sql.NullString
	OrderDate        sql.NullTime
	CustomerID       sql.NullInt64
	TotalGrossAmount sql.NullFloat64
	TotalNetAmount   sql.NullFloat64
}

func (SalesOrderEntity) TableName() string {
	return "sales_order"
}
