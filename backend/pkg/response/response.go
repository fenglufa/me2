package response

import (
	"net/http"

	"github.com/me2/pkg/errcode"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

// Error 错误响应
func Error(err *errcode.Error) *Response {
	return &Response{
		Code: err.Code,
		Msg:  err.Msg,
		Data: nil,
	}
}

// HttpSuccess HTTP 成功响应
func HttpSuccess(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Success(data))
}

// HttpError HTTP 错误响应
func HttpError(w http.ResponseWriter, err *errcode.Error) {
	httpx.WriteJson(w, http.StatusOK, Error(err))
}

// HttpErrorWithStatus HTTP 错误响应(自定义状态码)
func HttpErrorWithStatus(w http.ResponseWriter, status int, err *errcode.Error) {
	httpx.WriteJson(w, status, Error(err))
}

// PageResponse 分页响应
type PageResponse struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     interface{} `json:"list"`
}

// PageSuccess 分页成功响应
func PageSuccess(total int64, page, pageSize int, list interface{}) *Response {
	return Success(&PageResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	})
}
