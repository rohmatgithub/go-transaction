package restapi

import (
	"fmt"
	"go-transaction/config"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Router() error {
	app := fiber.New(fiber.Config{
		//Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Transaction App v1.0.0",
		ColorScheme:   fiber.Colors{Green: ""},
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})
	app.Use(requestid.New())
	//app.Use(recoverfiber.New())
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				customErrorHandler(c, fmt.Errorf("%v", r))
			}
		}()
		return c.Next()
	})
	app.Use(middleware)

	v1 := app.Group("/v1/trans")
	RouteSalesOrder(v1)
	RouteSalesInvoice(v1)

	// exampleRepository := example_repository.NewExampleRepository(common.GormDB)
	// exampleService := example_service.NewExampleService(exampleRepository)
	// exampleController := example_controller.NewExampleController(exampleService)
	// exampleController.Route(v1)

	app.Use(NotFoundHandler)
	return app.Listen(fmt.Sprintf(":%d", config.ApplicationConfiguration.GetServerConfig().Port))
}
