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
		save["qty"] = model.UpsertModel{Value: entitySo.ID}
		save["selling_price"] = model.UpsertModel{Value: entitySo.ID}
		save["line_gross_amount"] = model.UpsertModel{Value: entitySo.ID}
		save["line_net_amount"] = model.UpsertModel{Value: entitySo.ID}
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
