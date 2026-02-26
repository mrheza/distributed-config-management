package httpresponse

import (
	sharedhttpresponse "github.com/mrheza/distributed-config-management/shared/httpresponse"

	"github.com/gin-gonic/gin"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidationFieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code,omitempty"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message"`
}

type ValidationErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Fields  []ValidationFieldError `json:"fields,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ValidationErrorResponse struct {
	Error ValidationErrorDetail `json:"error"`
}

func Error(c *gin.Context, status int, code, message string) {
	sharedhttpresponse.Error(c, status, code, message)
}

func InternalServerError(c *gin.Context, err error) {
	sharedhttpresponse.InternalServerError(c, err)
}

func Unauthorized(c *gin.Context) {
	sharedhttpresponse.Unauthorized(c)
}

func NotFound(c *gin.Context, message string) {
	sharedhttpresponse.NotFound(c, message)
}

func BadRequest(c *gin.Context, message string) {
	sharedhttpresponse.BadRequest(c, message)
}

func ValidationError(c *gin.Context, err error, requestStruct interface{}) {
	sharedhttpresponse.ValidationError(c, err, requestStruct)
}

func FromError(c *gin.Context, err error) {
	sharedhttpresponse.FromError(c, err)
}
