package model

import "fmt"

// BusinessError đại diện cho Business Exception
type BusinessError struct {
	Code    string
	Message string
}

// Hàm này giúp struct này thỏa mãn interface 'error' của Go
func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
