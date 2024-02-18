package model

// ==================== ERROR DTO ===================

func GenerateEmptyFieldError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-DTO-001", errorParam)
}

func GenerateFormatFieldError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-DTO-002", errorParam)
}

func GenerateFieldFormatWithRuleError(ruleName string, fieldName string, additionalInfo string) ErrorModel {
	errorParam := make([]ErrorParameter, 3)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	errorParam[1].ErrorParameterKey = "RuleName"
	errorParam[1].ErrorParameterValue = ruleName
	errorParam[2].ErrorParameterKey = "Other"
	errorParam[2].ErrorParameterValue = additionalInfo
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-DTO-003", errorParam)
}

func GenerateInvalidRequestError(causedBy error) ErrorModel {
	return GenerateErrorModel(400, "E-4-TRS-DTO-004", causedBy)
}

func GenerateFailedValidate() ErrorModel {
	return GenerateErrorModel(400, "E-4-TRS-DTO-005", nil)
}

// ====================== ERROR SERVICE  ===================

func GenerateUnauthorizedClientError() ErrorModel {
	return GenerateErrorModel(401, "E-1-TRS-SRV-001", nil)
}

func GenerateVerifyPasswordInvalidError() ErrorModel {
	return GenerateErrorModel(401, "E-1-TRS-SRV-002", nil)
}

func GenerateExpiredTokenError() ErrorModel {
	return GenerateErrorModel(401, "E-1-TRS-SRV-003", nil)
}

func GenerateHasUsedDataError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-SRV-001", errorParam)
}

func GenerateUnknownDataError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-SRV-002", errorParam)
}

func GenerateDataLockedError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-SRV-003", errorParam)
}

func GenerateNotAccessError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-SRV-004", errorParam)
}

func GenerateNotChangedDataError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-SRV-005", errorParam)
}

func GenerateNotDeleteDataError(fieldName string) ErrorModel {
	errorParam := make([]ErrorParameter, 1)
	errorParam[0].ErrorParameterKey = "FieldName"
	errorParam[0].ErrorParameterValue = fieldName
	return GenerateErrorModelWithErrorParam(400, "E-4-TRS-SRV-006", errorParam)
}

func GenerateUnknownError(causedBy error) ErrorModel {
	return GenerateErrorModel(500, "E-5-TRS-SRV-001", causedBy)
}

func GenerateInternalDBServerError(causedBy error) ErrorModel {
	return GenerateErrorModel(500, "E-5-TRS-DBS-001", causedBy)
}
