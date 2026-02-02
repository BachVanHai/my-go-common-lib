package model

type CommonResponseHttp struct {
	Status          int    `json:"status"`
	Message         string `json:"message"`
	Code            string `json:"code,omitempty"` // Mã lỗi riêng của hệ thống
	Data            any    `json:"data,omitempty"`
	ClientMessageId string `json:"clientMessageId,omitempty"`
}
