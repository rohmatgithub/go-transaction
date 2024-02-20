package dto

import "time"

type Header struct {
	RequestID string    `json:"request_id"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

type Status struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
}

type CompanyBranchResponse struct {
	Header  Header `json:"header"`
	Payload struct {
		Status Status `json:"status"`
		Data   []struct {
			ID        int       `json:"id"`
			Code      string    `json:"code"`
			Name      string    `json:"name"`
			Address1  string    `json:"address_1"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"data"`
		Other interface{} `json:"other"`
	} `json:"payload"`
}

type ListCustomerResponse struct {
	Header  Header `json:"header"`
	Payload struct {
		Status Status `json:"status"`
		Data   []struct {
			ID        int64     `json:"id"`
			Code      string    `json:"code"`
			Name      string    `json:"name"`
			Phone     string    `json:"phone"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"data"`
		Other interface{} `json:"other"`
	} `json:"payload"`
}

type ListProductResponse struct {
	Header  Header `json:"header"`
	Payload struct {
		Status Status `json:"status"`
		Data   []struct {
			ID           int64     `json:"id"`
			Code         string    `json:"code"`
			Name         string    `json:"name"`
			SellingPrice float64   `json:"selling_price"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			// Category     dto.StructGeneral `json:"category"`
		} `json:"data"`
		Other interface{} `json:"other"`
	} `json:"payload"`
}
