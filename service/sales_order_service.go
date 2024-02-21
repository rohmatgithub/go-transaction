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

func InsertSalesOrder(request dto.SalesOrderRequest, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	orderDate := time.Now()
	var err error

	validated := request.ValidateInsert(ctxModel)
	if validated != nil {
		out.Status.Detail = validated
		errMdl = model.GenerateFailedValidate()
		return
	}

	if len(request.ListItem) == 0 {
		errMdl = model.GenerateEmptyFieldError(constanta.Item)
		return
	}
	if request.OrderDate != "" {
		orderDate, err = time.Parse("2006-01-02", request.OrderDate)
		if err != nil {
			errMdl = model.GenerateFormatFieldError(constanta.OrderDate)
			return
		}
	}

	if request.OrderNumber == "" {
		str := common.GenerateRandomString(6)
		request.OrderNumber = "SO" + orderDate.Format("060102") + strings.ToUpper(str)
	} else {
		var resultDB *entity.SalesOrderEntity
		resultDB, errMdl = repository.FindSalesOrder(common.GormDB, &entity.SalesOrderEntity{
			OrderNumber: sql.NullString{String: request.OrderNumber, Valid: true},
		})
		if errMdl.Error != nil {
			return
		}
		if resultDB.ID > 0 {
			errMdl = model.GenerateHasUsedDataError(constanta.OrderNumber)
			return
		}
	}

	timeNow := time.Now()
	entitySo := entity.SalesOrderEntity{
		CompanyID:        sql.NullInt64{Int64: ctxModel.AuthAccessTokenModel.CompanyID, Valid: true},
		CustomerID:       sql.NullInt64{Int64: request.CustomerID, Valid: true},
		OrderNumber:      sql.NullString{String: request.OrderNumber, Valid: true},
		OrderDate:        sql.NullTime{Time: orderDate, Valid: true},
		TotalGrossAmount: sql.NullFloat64{Float64: request.TotalGrossAmount, Valid: true},
		TotalNetAmount:   sql.NullFloat64{Float64: request.TotalNetAmount, Valid: true},
		AbstractEntity: entity.AbstractEntity{
			CreatedBy: sql.NullInt64{Int64: ctxModel.AuthAccessTokenModel.ResourceUserID, Valid: true},
			UpdatedBy: sql.NullInt64{Int64: ctxModel.AuthAccessTokenModel.ResourceUserID, Valid: true},
			CreatedAt: sql.NullTime{Time: timeNow, Valid: true},
			UpdatedAt: sql.NullTime{Time: timeNow, Valid: true},
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

	errGorm := tx.Create(&entitySo)
	if errGorm.Error != nil {
		errMdl = model.GenerateInternalDBServerError(errGorm.Error)
		return
	}

	listSave := []map[string]model.UpsertModel{}

	for _, item := range request.ListItem {
		save := make(map[string]model.UpsertModel)
		save["order_id"] = model.UpsertModel{Value: entitySo.ID, PrimaryKey: true}
		save["product_id"] = model.UpsertModel{Value: item.ProductID, PrimaryKey: true}
		save["qty"] = model.UpsertModel{Value: item.Qty}
		save["selling_price"] = model.UpsertModel{Value: item.SellingPrice}
		save["line_gross_amount"] = model.UpsertModel{Value: item.LineGrossAmount}
		save["line_net_amount"] = model.UpsertModel{Value: item.LineNetAmount}
		save["created_by"] = model.UpsertModel{Value: ctxModel.AuthAccessTokenModel.ResourceUserID}
		save["updated_by"] = model.UpsertModel{Value: ctxModel.AuthAccessTokenModel.ResourceUserID}

		listSave = append(listSave, save)
	}

	errMdl = repository.InsertSalesOrderItem(tx, listSave)
	if errMdl.Error != nil {
		return
	}

	out.Status.Message = InsertI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}

func ListSalesOrder(dtoList dto.GetListRequest, searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {

	resultDB, errMdl := repository.ListSalesOrder(common.GormDB, dtoList, searchParam, ctxModel)
	if errMdl.Error != nil {
		return
	}

	var result []dto.ListSalesOrderResponse
	var listCustomerID []int64
	for _, v := range resultDB {
		temp := v.(entity.SalesOrderEntity)
		listCustomerID = append(listCustomerID, temp.CustomerID.Int64)
		result = append(result, dto.ListSalesOrderResponse{
			ID:               temp.ID,
			OrderNumber:      temp.OrderNumber.String,
			OrderDate:        temp.OrderDate.Time.Format("2006-01-02"),
			TotalGrossAmount: temp.TotalGrossAmount.Float64,
			TotalNetAmount:   temp.TotalNetAmount.Float64,
			CustomerID:       temp.CustomerID.Int64,
			IsGenerated:      temp.IsGenerated.Bool,
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

func CountListSalesOrder(searchParam []dto.SearchByParam, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {

	resultDB, errMdl := repository.CountListSalesOrder(common.GormDB, searchParam, ctxModel)
	if errMdl.Error != nil {
		return
	}

	out.Data = resultDB
	out.Status.Message = ListI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}

func GetDetailSalesOrder(id int64, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	// get sales order header
	resultDB, errMdl := repository.FindSalesOrder(common.GormDB, &entity.SalesOrderEntity{
		AbstractEntity: entity.AbstractEntity{
			ID: id,
		},
	})
	if errMdl.Error != nil {
		return
	}
	if resultDB.ID == 0 {
		errMdl = model.GenerateUnknownDataError(constanta.OrderNumber)
		return
	}
	// get sales order item
	listItem, errMdl := repository.GetListSalesOrderItem(common.GormDB, &entity.SalesOrderItemEntity{
		OrderID: sql.NullInt64{Int64: resultDB.ID},
	})
	if errMdl.Error != nil {
		return
	}
	var listItemResponse []dto.ListSalesOrderItemResponse
	var listProductID []int64
	for _, v := range listItem {
		temp := v.(entity.SalesOrderItemEntity)
		listItemResponse = append(listItemResponse, dto.ListSalesOrderItemResponse{
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
	resultFinal := dto.DetailSalesOrder{
		ListSalesOrderResponse: dto.ListSalesOrderResponse{
			ID:               resultDB.ID,
			UpdatedAt:        resultDB.UpdatedAt.Time,
			OrderNumber:      resultDB.OrderNumber.String,
			OrderDate:        resultDB.OrderDate.Time.Format("2006-01-02"),
			TotalGrossAmount: resultDB.TotalGrossAmount.Float64,
			TotalNetAmount:   resultDB.TotalNetAmount.Float64,
			CustomerID:       resultDB.CustomerID.Int64,
			CustomerCode:     customerCode,
			CustomerName:     customerName,
			IsGenerated:      resultDB.IsGenerated.Bool,
		},
		ListItem: listItemResponse,
	}

	out.Data = resultFinal
	out.Status.Message = ViewI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}

func UpdateSalesOrder(request dto.SalesOrderRequest, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	orderDate := time.Now()
	var err error

	validated, errMdl := request.ValidateUpdate(ctxModel)
	if errMdl.Error != nil {
		return
	}

	if validated != nil {
		out.Status.Detail = validated
		errMdl = model.GenerateFailedValidate()
		return
	}

	if len(request.ListItem) == 0 {
		errMdl = model.GenerateEmptyFieldError(constanta.Item)
		return
	}
	if request.OrderDate != "" {
		orderDate, err = time.Parse("2006-01-02", request.OrderDate)
		if err != nil {
			errMdl = model.GenerateFormatFieldError(constanta.OrderDate)
			return
		}
	}
	resultDB, errMdl := repository.FindSalesOrder(common.GormDB, &entity.SalesOrderEntity{
		AbstractEntity: entity.AbstractEntity{
			ID: request.ID,
		},
	})

	if resultDB.ID == 0 {
		errMdl = model.GenerateUnknownDataError(constanta.SalesOrder)
		return
	}
	if resultDB.OrderNumber.String != request.OrderNumber {
		errMdl = model.GenerateNotChangedDataError(constanta.OrderNumber)
		return
	}
	if resultDB.UpdatedAt.Time != request.UpdatedAt {
		errMdl = model.GenerateDataLockedError(constanta.SalesOrder)
		return
	}
	timeNow := time.Now()
	entitySo := entity.SalesOrderEntity{
		CompanyID:        sql.NullInt64{Int64: ctxModel.AuthAccessTokenModel.CompanyID, Valid: true},
		CustomerID:       sql.NullInt64{Int64: request.CustomerID, Valid: true},
		OrderNumber:      sql.NullString{String: request.OrderNumber, Valid: true},
		OrderDate:        sql.NullTime{Time: orderDate, Valid: true},
		TotalGrossAmount: sql.NullFloat64{Float64: request.TotalGrossAmount, Valid: true},
		TotalNetAmount:   sql.NullFloat64{Float64: request.TotalNetAmount, Valid: true},
		AbstractEntity: entity.AbstractEntity{
			ID:        resultDB.ID,
			CreatedBy: sql.NullInt64{Int64: resultDB.CreatedBy.Int64, Valid: true},
			UpdatedBy: sql.NullInt64{Int64: ctxModel.AuthAccessTokenModel.ResourceUserID, Valid: true},
			CreatedAt: sql.NullTime{Time: resultDB.CreatedAt.Time, Valid: true},
			UpdatedAt: sql.NullTime{Time: timeNow, Valid: true},
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

	errGorm := tx.Save(&entitySo)
	if errGorm.Error != nil {
		errMdl = model.GenerateInternalDBServerError(errGorm.Error)
		return
	}

	// delete all item
	errMdl = repository.DeleteSalesOrderItemByOrderID(tx, request.ID)
	if errMdl.Error != nil {
		return
	}

	listSave := []map[string]model.UpsertModel{}

	for _, item := range request.ListItem {
		save := make(map[string]model.UpsertModel)
		save["order_id"] = model.UpsertModel{Value: entitySo.ID, PrimaryKey: true}
		save["product_id"] = model.UpsertModel{Value: item.ProductID, PrimaryKey: true}
		save["qty"] = model.UpsertModel{Value: item.Qty}
		save["selling_price"] = model.UpsertModel{Value: item.SellingPrice}
		save["line_gross_amount"] = model.UpsertModel{Value: item.LineGrossAmount}
		save["line_net_amount"] = model.UpsertModel{Value: item.LineNetAmount}
		save["created_by"] = model.UpsertModel{Value: ctxModel.AuthAccessTokenModel.ResourceUserID}
		save["updated_by"] = model.UpsertModel{Value: ctxModel.AuthAccessTokenModel.ResourceUserID}

		listSave = append(listSave, save)
	}

	errMdl = repository.InsertSalesOrderItem(tx, listSave)
	if errMdl.Error != nil {
		return
	}

	out.Status.Message = InsertI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}
