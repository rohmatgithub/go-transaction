package repository

import (
	"database/sql"
	"fmt"
	"go-transaction/dto"
	"go-transaction/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func GetListDataDefault(gormDB *gorm.DB, query string, queryParam []interface{},
	dtoList dto.GetListRequest, searchBy []dto.SearchByParam,
	wrap func(rows *sql.Rows) (interface{}, error)) (result []interface{}, errMdl model.ErrorModel) {

	queryParam, queryCondition := SearchByParamToQuery(searchBy, queryParam)
	query += queryCondition + fmt.Sprintf(" ORDER BY %s ", dtoList.OrderBy)

	if dtoList.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d ", dtoList.Limit, countOffset(dtoList.Page, dtoList.Limit))
	}

	return ExecuteQuery(gormDB, query, queryParam, wrap)
}

func GetCountDataDefault(gormDB *gorm.DB, query string, queryParam []interface{}, searchBy []dto.SearchByParam) (result int64, errMdl model.ErrorModel) {

	queryParam, queryCondition := SearchByParamToQuery(searchBy, queryParam)
	query += queryCondition

	var temp sql.NullInt64
	gormCallBack := gormDB.Raw(query, queryParam...).Scan(&temp)
	if gormCallBack.Error != nil {
		errMdl = model.GenerateUnknownError(gormCallBack.Error)
		return
	}

	result = temp.Int64
	return
}

func SearchByParamToQuery(searchByParam []dto.SearchByParam, queryParam []interface{}) (resultQueryParam []interface{}, result string) {
	if len(queryParam) == 0 && len(searchByParam) > 0 {
		result = "WHERE \n"
	}
	var (
		operator          string
		searchConditionOr []dto.SearchByParam
	)
	index := len(queryParam)
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].Condition == "OR" {
			searchConditionOr = append(searchConditionOr, searchByParam[i])
			continue
		}
		if index > 0 {
			result += " AND "
		}
		index++
		if searchByParam[i].DataType == "enum" {
			searchByParam[i].SearchKey = "cast( " + searchByParam[i].SearchKey + " AS VARCHAR)"
		}
		searchByParam[i], operator = getOperator(searchByParam[i])

		if searchByParam[i].SearchOperator == "between" {
			operator = "between"
			result += fmt.Sprintf(" %s %s $%d AND $%d ", searchByParam[i].SearchKey, operator, len(queryParam)+1, len(queryParam)+2)
			searchValueSplit := strings.Fields(searchByParam[i].SearchValue)
			queryParam = append(queryParam, searchValueSplit[0])
			if len(searchValueSplit) > 1 {
				queryParam = append(queryParam, searchValueSplit[1])
			}
		} else if searchByParam[i].SearchOperator == "in" || searchByParam[i].SearchOperator == "not_in" {
			queryList := getQueryListValue(len(queryParam)+1, searchByParam[i].ListValue)
			result += " " + searchByParam[i].SearchKey + " " + operator + " (" + queryList + ")"
			queryParam = append(queryParam, searchByParam[i].ListValue...)
		} else {
			result += fmt.Sprintf(" %s %s $%d ", searchByParam[i].SearchKey, operator, len(queryParam)+1)
			queryParam = append(queryParam, searchByParam[i].SearchValue)
		}

	}

	queryOR := ""
	for i := 0; i < len(searchConditionOr); i++ {
		searchConditionOr[i], operator = getOperator(searchConditionOr[i])
		if searchConditionOr[i].SearchOperator == "between" {
			queryOR += fmt.Sprintf(" %s %s $%d AND $%d", searchConditionOr[i].SearchKey, searchConditionOr[i].SearchOperator, len(queryParam)+1, len(queryParam)+2)
			searchValueSplit := strings.Fields(searchConditionOr[i].SearchValue)
			queryParam = append(queryParam, searchValueSplit[0])
			if len(searchValueSplit) > 1 {
				queryParam = append(queryParam, searchValueSplit[1])
			}
		} else {
			switch operator {
			case "in", "not in":
				queryList := getQueryListValue(len(queryParam)+1, searchConditionOr[i].ListValue)
				queryOR += fmt.Sprintf(" %s %s (%s) ", searchConditionOr[i].SearchKey, operator, queryList)
				queryParam = append(queryParam, searchConditionOr[i].ListValue...)
			default:
				queryOR += fmt.Sprintf(" %s %s $%d ", searchConditionOr[i].SearchKey, operator, len(queryParam)+1)
				queryParam = append(queryParam, searchConditionOr[i].SearchValue)
			}
		}
		if i < len(searchConditionOr)-1 {
			queryOR += " OR "
		}
	}
	if len(searchConditionOr) > 0 {
		queryOR = fmt.Sprintf(" ( %s )", queryOR)
		if index > 0 {
			queryOR = " AND " + queryOR
		}
		result += queryOR
	}
	resultQueryParam = queryParam
	return
}

