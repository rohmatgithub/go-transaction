package restapi

import (
	"go-transaction/common"
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
	app.Get("/order", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", ListSalesOrder)
	})
	app.Get("/order/count", func(c *fiber.Ctx) error {
		return ae.ServeJwtToken(c, "", CountListSalesOrder)
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
	listParam, errMdl := util_controller.ValidateCount(c, dto.ValidOperatorGeneral)
	if errMdl.Error != nil {
		return
	}
	out, errMdl = service.CountListSalesOrder(listParam, ctx)
	if errMdl.Error != nil {
		return
	}

	return
}
