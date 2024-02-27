package repository

import (
	"database/sql"
	"go-transaction/common"
	"go-transaction/dto"
	"go-transaction/entity"
	"go-transaction/model"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func CreateInvoiceItem(tx *gorm.DB, data []map[string]model.UpsertModel) (errMdl model.ErrorModel) {
	queryString, queryParam := getQueryUpsertMultiValues("sales_invoice_item", data)

	result := tx.Exec(queryString, queryParam...)
	if result.Error != nil {
		return model.GenerateInternalDBServerError(result.Error)
	}
	return
}

func ListSalesInvoice(db *gorm.DB, dtoList dto.GetListRequest, searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (result []interface{}, errMdl model.ErrorModel) {
	// for i := 0; i < len(searchParam); i++ {
	// 	switch searchParam[i].SearchKey {
	// 	case "name":
	// 		searchParam[i].SearchKey = "c." + searchParam[i].SearchKey
	// 	}
	// }
	searchParam = append(searchParam, dto.SearchByParam{
		SearchKey:      "si.company_id",
		SearchOperator: "eq",
		SearchValue:    strconv.Itoa(int(ctxModel.AuthAccessTokenModel.CompanyID)),
	})
	dtoList.OrderBy = "si." + dtoList.OrderBy
	query := "SELECT si.id, si.invoice_number, si.invoice_date, " +
		"si.total_gross_amount, si.total_net_amount, si.customer_id " +
		"FROM sales_invoice si "

	return GetListDataDefault(db, query, nil, dtoList, searchParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp entity.SalesInvoiceEntity
			err := rows.Scan(&temp.ID, &temp.InvoiceNumber, &temp.InvoiceDate,
				&temp.TotalGrossAmount, &temp.TotalNetAmount, &temp.CustomerID)
			return temp, err
		})

}
func CountListSalesInvoice(db *gorm.DB, searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (result int64, errMdl model.ErrorModel) {
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
	query := "SELECT COUNT(0) FROM sales_invoice "

	return GetCountDataDefault(db, query, nil, searchParam)

}

func FindSalesInvoice(db *gorm.DB, entity *entity.SalesInvoiceEntity) (result *entity.SalesInvoiceEntity, errMdl model.ErrorModel) {
	gormResult := db.Where(entity).Find(&result)
	if gormResult.Error != nil {
		errMdl = model.GenerateInternalDBServerError(gormResult.Error)
		return
	}

	return
}

func GetListSalesInvoiceItem(db *gorm.DB, e *entity.SalesInvoiceItemEntity) (result []interface{}, errMdl model.ErrorModel) {
	query := "SELECT Invoice_id, product_id, qty, " +
		"selling_price, line_gross_amount, line_net_amount " +
		"FROM sales_invoice_item " +
		"WHERE invoice_id = $1 "

	queryParam := []interface{}{e.InvoiceID.Int64}
	return ExecuteQuery(db, query, queryParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp entity.SalesInvoiceItemEntity
			err := rows.Scan(&temp.InvoiceID, &temp.ProductID, &temp.Qty,
				&temp.SellingPrice, &temp.LineGrossAmount, &temp.LineNetAmount)
			return temp, err
		})

}

func GetReportCustomerPurchase(db *gorm.DB, companyID, customerID int64, fromDate time.Time, thurDate time.Time) (result []interface{}, errMdl model.ErrorModel) {
	query := " SELECT si.customer_id, sii.product_id, sum(sii.qty), sum(line_gross_amount), sum(line_net_amount) " +
		"FROM sales_invoice_item sii " +
		"LEFT JOIN sales_invoice si ON sii.invoice_id = si.id " +
		"WHERE si.company_id = $1 AND si.invoice_date BETWEEN $2 AND $3 "
	queryParam := []interface{}{companyID, fromDate, thurDate}
	if customerID > 0 {
		query += "AND si.customer_id = $4 "
		queryParam = append(queryParam, customerID)
	}

	query += " GROUP BY si.customer_id, sii.product_id  ORDER BY si.customer_id "

	return ExecuteQuery(db, query, queryParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp entity.ReportSalesInvoiceEntity
			err := rows.Scan(&temp.CustomerID, &temp.ProductID, &temp.Qty, &temp.LineGrossAmount, &temp.LineNetAmount)
			return temp, err
		})
}
