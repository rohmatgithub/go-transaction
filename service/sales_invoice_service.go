package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-transaction/common"
	"go-transaction/config"
	"go-transaction/constanta"
	"go-transaction/dto"
	"go-transaction/entity"
	"go-transaction/model"
	"go-transaction/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func CreateSalesInvoice(salesOrderID int64, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	invoiceDate := time.Now()

	entitySo, errMdl := repository.FindSalesOrder(common.GormDB, &entity.SalesOrderEntity{
		AbstractEntity: entity.AbstractEntity{
			ID: salesOrderID,
		},
	})
	if errMdl.Error != nil {
		return
	}

	if entitySo.ID == 0 {
		errMdl = model.GenerateUnknownDataError(constanta.SalesOrder)
		return
	}

	if entitySo.IsGenerated.Bool {
		errMdl = model.GenerateDataLockedError(constanta.SalesOrder)
		return
	}
	invoiceNumber := "SI" + invoiceDate.Format("060102") + strings.ToUpper(common.GenerateRandomString(6))

	listItem, errMdl := repository.GetListSalesOrderItem(common.GormDB, &entity.SalesOrderItemEntity{
		OrderID: sql.NullInt64{Int64: salesOrderID},
	})
	if errMdl.Error != nil {
		return
	}

	entitySi := &entity.SalesInvoiceEntity{
		CompanyID:        sql.NullInt64{Int64: entitySo.CompanyID.Int64, Valid: true},
		CustomerID:       sql.NullInt64{Int64: entitySo.CustomerID.Int64, Valid: true},
		OrderID:          sql.NullInt64{Int64: salesOrderID, Valid: true},
		InvoiceNumber:    sql.NullString{String: invoiceNumber, Valid: true},
		InvoiceDate:      sql.NullTime{Time: invoiceDate, Valid: true},
		TotalGrossAmount: sql.NullFloat64{Float64: entitySo.TotalGrossAmount.Float64, Valid: true},
		TotalNetAmount:   sql.NullFloat64{Float64: entitySo.TotalNetAmount.Float64, Valid: true},
		AbstractEntity: entity.AbstractEntity{
			CreatedBy: sql.NullInt64{Int64: ctxModel.AuthAccessTokenModel.ResourceUserID, Valid: true},
			UpdatedBy: sql.NullInt64{Int64: ctxModel.AuthAccessTokenModel.ResourceUserID, Valid: true},
			CreatedAt: sql.NullTime{Time: invoiceDate, Valid: true},
			UpdatedAt: sql.NullTime{Time: invoiceDate, Valid: true},
		},
	}

	// insert token to db
	tx := common.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil || errMdl.Error != nil {
			tx.Rollback()
		} else {
			// save to redis
			err := tx.Commit().Error
			if err != nil {
				errMdl = model.GenerateInternalDBServerError(err)
				return
			}
		}
	}()

	errGorm := tx.Create(&entitySi)
	if errGorm.Error != nil {
		errMdl = model.GenerateInternalDBServerError(errGorm.Error)
		return
	}

	// update sales order
	entitySo.IsGenerated = sql.NullBool{Bool: true, Valid: true}
	errGorm = tx.Save(&entitySo)
	if errGorm.Error != nil {
		errMdl = model.GenerateInternalDBServerError(errGorm.Error)
		return
	}

	listSave := []map[string]model.UpsertModel{}
	for _, v := range listItem {
		temp := v.(entity.SalesOrderItemEntity)
		save := make(map[string]model.UpsertModel)
		save["invoice_id"] = model.UpsertModel{Value: entitySi.ID, PrimaryKey: true}
		save["product_id"] = model.UpsertModel{Value: temp.ProductID, PrimaryKey: true}
		save["qty"] = model.UpsertModel{Value: temp.Qty}
		save["selling_price"] = model.UpsertModel{Value: temp.SellingPrice}
		save["line_gross_amount"] = model.UpsertModel{Value: temp.LineGrossAmount}
		save["line_net_amount"] = model.UpsertModel{Value: temp.LineNetAmount}
		save["created_by"] = model.UpsertModel{Value: ctxModel.AuthAccessTokenModel.ResourceUserID}
		save["updated_by"] = model.UpsertModel{Value: ctxModel.AuthAccessTokenModel.ResourceUserID}

		listSave = append(listSave, save)
	}

	errMdl = repository.CreateInvoiceItem(tx, listSave)
	if errMdl.Error != nil {
		return
	}

	out.Status.Message = InsertI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}

