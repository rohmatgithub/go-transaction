package common

import ut "github.com/go-playground/universal-translator"

type ContextModel struct {
	LoggerModel          LoggerModel
	AuthAccessTokenModel AuthAccessTokenModel
	PermissionHave       string
	Translator           ut.Translator
}

type AuthAccessTokenModel struct {
	ResourceUserID int64
	CompanyID      int64
	BranchID       int64
	ClientID       string
	Locale         string
}
