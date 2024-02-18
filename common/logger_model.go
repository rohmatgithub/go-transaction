package common

import (
	"encoding/json"
	"fmt"
)

type LoggerModel struct {
	Pid         string `json:"pid"`
	RequestID   string `json:"request_id"`
	Source      string `json:"source"`
	ClientID    string `json:"client_id"`
	UserID      string `json:"user_id" `
	Resource    string `json:"resource"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Class       string `json:"class"`
	ByteIn      int    `json:"byte_in"`
	ByteOut     int    `json:"byte_out"`
	Status      int    `json:"status" `
	Code        string `json:"code"`
	Message     string `json:"message"`
	Path        string `json:"path"`
	Method      string `json:"method"`
}

func GenerateLogModel(model LoggerModel) (output string) {
	return fmt.Sprintf("method:%s, path:%s, status:%d, pid:%s, request_id:%s, source:%s, client_id:%s, user_id:%s, resource:%s,"+
		"application:%s, version:%s, byte_in:%d, byte_out:%d, code:%s, message:%s, class:[%s]\n",
		model.Method, model.Path, model.Status, model.Pid, model.RequestID, model.Source, model.ClientID, model.UserID, model.Resource,
		model.Application, model.Version, model.ByteIn, model.ByteOut, model.Code, model.Message, model.Class)
}

func (object LoggerModel) String() string {
	b, err := json.Marshal(object)
	if err != nil {
		fmt.Println("error coy ", err)
		return ""
	}
	return string(b)
}
