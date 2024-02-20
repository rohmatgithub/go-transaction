package restapi

import (
	"fmt"
	"go-transaction/common"
	"go-transaction/constanta"
	"go-transaction/controller/restapi/util_controller"
	"go-transaction/dto"
	"go-transaction/model"
	"go-transaction/service"

	"github.com/gofiber/fiber/v2"
)

func RouteSalesInvoice(app fiber.Router) {
	var ae util_controller.AbstractController
	app.Post(fmt.Sprintf("/invoice/:%s", constanta.ParamID), func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", CreateInvoice)
	})
	app.Get("/invoice", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", ListSalesInvoice)
	})
	app.Get("/invoice/count", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", CountListSalesInvoice)
	})
	app.Get(fmt.Sprintf("/invoice/:%s", constanta.ParamID), func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", GetDetailSalesInvoice)
	})
}
func CreateInvoice(c *fiber.Ctx, contextModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	id, errMdl := util_controller.GetParamID(c)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.CreateSalesInvoice(id, contextModel)
	if errMdl.Error != nil {
		return
	}

	return
}

func ListSalesInvoice(c *fiber.Ctx, ctx *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	// set to search param
	dtoList, listParam, errMdl := util_controller.ValidateList(c, []string{"id", "Invoice_date", "updated_at"}, dto.ValidOperatorSalesInvoice)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.ListSalesInvoice(dtoList, listParam, ctx)
	if errMdl.Error != nil {
		return
	}
	return
}

func CountListSalesInvoice(c *fiber.Ctx, ctx *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	// set to search param
	listParam, errMdl := util_controller.ValidateCount(c, dto.ValidOperatorGeneral)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.CountListSalesInvoice(listParam, ctx)
	if errMdl.Error != nil {
		return
	}

	return
}

func GetDetailSalesInvoice(c *fiber.Ctx, ctx *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	id, errMdl := util_controller.GetParamID(c)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.GetDetailSalesInvoice(id, ctx)
	if errMdl.Error != nil {
		return
	}

	return
}
