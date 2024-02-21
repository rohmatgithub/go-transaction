package repository

import (
	"database/sql"
	"go-transaction/common"
	"go-transaction/dto"
	"go-transaction/entity"
	"go-transaction/model"
	"strconv"

	"gorm.io/gorm"
)

func InsertSalesOrderItem(tx *gorm.DB, data []map[string]model.UpsertModel) (errMdl model.ErrorModel) {
	queryString, queryParam := getQueryUpsertMultiValues("sales_order_item", data)

	result := tx.Exec(queryString, queryParam...)
	if result.Error != nil {
		return model.GenerateInternalDBServerError(result.Error)
	}
	return
}

func ListSalesOrder(db *gorm.DB, dtoList dto.GetListRequest, searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (result []interface{}, errMdl model.ErrorModel) {
	// for i := 0; i < len(searchParam); i++ {
	// 	switch searchParam[i].SearchKey {
	// 	case "name":
	// 		searchParam[i].SearchKey = "c." + searchParam[i].SearchKey
	// 	}
	// }
	searchParam = append(searchParam, dto.SearchByParam{
		SearchKey:      "so.company_id",
		SearchOperator: "eq",
		SearchValue:    strconv.Itoa(int(ctxModel.AuthAccessTokenModel.CompanyID)),
	})
	dtoList.OrderBy = "so." + dtoList.OrderBy
	query := "SELECT so.id, so.order_number, so.order_date, " +
		"so.total_gross_amount, so.total_net_amount, so.customer_id, " +
		"so.is_generated " +
		"FROM sales_order so "

	return GetListDataDefault(db, query, nil, dtoList, searchParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp entity.SalesOrderEntity
			err := rows.Scan(&temp.ID, &temp.OrderNumber, &temp.OrderDate,
				&temp.TotalGrossAmount, &temp.TotalNetAmount, &temp.CustomerID,
				&temp.IsGenerated)
			return temp, err
		})

}
func CountListSalesOrder(db *gorm.DB, searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (result int64, errMdl model.ErrorModel) {
	searchParam = append(searchParam, dto.SearchByParam{
		SearchKey:      "company_id",
		SearchOperator: "eq",
		SearchValue:    strconv.Itoa(int(ctxModel.AuthAccessTokenModel.CompanyID)),
	},
	//  dto.SearchByParam{
	// 	SearchKey:      "c.branch_id",
	// 	SearchOperator: "eq",
	// 	SearchValue:    strconv.Itoa(int(ctxModel.AuthAccessTokenModel.BranchID)),
	// }
	)
	query := "SELECT COUNT(0) FROM sales_order "

	return GetCountDataDefault(db, query, nil, searchParam)

}

func FindSalesOrder(db *gorm.DB, entity *entity.SalesOrderEntity) (result *entity.SalesOrderEntity, errMdl model.ErrorModel) {
	gormResult := db.Where(entity).Find(&result)
	if gormResult.Error != nil {
		errMdl = model.GenerateInternalDBServerError(gormResult.Error)
		return
	}

	return
}

func GetListSalesOrderItem(db *gorm.DB, e *entity.SalesOrderItemEntity) (result []interface{}, errMdl model.ErrorModel) {
	query := "SELECT order_id, product_id, qty, " +
		"selling_price, line_gross_amount, line_net_amount " +
		"FROM sales_order_item " +
		"WHERE order_id = $1 "

	queryParam := []interface{}{e.OrderID.Int64}
	return ExecuteQuery(db, query, queryParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp entity.SalesOrderItemEntity
			err := rows.Scan(&temp.OrderID, &temp.ProductID, &temp.Qty,
				&temp.SellingPrice, &temp.LineGrossAmount, &temp.LineNetAmount)
			return temp, err
		})

}

func DeleteSalesOrderItemByOrderID(tx *gorm.DB, orderID int64) (errMdl model.ErrorModel) {
	queryString := "DELETE FROM sales_order_item WHERE order_id = $1"

	result := tx.Exec(queryString, orderID)
	if result.Error != nil {
		return model.GenerateInternalDBServerError(result.Error)
	}
	return
}
