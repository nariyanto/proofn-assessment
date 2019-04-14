package helper

import "net/http"

type BaseResponse struct {
	Success    bool        `json:"success"`
	Status     int         `json:"status"`
	Message    interface{} `json:"message"`
	Data       interface{} `json:"data"`
	Exceptions error       `json:"exceptions"`
}

func CreateBaseResponse(success bool, status int, message interface{}, data interface{}, exceptions error) BaseResponse {
	baseResponse := BaseResponse{}
	baseResponse.Success = success
	baseResponse.Status = status
	baseResponse.Message = message
	baseResponse.Data = data
	baseResponse.Exceptions = exceptions
	return baseResponse
}

func CreateSuccessResponse(data interface{}, message interface{}) BaseResponse {
	successResponse := CreateBaseResponse(true, http.StatusOK, message, data, nil)
	return successResponse
}

func CreateErrorResponse(status int, message interface{}, exceptions error) BaseResponse {
	errorResponse := CreateBaseResponse(false, status, message, nil, exceptions)
	return errorResponse
}
