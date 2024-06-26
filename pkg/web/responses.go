package web

type status string

const (
	OK    = status("OK")
	ERROR = status("ERROR")
)

type (
	Response struct {
		Status  status `json:"status"`
		MsgCode string `json:"msg_code"`
		Data    any    `json:"data,omitempty"`
	}
)

func OKResponse(msgCode string, body, meta any) Response {
	return Response{
		Status:  OK,
		MsgCode: msgCode,
		Data:    body,
	}
}

func ErrorResponse(msgCode string, body, meta any) Response {
	return Response{
		Status:  ERROR,
		MsgCode: msgCode,
		Data:    body,
	}
}
