package common

import (
	"errors"
	"go-transaction/constanta"
	"reflect"
	"strings"

	en_US "github.com/go-playground/locales/en"
	id_ID "github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	idtranslations "github.com/go-playground/validator/v10/translations/id"
	"github.com/gofiber/fiber/v2/log"
)

type goValidatorImpl struct {
	Validate     *validator.Validate
	TranslatorId ut.Translator
	TranslatorEn ut.Translator
}

func NewGoValidator() ValidationInterface {
	// set translator
	en := en_US.New()
	id := id_ID.New()
	uni := ut.New(en, id)

	validate := validator.New()
	translatorId, _ := uni.GetTranslator("id")
	translatorEn, _ := uni.GetTranslator("en")

	err := entranslations.RegisterDefaultTranslations(validate, translatorEn)
	if err != nil {
		log.Fatal(err)
	}
	err = idtranslations.RegisterDefaultTranslations(validate, translatorId)
	if err != nil {
		log.Fatal(err)
	}
	return &goValidatorImpl{
		Validate:     validate,
		TranslatorId: translatorId,
		TranslatorEn: translatorEn,
	}
}

func (v *goValidatorImpl) ValidationAll(input interface{}, contextModel *ContextModel) map[string]string {
	err := v.Validate.Struct(input)
	if err != nil {

		// translate all error at once
		var errs validator.ValidationErrors
		errors.As(err, &errs)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned, and you'll see something surprising
		// translations are i18n aware!!!!
		result := make(map[string]string)
		for _, fieldError := range errs {
			field, _ := reflect.TypeOf(input).FieldByName(fieldError.Field())
			jsonTag := field.Tag.Get("json")
			//result[jsonTag] = fieldError.Translate(v.getTranslator(contextModel.AuthAccessTokenModel.Locale))
			result[jsonTag] = strings.Replace(fieldError.Translate(v.getTranslator(contextModel.AuthAccessTokenModel.Locale)), fieldError.Field()+" ", "", 1)
		}

		return result
	}
	return nil
}

func (v *goValidatorImpl) ValidationCustom(name string, tag string, contextModel *ContextModel) string {
	err := v.Validate.Var(name, tag)
	if err != nil {
		// translate all error at once
		var errs validator.ValidationErrors
		errors.As(err, &errs)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned, and you'll see something surprising
		// translations are i18n aware!!!!
		for _, fieldError := range errs {
			//result[jsonTag] = fieldError.Translate(v.getTranslator(contextModel.AuthAccessTokenModel.Locale))
			return fieldError.Translate(v.getTranslator(contextModel.AuthAccessTokenModel.Locale))
		}
	}
	return ""
}

func (v *goValidatorImpl) getTranslator(locale string) ut.Translator {
	if locale == constanta.LanguageId {
		return v.TranslatorId
	}

	return v.TranslatorEn
}
