package apierror

import (
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper-API/pkg/apperrors"
	"net/http"
)

type ApiError struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func New(err error) ApiError {
	switch errors.Code(err) {
	case errors.Code(apperrors.NotFound):
		return ApiError{Status: http.StatusNotFound, Code: "not_found", Message: err.Error(), Data: errors.Data(err)}
	default:
		return ApiError{Status: http.StatusInternalServerError, Code: "internal", Message: err.Error(), Data: errors.Data(err)}
	}
}