func getOperator(searchByParam dto.SearchByParam) (result dto.SearchByParam, operator string) {
	operator = searchByParam.SearchOperator
	switch operator {
	case "like":
		//searchByParam.SearchKey = "LOWER(" + searchByParam.SearchKey + ")"
		operator = "ilike"
		//searchByParam.SearchValue = strings.ToLower(searchByParam.SearchValue)
		searchByParam.SearchValue = "%" + searchByParam.SearchValue + "%"
	case "eq":
		operator = "="
	case "not_eq":
		operator = "!="
	case "not_like":
		operator = "not ilike "
		searchByParam.SearchKey = "LOWER(" + searchByParam.SearchKey + ")"
		//searchByParam.SearchValue = strings.ToLower(searchByParam.SearchValue)
		searchByParam.SearchValue = "%" + searchByParam.SearchValue + "%"
	case "between":
		operator = "between"
	//20-02-2022 -NEXCORE
	//-- Start Perubahan
	case "in":
		operator = "in"
	case "not_in":
		operator = "not in"
	}
	return searchByParam, operator
}

func getQueryListValue(currentIndex int, listValue []interface{}) (result string) {
	index := currentIndex
	for i := 0; i < len(listValue); i++ {
		if i != len(listValue)-1 {
			result += " $" + strconv.Itoa(index) + ","
		} else {
			result += " $" + strconv.Itoa(index)
		}
		index++
	}
	return result
}

func ExecuteQuery(gormDB *gorm.DB, query string, queryParam []interface{},
	wrap func(rows *sql.Rows) (interface{}, error)) (result []interface{}, errMdl model.ErrorModel) {

	rows, err := gormDB.Raw(query, queryParam...).Rows()
	if err != nil {
		errMdl = model.GenerateUnknownError(err)
		return
	}

	if rows != nil {
		defer func() {
			errMdl = closeRow(rows, errMdl)
		}()
		var temp interface{}
		for rows.Next() {
			temp, err = wrap(rows)
			if err != nil {
				errMdl = model.GenerateInternalDBServerError(err)
				return
			}
			result = append(result, temp)
		}
	}

	return
}

func closeRow(rows *sql.Rows, inputErr model.ErrorModel) (errMdl model.ErrorModel) {
	err := rows.Close()
	if err != nil {
		errMdl = model.GenerateInternalDBServerError(err)
	} else {
		errMdl = inputErr
	}
	return
}

func countOffset(page int, limit int) int {
	return (page - 1) * limit
}

