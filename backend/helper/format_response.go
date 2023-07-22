package helper

import "net/http"

type BaseResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SuccessResponse format on response success
func SuccessResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// SuccessCreatedResponse format on response success created
func SuccessCreatedResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// BadRequestResponse format on response error bad request
func BadRequestResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    http.StatusBadRequest,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

// NotFoundResponse format on response error not found
func NotFoundResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    http.StatusNotFound,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

// UnauthorizedResponse format on response error unauthorized
func UnauthorizedResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    http.StatusUnauthorized,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

// ForbiddenResponse format on response error forbidden
func ForbiddenResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    http.StatusForbidden,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

// ConflictResponse format on response error conflict
func ConflictResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    http.StatusConflict,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

// InternalServerErrorResponse format on response internal server error
func InternalServerErrorResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    http.StatusInternalServerError,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}
