package dto

import (
	"go-transaction/common"
	"go-transaction/constanta"
	"go-transaction/model"
	"strconv"
	"strings"
	"time"
)

var (
	ValidOperatorProduct    map[string]DefaultOperator
	ValidOperatorGeneral    map[string]DefaultOperator
	ValidOperatorSalesOrder map[string]DefaultOperator
	DefaultOrder            = []string{"id", "code", "name", "updated_at"}
)

func GenerateValidOperator() {
	ValidOperatorGeneral = map[string]DefaultOperator{
		"code":        {DataType: "char", Operator: []string{"eq", "like"}},
		"name":        {DataType: "char", Operator: []string{"eq", "like"}},
		"division_id": {DataType: "number", Operator: []string{"eq"}},
	}
	ValidOperatorSalesOrder = map[string]DefaultOperator{
		"order_number": {DataType: "char", Operator: []string{"eq", "like"}},
		"name":         {DataType: "char", Operator: []string{"eq", "like"}},
	}
}

type GetListRequest struct {
	Page    int           `json:"page"`
	Limit   int           `json:"limit"`
	OrderBy string        `json:"order_by"`
	Filter  string        `json:"filter"`
	ListID  []interface{} `json:"list_id"`
}

type SearchByParam struct {
	SearchKey      string
	DataType       string
	SearchOperator string
	SearchValue    string
	Condition      string
	ListValue      []interface{}
}

type DefaultOperator struct {
	DataType string   `json:"data_type"`
	Operator []string `json:"operator"`
}

func (input *GetListRequest) ValidateInputPageLimitAndOrderBy(validLimit []int, validOrderBy []string) model.ErrorModel {
	if input.Page < 1 {
		return model.GenerateFieldFormatWithRuleError("NEED_MORE_THAN", constanta.Page, "0")
	}

	input.Limit = checkLimit(validLimit, input.Limit)
	if input.Limit < 1 && input.Limit != -99 {
		return model.GenerateFieldFormatWithRuleError("NEED_MORE_THAN", constanta.Limit, "0")
	}

	input.OrderBy = strings.Trim(input.OrderBy, " ")
	if input.OrderBy == "" {
		input.OrderBy = validOrderBy[0]
	} else {
		orderBySplit := strings.Split(input.OrderBy, " ")
		var isAscending bool

		if !(len(orderBySplit) >= 1 && len(orderBySplit) <= 2) {
			return model.GenerateFormatFieldError(constanta.OrderBy)
		}

		if len(orderBySplit) == 1 {
			isAscending = true
		} else {
			if strings.ToUpper(orderBySplit[1]) == "ASC" {
				isAscending = true
			} else if strings.ToUpper(orderBySplit[1]) == "DESC" {
				isAscending = false
			} else {
				return model.GenerateFormatFieldError(constanta.OrderBy)
			}
		}

		if !common.ValidateStringContainInStringArray(validOrderBy, orderBySplit[0]) {
			return model.GenerateFormatFieldError(constanta.OrderBy)
		}

		input.OrderBy = orderBySplit[0] + " "
		if isAscending {
			input.OrderBy += "ASC"
		} else {
			input.OrderBy += "DESC"
		}
	}

	return model.ErrorModel{}
}

func checkLimit(validLimit []int, limit int) (result int) {
	if len(validLimit) == 0 {
		return 0
	}
	if limit != -99 && limit < validLimit[0] {
		return validLimit[0]
	}
	for i := 0; i < len(validLimit); i++ {
		if validLimit[i] == limit {
			return limit
		}
	}
	return 0
}

func (input *GetListRequest) ValidateFilter(validOperator map[string]DefaultOperator) (searchBy []SearchByParam, errMdl model.ErrorModel) {
	filter := input.Filter
	if filter != "" {
		filterSplitComma := strings.Split(filter, ",")
		for i := 0; i < len(filterSplitComma); i++ {
			filterIndex := strings.TrimSpace(filterSplitComma[i])
			filterIndexSplitSpace := strings.Split(filterIndex, " ")
			if len(filterIndexSplitSpace) > 2 {
				searchKey := strings.Trim(filterIndexSplitSpace[0], " ")
				operator := strings.Trim(filterIndexSplitSpace[1], " ")
				searchValue := ""
				for j := 2; j < len(filterIndexSplitSpace); j++ {
					searchValue += filterIndexSplitSpace[j] + " "
				}
				searchValue = strings.Trim(searchValue, " ")

				searchBy = append(searchBy, SearchByParam{
					DataType:       validOperator[searchKey].DataType,
					SearchKey:      searchKey,
					SearchOperator: operator,
					SearchValue:    searchValue,
				})

				if !isOperatorValid(searchKey, searchValue, operator, validOperator) {
					errMdl = model.GenerateFormatFieldError(constanta.Filter)
					return
				}
			} else {
				errMdl = model.GenerateFormatFieldError(constanta.Filter)
				return
			}
		}
	}
	errMdl = model.ErrorModel{}
	return
}

func isOperatorValid(key string, value string, operator string, validOperator map[string]DefaultOperator) bool {
	if validOperator[key].Operator == nil {
		return false
	} else {
		if validOperator[key].DataType == "number" {
			_, err := strconv.Atoi(value)
			if err != nil {
				return false
			}
		} else if validOperator[key].DataType == "date" {
			dateSplit := strings.Fields(value)
			for _, element := range dateSplit {
				_, errS := time.Parse("2006-01-02", element)
				if errS != nil {
					return false
				}
			}
		}
		if operator == "between" {
			if len(strings.Fields(value)) < 2 {
				return false
			}
		}
		return common.ValidateStringContainInStringArray(validOperator[key].Operator, operator)
	}
}
