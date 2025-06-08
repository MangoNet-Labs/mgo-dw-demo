package response

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"time"
)

// BaseResponse Common response structure
type BaseResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// PageResult Paging data structure
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// OkJson Successful response
func OkJson(w http.ResponseWriter, data interface{}) {
	resp := BaseResponse{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	httpx.OkJson(w, resp)
}

// OkPage Successful paging response
func OkPage(w http.ResponseWriter, list interface{}, total int64, page, pageSize int) {
	resp := BaseResponse{
		Code:    0,
		Message: "success",
		Data: PageResult{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		Timestamp: time.Now().Unix(),
	}
	httpx.OkJson(w, resp)
}

// FailJson Error Response
func FailJson(w http.ResponseWriter, msg string, code int) {
	resp := BaseResponse{
		Code:      code,
		Message:   msg,
		Timestamp: time.Now().Unix(),
	}
	httpx.OkJson(w, resp)
}
