package app

import (
	"strconv"

	"github.com/golang/glog"
)

const (
	success = iota
	faield
)

// APIMessage 封装消息给前端
type APIMessage struct {
	Data    interface{} `json:"data,omitempty"`
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
}

// NewMsg 统一封装消息
func NewMsg(code int, data interface{}, details ...interface{}) (int, *APIMessage) {
	m := &APIMessage{}
	m.Code = strconv.Itoa(code)

	if code == 200 {
		m.Status = success
		switch d := data.(type) {
		case string:
			m.Message = d
		default:
			m.Data = data
		}

	} else {
		m.Status = faield
		switch d := data.(type) {
		case string:
			glog.ErrorDepth(1, d)
			m.Message = d
		case error:
			// TODO: handle sql error
			glog.ErrorDepth(1, d)
			m.Message = d.Error()
		default:
			m.Message = "发生不可知错误."
		}
	}

	return 200, m
}
