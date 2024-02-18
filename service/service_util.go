package service

import (
	"go-transaction/common"
	"go-transaction/constanta"
	"go-transaction/model"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
)

func InsertI18NMessage(language string) string {
	return GenerateI18NMessage("SUCCESS_INSERT_MESSAGE", language)
}

func UpdateI18NMessage(language string) string {
	return GenerateI18NMessage("SUCCESS_UPDATE_MESSAGE", language)
}

func ListI18NMessage(language string) string {
	return GenerateI18NMessage("SUCCESS_LIST_MESSAGE", language)
}

func CountI18NMessage(language string) string {
	return GenerateI18NMessage("SUCCESS_COUNT_MESSAGE", language)
}

func DeleteI18NMessage(language string) string {
	return GenerateI18NMessage("SUCCESS_DELETE_MESSAGE", language)
}

func ViewI18NMessage(language string) string {
	return GenerateI18NMessage("SUCCESS_VIEW_MESSAGE", language)
}
func GenerateI18NMessage(messageID string, language string) (output string) {
	return common.GenerateI18NServiceMessage(common.CommonBundle, messageID, language, nil)
}

func HitToResourceOther(uri string, method string, ctxModel *common.ContextModel) (output []byte, errMdl model.ErrorModel) {
	// Create an HTTP client
	client := &http.Client{}

	tokenInternal, errMdl := model.GetTokenInternal(ctxModel.AuthAccessTokenModel.ResourceUserID, ctxModel.AuthAccessTokenModel.CompanyID)
	if errMdl.Error != nil {
		return
	}
	// uri := config.ApplicationConfiguration.GetUriResouce().MasterData + "/v1/master/company/" + strconv.Itoa(int(userDB.CompanyID.Int64))
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		log.Error("Error creating request:", err)
		errMdl = model.GenerateUnknownError(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(constanta.TokenHeaderNameConstanta, tokenInternal)

	// Make the request
	response, err := client.Do(req)
	if err != nil {
		log.Error("Error making request:", err)
		errMdl = model.GenerateUnknownError(err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	output, err = io.ReadAll(response.Body)
	if err != nil {
		log.Error("Error reading response:", err)
		errMdl = model.GenerateUnknownError(err)
		return
	}

	return
}