func ListSalesInvoice(dtoList dto.GetListRequest, searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {

	resultDB, errMdl := repository.ListSalesInvoice(common.GormDB, dtoList, searchParam, ctxModel)
	if errMdl.Error != nil {
		return
	}

	var result []dto.ListSalesInvoiceResponse
	var listCustomerID []int64
	for _, v := range resultDB {
		temp := v.(entity.SalesInvoiceEntity)
		listCustomerID = append(listCustomerID, temp.CustomerID.Int64)
		result = append(result, dto.ListSalesInvoiceResponse{
			ID:               temp.ID,
			InvoiceNumber:    temp.InvoiceNumber.String,
			InvoiceDate:      temp.InvoiceDate.Time.Format("2006-01-02"),
			TotalGrossAmount: temp.TotalGrossAmount.Float64,
			TotalNetAmount:   temp.TotalNetAmount.Float64,
			CustomerID:       temp.CustomerID.Int64,
			// CustomerCode:     temp.CustomerCode.String,
			// CustomerName:     temp.CustomerName.String,
		})
	}

	// get customer to master data service
	uri := config.ApplicationConfiguration.GetUriResouce().MasterData + "/v1/master/customer?page=1&list_id=" + strings.Join(strings.Fields(fmt.Sprint(listCustomerID)), ",")
	response, errMdl := HitToResourceOther(uri, "GET", ctxModel)
	if errMdl.Error != nil {
		return
	}
	var customer dto.ListCustomerResponse
	err := json.Unmarshal(response, &customer)
	if err != nil {
		log.Error("Error unmarshalling JSON:", err)
		errMdl = model.GenerateUnknownError(err)
		return
	}

	for i := 0; i < len(result); i++ {
		for _, v := range customer.Payload.Data {
			if result[i].CustomerID == v.ID {
				result[i].CustomerCode = v.Code
				result[i].CustomerName = v.Name
				break
			}
		}
	}

	out.Data = result
	out.Status.Message = ListI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}

func CountListSalesInvoice(searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {

	resultDB, errMdl := repository.CountListSalesInvoice(common.GormDB, searchParam, ctxModel)
	if errMdl.Error != nil {
		return
	}

	out.Data = resultDB
	out.Status.Message = ListI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}

func GetDetailSalesInvoice(id int64, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	// get sales Invoice header
	resultDB, errMdl := repository.FindSalesInvoice(common.GormDB, &entity.SalesInvoiceEntity{
		AbstractEntity: entity.AbstractEntity{
			ID: id,
		},
	})
	if errMdl.Error != nil {
		return
	}
	if resultDB.ID == 0 {
		errMdl = model.GenerateUnknownDataError(constanta.InvoiceNumber)
		return
	}
	// get sales Invoice item
	listItem, errMdl := repository.GetListSalesInvoiceItem(common.GormDB, &entity.SalesInvoiceItemEntity{
		InvoiceID: sql.NullInt64{Int64: resultDB.ID},
	})
	if errMdl.Error != nil {
		return
	}
	var listItemResponse []dto.ListSalesInvoiceItemResponse
	var listProductID []int64
	for _, v := range listItem {
		temp := v.(entity.SalesInvoiceItemEntity)
		listItemResponse = append(listItemResponse, dto.ListSalesInvoiceItemResponse{
			ProductID:       temp.ProductID.Int64,
			Qty:             temp.Qty.Int64,
			SellingPrice:    temp.SellingPrice.Float64,
			LineGrossAmount: temp.LineGrossAmount.Float64,
			LineNetAmount:   temp.LineNetAmount.Float64,
		})
		listProductID = append(listProductID, temp.ProductID.Int64)
	}

	// get customer to master data service
	uri := config.ApplicationConfiguration.GetUriResouce().MasterData + "/v1/master/customer?page=1&list_id=" + strconv.FormatInt(resultDB.CustomerID.Int64, 10)
	response, errMdl := HitToResourceOther(uri, "GET", ctxModel)
	if errMdl.Error != nil {
		return
	}
	var customer dto.ListCustomerResponse
	err := json.Unmarshal(response, &customer)
	if err != nil {
		log.Error("Error unmarshalling JSON:", err)
		errMdl = model.GenerateUnknownError(err)
		return
	}
	// get product to master data
	uri = config.ApplicationConfiguration.GetUriResouce().MasterData + "/v1/master/product?page=1&list_id=" + strings.Join(strings.Fields(fmt.Sprint(listProductID)), ",")
	response, errMdl = HitToResourceOther(uri, "GET", ctxModel)
	if errMdl.Error != nil {
		return
	}
	var listProduct dto.ListProductResponse
	err = json.Unmarshal(response, &listProduct)
	if err != nil {
		log.Error("Error unmarshalling JSON:", err)
		errMdl = model.GenerateUnknownError(err)
		return
	}

	for _, product := range listProduct.Payload.Data {
		for i := 0; i < len(listItemResponse); i++ {
			if listItemResponse[i].ProductID == product.ID {
				listItemResponse[i].ProductCode = product.Code
				listItemResponse[i].ProductName = product.Name
				break
			}
		}
	}

	var customerCode, customerName string
	if len(customer.Payload.Data) > 0 {
		customerCode = customer.Payload.Data[0].Code
		customerName = customer.Payload.Data[0].Name
	}
	resultFinal := dto.DetailSalesInvoice{
		ListSalesInvoiceResponse: dto.ListSalesInvoiceResponse{
			ID:               resultDB.ID,
			InvoiceNumber:    resultDB.InvoiceNumber.String,
			InvoiceDate:      resultDB.InvoiceDate.Time.Format("2006-01-02"),
			TotalGrossAmount: resultDB.TotalGrossAmount.Float64,
			TotalNetAmount:   resultDB.TotalNetAmount.Float64,
			CustomerID:       resultDB.CustomerID.Int64,
			CustomerCode:     customerCode,
			CustomerName:     customerName,
		},
		ListItem: listItemResponse,
	}

	out.Data = resultFinal
	out.Status.Message = ViewI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}
