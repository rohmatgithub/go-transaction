package util_controller

import (
	"context"
	"go-transaction/common"
	"go-transaction/config"
	"go-transaction/constanta"
	"go-transaction/dto"
	"go-transaction/model"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type AbstractController struct {
}

func (ae AbstractController) ServeJwtToken(c *fiber.Ctx, menuConst string, runFunc func(*fiber.Ctx, *common.ContextModel) (dto.Payload, model.ErrorModel)) error {
	// validate client_id
	tokenStr := c.Get(constanta.TokenHeaderNameConstanta)

	validateFunc := func(contextModel *common.ContextModel) (errMdl model.ErrorModel) {
		if tokenStr == "" {
			errMdl = model.GenerateUnauthorizedClientError()
			return
		}
		// cek token expired
		tokenModel, errMdl := model.JWTToken{}.ParsingJwtTokenInternal(tokenStr)
		if errMdl.Error != nil {
			return
		}

		contextModel.AuthAccessTokenModel.ResourceUserID = tokenModel.UserID
		contextModel.AuthAccessTokenModel.CompanyID = tokenModel.CompanyID
		contextModel.AuthAccessTokenModel.BranchID = tokenModel.BranchID
		return
	}

	return ae.serve(c, validateFunc, runFunc)
}

func (ae AbstractController) serve(c *fiber.Ctx,
	validateFunc func(contextModel *common.ContextModel) model.ErrorModel,
	runFunc func(*fiber.Ctx, *common.ContextModel) (dto.Payload, model.ErrorModel)) (err error) {
	var (
		response     dto.StandardResponse
		payload      dto.Payload
		contextModel common.ContextModel
	)

	requestID := c.Locals("requestid").(string)
	logModel := c.Context().Value(constanta.ApplicationContextConstanta).(*common.LoggerModel)

	contextModel.LoggerModel = *logModel
	response.Header = dto.HeaderResponse{
		RequestID: requestID,
		Version:   config.ApplicationConfiguration.GetServerConfig().Version,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	defer func() {
		if r := recover(); r != nil {
			contextModel.LoggerModel.Message = string(debug.Stack())
			generateEResponseError(c, &contextModel, &payload, model.GenerateUnknownError(nil))
		}
		response.Payload = payload

		adaptor.CopyContextToFiberContext(context.WithValue(c.Context(), constanta.ApplicationContextConstanta, &contextModel.LoggerModel), c.Context())
		err = c.JSON(response)
	}()
	// validate
	errMdl := validateFunc(&contextModel)
	if errMdl.Error != nil {
		generateEResponseError(c, &contextModel, &payload, errMdl)
		return
	}
	payload, errMdl = runFunc(c, &contextModel)
	if errMdl.Error != nil {
		generateEResponseError(c, &contextModel, &payload, errMdl)
	} else {
		payload.Status.Success = true
		payload.Status.Code = "OK"
	}
	return
}

func generateEResponseError(c *fiber.Ctx, ctxModel *common.ContextModel, payload *dto.Payload, errMdl model.ErrorModel) {
	ctxModel.LoggerModel.Code = errMdl.Error.Error()
	ctxModel.LoggerModel.Class = errMdl.Line
	if errMdl.CausedBy != nil {
		ctxModel.LoggerModel.Message = errMdl.CausedBy.Error()
	}
	// write failed
	c.Status(errMdl.Code)
	payload.Status.Success = false
	payload.Status.Code = errMdl.Error.Error()
	payload.Status.Message = common.GenerateI18NErrorMessage(errMdl, ctxModel.AuthAccessTokenModel.Locale)
}

func GetParamID(c *fiber.Ctx) (id int64, errMdl model.ErrorModel) {
	str := c.Params(constanta.ParamID)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		errMdl = model.GenerateInvalidRequestError(err)
		return
	}
	return
}
