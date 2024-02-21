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

func RouteSalesOrder(app fiber.Router) {
	var ae util_controller.AbstractController
	app.Post("/order", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", CreateOrder)
	})
	app.Put("/order", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", UpdateOrder)
	})
	app.Get("/order", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", ListSalesOrder)
	})
	app.Get("/order/count", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", CountListSalesOrder)
	})
	app.Get(fmt.Sprintf("/order/:%s", constanta.ParamID), func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", GetDetailSalesOrder)
	})
}
func CreateOrder(c *fiber.Ctx, contextModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	var request dto.SalesOrderRequest
	err := c.BodyParser(&request)
	if err != nil {
		errMdl = model.GenerateInvalidRequestError(err)
		return
	}
	out, errMdl = service.InsertSalesOrder(request, contextModel)
	if errMdl.Error != nil {
		return
	}

	return
}

func UpdateOrder(c *fiber.Ctx, contextModel *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	var request dto.SalesOrderRequest
	err := c.BodyParser(&request)
	if err != nil {
		errMdl = model.GenerateInvalidRequestError(err)
		return
	}
	out, errMdl = service.UpdateSalesOrder(request, contextModel)
	if errMdl.Error != nil {
		return
	}

	return
}

func ListSalesOrder(c *fiber.Ctx, ctx *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	// set to search param
	dtoList, listParam, errMdl := util_controller.ValidateList(c, []string{"id", "order_date", "updated_at"}, dto.ValidOperatorSalesOrder)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.ListSalesOrder(dtoList, listParam, ctx)
	if errMdl.Error != nil {
		return
	}
	return
}

func CountListSalesOrder(c *fiber.Ctx, ctx *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	// set to search param
	listParam, errMdl := util_controller.ValidateCount(c, dto.ValidOperatorSalesOrder)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.CountListSalesOrder(listParam, ctx)
	if errMdl.Error != nil {
		return
	}

	return
}

func GetDetailSalesOrder(c *fiber.Ctx, ctx *common.ContextModel) (out dto.Payload, errMdl model.ErrorModel) {
	id, errMdl := util_controller.GetParamID(c)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.GetDetailSalesOrder(id, ctx)
	if errMdl.Error != nil {
		return
	}

	return
}
