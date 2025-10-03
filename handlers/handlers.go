package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

var (
	ErrMissingAPIArguments    = errors.New("MISSING_API_ARGUMENTS")
	ErrInsufficientPermission = errors.New("INSUFFICIENT_PERMISSION")
	ErrUnknown                = errors.New("UNKNOWN_ERROR")
)

var (
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	JSON_PARSE_ERROR      = "JSON_PARSE_ERROR"
)

type APIResponseErrJson struct {
	Errcode string      `json:"errcode"`
	ErrData interface{} `json:"errdata"`
}

type APIResponse struct {
	Err  interface{} `json:"err"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type EmptyResp struct{}

func APIResponseOK(c *gin.Context, data interface{}, msg string) {
	responseObj := &APIResponse{
		Err:  nil,
		Data: data,
		Msg:  msg,
	}
	c.JSON(http.StatusOK, responseObj)
}

func apiResponseErr(status int, c *gin.Context, errcode string, errdata interface{}, msg string) {
	responseObj := &APIResponse{
		Err: &APIResponseErrJson{
			Errcode: errcode,
			ErrData: errdata,
		},
		Data: nil,
		Msg:  msg,
	}
	c.JSON(status, responseObj)
}

func APIResponseBadRequest(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusBadRequest, c, errcode, errdata, msg)
}

func APIResponseUnauthorized(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusUnauthorized, c, errcode, errdata, msg)
}

func APIResponseForbidden(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusForbidden, c, errcode, errdata, msg)
}

func APIResponseConflict(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusConflict, c, errcode, errdata, msg)
}

func APIResponseGone(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusGone, c, errcode, errdata, msg)
}

func APIResponseUnprocessableEntity(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusUnprocessableEntity, c, errcode, errdata, msg)
}

func APIResponseNotAcceptable(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusNotAcceptable, c, errcode, errdata, msg)
}

func APIResponseInternalServerError(c *gin.Context, errcode string, errdata interface{}, msg string) {
	apiResponseErr(http.StatusInternalServerError, c, errcode, errdata, msg)
}

func APIFailedInternalAPICall(c *gin.Context, errcode string, errdata interface{}, msg string, statusCode int) {
	apiResponseErr(statusCode, c, errcode, errdata, msg)
}
