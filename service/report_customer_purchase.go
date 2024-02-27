package service

import (
	"go-transaction/common"
	"go-transaction/dto"
	"go-transaction/entity"
	"go-transaction/model"
	"go-transaction/repository"
	"time"
)

func ReportCustomerPurchase(request dto.CustomerPurchaseRequest, ctxModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	fromDate, err := time.Parse("2006-01-02", request.FromDateStr)
	if err != nil {
		errMdl = model.GenerateUnknownError(err)
		return
	}

	thruDate, err := time.Parse("2006-01-02", request.ThruDateStr)
	if err != nil {
		errMdl = model.GenerateUnknownError(err)
		return
	}
	// get report from db
	resultDB, errMdl := repository.GetReportCustomerPurchase(common.GormDB, ctxModel.AuthAccessTokenModel.CompanyID, request.CustomerID, fromDate, thruDate)
	if errMdl.Error != nil {
		return
	}

	var listData []dto.CustomerPurchaseResponse
	for _, v := range resultDB {
		temp := v.(entity.ReportSalesInvoiceEntity)
		listData = append(listData, dto.CustomerPurchaseResponse{
			CustomerID:      temp.CustomerID.Int64,
			ProductID:       temp.ProductID.Int64,
			Qty:             temp.Qty.Int64,
			LineGrossAmount: temp.LineGrossAmount.Float64,
			LineNetAmount:   temp.LineNetAmount.Float64,
		})
	}

	out.Data = listData
	out.Status.Message = ViewI18NMessage(ctxModel.AuthAccessTokenModel.Locale)
	return
}
