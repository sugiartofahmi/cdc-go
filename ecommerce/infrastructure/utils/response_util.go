package utils

import "math"

const ApiVersion = "1.0.0"

type BaseResponse struct {
	Version    string               `json:"version"`
	StatusCode int                  `json:"status_code"`
	Message    string               `json:"message,omitempty"`
	Data       *any                 `json:"data,omitempty"`
	Items      *any                 `json:"items,omitempty"`
	Count      *int                 `json:"count,omitempty"`
	Meta       *MetaResponse        `json:"meta,omitempty"`
	Errors     *[]ErrorItemResponse `json:"errors,omitempty"`
}

type ErrorItemResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MetaResponse struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func SuccessResponse(statusCode int, message string, data any) *BaseResponse {
	return &BaseResponse{
		Version:    ApiVersion,
		StatusCode: statusCode,
		Message:    message,
		Data:       &data,
	}
}

func SuccessResponseList(statusCode int, message string, items any, count int) *BaseResponse {
	return &BaseResponse{
		Version:    ApiVersion,
		StatusCode: statusCode,
		Message:    message,
		Items:      &items,
		Count:      &count,
	}
}

func SuccessResponsePagination(statusCode int, message string, items any, meta MetaResponse) *BaseResponse {
	return &BaseResponse{
		Version:    ApiVersion,
		StatusCode: statusCode,
		Message:    message,
		Items:      &items,
		Meta:       &meta,
	}
}

func ErrorResponse(statusCode int, message string, errors *[]ErrorItemResponse) *BaseResponse {
	return &BaseResponse{
		Version:    ApiVersion,
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}
}

func PaginationMetaBuilder(page int, perPage int, total int) *MetaResponse {
	if perPage <= 0 {
		perPage = 10
	}
	if page <= 0 {
		page = 1
	}
	totalPage := int(math.Ceil(float64(total) / float64(perPage)))

	return &MetaResponse{
		Page:      page,
		PerPage:   perPage,
		Total:     total,
		TotalPage: totalPage,
	}
}
