package model

import (
	"errors"
	"fmt"
	"runtime"
)

type ErrorModel struct {
	Code                  int
	Error                 error
	CausedBy              error
	ErrorParameter        []ErrorParameter
	AdditionalInformation interface{}
	Line                  string
}

type ErrorParameter struct {
	ErrorParameterKey   string
	ErrorParameterValue string
}

func customLogger() string {
	// Get the program counter of the current function.
	_, f, l, _ := runtime.Caller(3)
	return fmt.Sprintf("%s:%d", f, l)
}

func GenerateErrorModel(code int, err string, causedBy error) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.CausedBy = causedBy
	errModel.Line = customLogger()
	return errModel
}

func GenerateErrorModelWithErrorParam(code int, err string, errorParam []ErrorParameter) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.ErrorParameter = errorParam
	errModel.Line = customLogger()
	return errModel
}
