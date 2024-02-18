package common

type ValidationInterface interface {
	ValidationAll(interface{}, *ContextModel) map[string]string
	ValidationCustom(string, string, *ContextModel) string
}
