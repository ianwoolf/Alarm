package output

import (
	"github.com/lodastack/alarm/common"
)

type AlertMsg struct {
	Users []string
	Msg   string
}

type HandleFunc func(alertMsg AlertMsg) error

var Handlers map[string]HandleFunc

func init() {
	Handlers = make(map[string]HandleFunc)
	Handlers["mail"] = SendEMail // TODO
}

func NewAlertMsg(Users []string, msg string) AlertMsg {
	return AlertMsg{Users: Users, Msg: msg}
}

func Send(alertType []string, alertMsg AlertMsg) error {
	//
	alertType = append(alertType, "mail")
	alertType = common.RemoveDuplicateAndEmpty(alertType)
	//
	for _, handler := range alertType {
		handlerFunc, ok := Handlers[handler]
		if !ok {
			// TODO
			continue
		}
		if err := handlerFunc(alertMsg); err != nil {
			return err
		}
	}
	return nil
}
