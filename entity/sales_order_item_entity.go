package entity

import "database/sql"

type SalesOrderItemEntity struct {
	AbstractEntity
	OrderID         sql.NullInt64
	ProductID       sql.NullInt64
	Qty             sql.NullInt64
	SellingPrice    sql.NullFloat64
	LineGrossAmount sql.NullFloat64
	LineNetAmount   sql.NullFloat64
}

func (SalesOrderItemEntity) TableName() string {
	return "sales_order_item"
}
