package model

import "fmt"

type CommonResponseHttp struct {
	Status          int    `json:"status"`
	Message         string `json:"message"`
	Code            string `json:"code,omitempty"` // Mã lỗi riêng của hệ thống
	Data            any    `json:"data,omitempty"`
	ClientMessageId string `json:"clientMessageId,omitempty"`
}

// Hàm này giúp struct này thỏa mãn interface 'error' của Go
func (e *CommonResponseHttp) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
