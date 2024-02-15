package common

import (
	"encoding/json"
	"github.com/julioisaac/users/logger"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/julioisaac/users/domainerror"
)

var mapCodes = map[domainerror.ErrorCode]int{
	// BadRequest statuses
	domainerror.CreateUserError:      http.StatusBadRequest,
	domainerror.GetUserDefaultError:  http.StatusBadRequest,
	domainerror.GetUserNotFoundError: http.StatusBadRequest,

	// Internal statuses
	domainerror.Default: http.StatusInternalServerError,
}

func toMapStatusResponse(ec domainerror.ErrorCode) int {
	if status, ok := mapCodes[ec]; ok {
		return status
	}

	return http.StatusInternalServerError
}

// ToDomainErrorResponse given error
func ToDomainErrorResponse(err error) ApiError {
	apiError := ApiError{
		Code:       string(domainerror.Default),
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}

	if de, ok := err.(*domainerror.Error); ok {
		if de.Code == domainerror.ErrorCodeNone {
			de.Code = domainerror.Default
			de.Message = "unknown error, check if its correctly mapped ;)"
		}
		apiError.StatusCode = toMapStatusResponse(de.Code)
		apiError.Code = string(de.Code)
		apiError.Message = de.Message
		apiError.Detail = de.Detail
	}

	if ve, ok := err.(validation.Errors); ok {
		apiError.StatusCode = http.StatusBadRequest
		apiError.Code = string(domainerror.ValidationError)
		apiError.Message = ve.Error()
		apiError.Detail = make(map[string]any, len(ve))
		for k, v := range ve {
			apiError.Detail[k] = v.Error()
		}
	}

	return apiError
}

// WriteHTTPResponse marshals the given body and write it to given response writer, if the request
// has a query string named `fields` a json mask is applied to response
func WriteHTTPResponse(w http.ResponseWriter, r *http.Request, body any, status int) {
	if body == nil {
		w.WriteHeader(status)
		return
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		logger.Logger.WithError(err).Error("error marshaling json body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	if err != nil {
		logger.Logger.WithError(err).Error("error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// WriteHTTPError to response
func WriteHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.WithError(err).Error("unknown error")
	response := ToDomainErrorResponse(err)

	WriteHTTPResponse(w, r, response, response.StatusCode)
}

// WriteResponse to fiber context
func WriteResponse(c *fiber.Ctx, body any, status int) error {
	if body == nil {
		return c.SendStatus(status)
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		logger.Logger.WithError(err).Error("error marshaling json body")
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Context().Response.SetBodyRaw(bytes)
	c.Context().Response.Header.SetContentType("application/json")
	return c.SendStatus(status)
}

// WriteError to fiber context
func WriteError(c *fiber.Ctx, err error) error {
	logger.Logger.WithError(err).Error("unknown error")
	resp := ToDomainErrorResponse(err)
	return WriteResponse(c, resp, resp.StatusCode)
}

// ApiError data, all errors of our HTTP API returns this format
// @Description default error API format
type ApiError struct {
	// Code of the error
	Code string `json:"code"`
	// Message of the error
	Message string `json:"message"`
	// Detail of the error
	Detail map[string]any `json:"detail,omitempty"`
	// StatusCode
	StatusCode int `json:"-"`
}