func getQueryUpsert(tableName string, data map[string]model.UpsertModel) (query string, queryParam []interface{}) {
	var insertQuery, paramStr, primaryKey, updateQuery string
	i := 0
	for field, value := range data {
		i++
		insertQuery += fmt.Sprintf(" %s", field)
		paramStr += fmt.Sprintf(" $%d", i)
		if i < len(data) {
			insertQuery += ","
			paramStr += ","
		}

		queryParam = append(queryParam, value.Value)
		if value.PrimaryKey {
			if primaryKey != "" {
				primaryKey += ","
			}
			primaryKey += fmt.Sprintf(" %s", field)
		} else {
			switch field {
			case "created_by", "created_at":
				continue
			}
			if updateQuery != "" {
				updateQuery += ","
			}
			updateQuery += fmt.Sprintf(" %s=EXCLUDED.%s", field, field)
		}
	}
	//INSERT INTO "nexchief"."nexchief_account_parameter" ("nexchief_account_id","parameter_id","char_value","created_at","deleted")
	//VALUES (1,1,'YES','2020-09-01 15:21:51+00','FALSE')
	//ON CONFLICT ("nexchief_account_id","parameter_id") DO UPDATE SET "char_value"=EXCLUDED."char_value","created_at"=EXCLUDED."created_at","deleted"=EXCLUDED."deleted"
	query = "INSERT INTO " + tableName +
		fmt.Sprintf(" (%s) ", insertQuery) +
		"VALUES " +
		fmt.Sprintf(" (%s) ", paramStr) +
		"ON CONFLICT " +
		fmt.Sprintf(" (%s) ", primaryKey) +
		"DO UPDATE SET " + updateQuery + " RETURNING id"
	return
}

func getQueryUpsertMultiValues(tableName string, listData []map[string]model.UpsertModel) (query string, queryParam []interface{}) {
	var insertQuery, primaryKey, updateQuery string
	var values string
	i := 0
	var tempField []string

	tempMap := filterData(listData)
	index := 0
	for _, data := range tempMap {
		//data := listData[m]
		var paramStr string
		if index == 0 {
			for field, value := range data {
				i++
				insertQuery += fmt.Sprintf(" %s", field)
				paramStr += fmt.Sprintf(" $%d", len(queryParam)+1)
				if i < len(data) {
					paramStr += ","
					insertQuery += ","
				}
				tempField = append(tempField, field)
				queryParam = append(queryParam, value.Value)
				if value.PrimaryKey {
					if primaryKey != "" {
						primaryKey += ","
					}
					primaryKey += fmt.Sprintf(" %s", field)
				} else {
					switch field {
					case "created_by", "created_at":
						continue
					}
					if updateQuery != "" {
						updateQuery += ","
					}
					updateQuery += fmt.Sprintf(" %s=EXCLUDED.%s", field, field)
				}
			}
		} else {
			for j := 0; j < len(tempField); j++ {
				value := data[tempField[j]]
				paramStr += fmt.Sprintf(" $%d", len(queryParam)+1)
				if j < len(tempField)-1 {
					paramStr += ","
				}
				queryParam = append(queryParam, value.Value)
			}
		}

		i = 0
		values += fmt.Sprintf("(%s)", paramStr)
		if index < len(tempMap)-1 {
			values += ",\n"
		}
		index++
	}

	//INSERT INTO "nexchief"."nexchief_account_parameter" ("nexchief_account_id","parameter_id","char_value","created_at","deleted")
	//VALUES (1,1,'YES','2020-09-01 15:21:51+00','FALSE')
	//ON CONFLICT ("nexchief_account_id","parameter_id") DO UPDATE SET "char_value"=EXCLUDED."char_value","created_at"=EXCLUDED."created_at","deleted"=EXCLUDED."deleted"
	query = "INSERT INTO " + tableName +
		fmt.Sprintf(" (%s) \n", insertQuery) +
		"VALUES " + values + "\n" +
		//fmt.Sprintf(" (%s) ", paramStr) +
		"ON CONFLICT " +
		fmt.Sprintf(" (%s) \n", primaryKey) +
		"DO UPDATE SET " + updateQuery + " RETURNING id"
	return
}

func filterData(listData []map[string]model.UpsertModel) (result map[string]map[string]model.UpsertModel) {
	var pk []string
	for i := 0; i < len(listData); {
		dt := listData[i]
		for s, model := range dt {
			if model.PrimaryKey {
				pk = append(pk, s)
			}
		}
		break
	}

	// save di map biar gk duplicate
	tempMap := make(map[string]map[string]model.UpsertModel)
	for i := 0; i < len(listData); i++ {
		dt := listData[i]
		var pkStr string
		for j := 0; j < len(pk); j++ {
			pkStr += fmt.Sprintf("%v", dt[pk[j]].Value)
		}
		tempMap[pkStr] = dt
	}
	return tempMap
}
