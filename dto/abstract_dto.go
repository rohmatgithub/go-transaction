package dto

import (
	"go-transaction/constanta"
	"go-transaction/model"
	"time"
)

type StandardResponse struct {
	Header  HeaderResponse `json:"header"`
	Payload Payload        `json:"payload"`
}

type HeaderResponse struct {
	RequestID string `json:"request_id"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type Payload struct {
	Status StatusPayload `json:"status"`
	Data   interface{}   `json:"data"`
	Other  interface{}   `json:"other"`
}

type StatusPayload struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
}

type AbstractDto struct {
	ID           int64  `json:"id"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

type StructGeneral struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (dto *AbstractDto) ValidateUpdateGeneral() (errMdl model.ErrorModel) {
	if dto.ID < 1 {
		errMdl = model.GenerateUnknownDataError(constanta.ID)
		return
	}

	times, err := time.Parse(constanta.FormatDateGeneral, dto.UpdatedAtStr)
	if err != nil {
		errMdl = model.GenerateFormatFieldError(constanta.UpdatedAt)
		return
	}

	dto.UpdatedAt = times
	return
}
