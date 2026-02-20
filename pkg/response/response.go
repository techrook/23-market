package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool         `json:"success"`
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
}

type ErrorDetail struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}

func OK(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Envelope{
		Success: true,
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, Envelope{
		Success: true,
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func Accepted(c *gin.Context, message string, meta *Meta) {
	c.JSON(http.StatusAccepted, Envelope{
		Success: true,
		Code:    http.StatusAccepted,
		Message: message,
		Meta:    meta,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Paginated(c *gin.Context, data interface{}, page, pageSize, total int, message string) {
	c.JSON(http.StatusOK, Envelope{
		Success: true,
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: total,
		},
	})
}

func Error(c *gin.Context, status int, appCode, message string, details interface{}, isProd bool) {
	envelope := Envelope{
		Success: false,
		Code:    status,
		Message: message,
		Error: &ErrorDetail{
			Code:    appCode,
			Message: message,
		},
	}

	if !isProd && details != nil {
		envelope.Error.Details = details
	}

	c.JSON(status, envelope)
}

func BadRequest(c *gin.Context, message string, details interface{}, isProd bool) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details, isProd)
}

func Unauthorized(c *gin.Context, message string, isProd bool) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil, isProd)
}

func Forbidden(c *gin.Context, message string, isProd bool) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message, nil, isProd)
}

func NotFound(c *gin.Context, resource string, isProd bool) {
	msg := resource + " not found"
	Error(c, http.StatusNotFound, "NOT_FOUND", msg, nil, isProd)
}

func Conflict(c *gin.Context, message string, details interface{}, isProd bool) {
	Error(c, http.StatusConflict, "CONFLICT", message, details, isProd)
}

func InternalError(c *gin.Context, message string, err error, isProd bool) {
	var details interface{}
	if !isProd && err != nil {
		details = err.Error()
	}
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, details, isProd)
}

func IsProduction(c *gin.Context) bool {
	if val, exists := c.Get("environment"); exists {
		if env, ok := val.(string); ok {
			return env == "production"
		}
	}
	return gin.Mode() == gin.ReleaseMode
}
