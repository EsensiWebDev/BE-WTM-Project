package response

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"wtm-backend/pkg/utils"
)

// Response represents a standard error response
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ResponseWithData represents a standard response with data
type ResponseWithData struct {
	Response `json:",inline"`
	Data     interface{} `json:"data,omitempty"`
}

// ResponseWithPagination represents a standard response with data and pagination
type ResponseWithPagination struct {
	ResponseWithData `json:",inline"`
	Pagination       *Pagination `json:"pagination"`
}

type Pagination struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewPagination(limit, page, total int) *Pagination {
	totalPages := 1
	if total > 0 && limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
	return &Pagination{
		Limit:      limit,
		Page:       page,
		Total:      total,
		TotalPages: totalPages,
	}
}

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, ResponseWithData{
		Response: Response{
			Status:  http.StatusOK,
			Message: message,
		},
		Data: data,
	})
}

func SuccessWithPagination(c *gin.Context, data interface{}, message string, pagination *Pagination) {
	c.JSON(http.StatusOK, ResponseWithPagination{
		ResponseWithData: ResponseWithData{
			Response: Response{
				Status:  http.StatusOK,
				Message: message,
			},
			Data: data,
		},
		Pagination: pagination,
	})
}

func EmptyList(c *gin.Context, message string, pagination *Pagination) {
	c.JSON(http.StatusOK, ResponseWithPagination{
		ResponseWithData: ResponseWithData{
			Response: Response{
				Status:  http.StatusOK,
				Message: message,
			},
			Data: []interface{}{},
		},
		Pagination: pagination,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Status:  status,
		Message: utils.CapitalizeWords(message),
	})
}

func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, ResponseWithData{
		Response: Response{
			Status:  http.StatusBadRequest,
			Message: "Validation Error",
		},
		Data: errors,
	})
}
