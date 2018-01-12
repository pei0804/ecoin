package httputil

import "fmt"

// HTTPError エラー用
type HTTPError struct {
	Status  int   `json:"status"`
	Code    int   `json:"code"`
	Message error `json:"message"`
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}
