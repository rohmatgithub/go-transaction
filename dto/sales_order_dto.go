package dto

import (
	"go-transaction/common"
	"go-transaction/model"
)

type SalesOrderRequest struct {
	OrderNumber      string                  `json:"order_number"`
	OrderDate        string                  `json:"order_date" validate:"required"`
	CustomerID       int64                   `json:"customer_id" validate:"required"`
	TotalGrossAmount float64                 `json:"total_gross_amount" validate:"required"`
	TotalNetAmount   float64                 `json:"total_net_amount" validate:"required"`
	ListItem         []SalesOrderItemRequest `json:"list_item" validate:"required"`
	AbstractDto
}

type SalesOrderItemRequest struct {
	ProductID       int64   `json:"product_id" validate:"required"`
	Qty             int64   `json:"qty" validate:"required"`
	SellingPrice    float64 `json:"selling_price" validate:"required"`
	LineGrossAmount float64 `json:"line_gross_amount" validate:"required"`
	LineNetAmount   float64 `json:"line_net_amount" validate:"required"`
}

type ListSalesOrderResponse struct {
	ID               int64   `json:"id"`
	OrderNumber      string  `json:"order_number"`
	OrderDate        string  `json:"order_date"`
	TotalGrossAmount float64 `json:"total_gross_amount"`
	TotalNetAmount   float64 `json:"total_net_amount"`
	CustomerID       int64   `json:"customer_id"`
	CustomerCode     string  `json:"customer_code"`
	CustomerName     string  `json:"customer_name"`
}

func (c *SalesOrderRequest) ValidateInsert(contextModel *common.ContextModel) map[string]string {
	return common.Validation.ValidationAll(*c, contextModel)
}

func (c *SalesOrderRequest) ValidateUpdate(contextModel *common.ContextModel) (resultMap map[string]string, errMdl model.ErrorModel) {
	resultMap = common.Validation.ValidationAll(*c, contextModel)

	errMdl = c.ValidateUpdateGeneral()

	return
}
