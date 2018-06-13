// Package apihelper automate construction of http response
package apitool

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ApiError interface for all API error messages
type ApiError interface {
}

// Error is the default error message structure the API returns
type Error struct {
	apiError ApiError
	Param    string `json:"param"`
	Detail   string `json:"detail"`
	Message  string `json:"message"`
}

// ApiErrors carry the list of errors returned by the API from a request
type ApiErrors struct {
	Errors []ApiError `json:"errors"`
}

// Code describe the status that will be generated in the BuildResponseError for the http response
type Code uint

// Code describe the status that will be generated in the BuildResponseError for the http response
const (
	Processing      Code = 102
	NoContent            = 204
	BadRequest           = 400
	Unauthorized         = 401
	Forbidden            = 403
	NotFound             = 404
	AlreadyExist         = 409
	UnexpectedError      = 500
	NotImplemented       = 501
)

// Error describe the error object returned from services that can be passed directly to the BuildResponseError method
type ErrorDescription struct {
	Param   string
	Detail  error
	Message string
	Code    Code
}

// DefaultCORSConfig Generate CORS config for router
func DefaultCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowOrigins = []string{os.Getenv("WHITELISTED_DOMAIN")}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	CORSConfig.AddAllowHeaders("Authorization")
	return CORSConfig
}

// BuildRequestError build a usable JSON error object from an error string generated by the structure validator
func BuildRequestError(err error) (int, ApiErrors) {
	var apiErrors ApiErrors
	switch err.(type) {
	case validator.ValidationErrors:
		for _, v := range err.(validator.ValidationErrors) {
			var validationError Error
			validationError.Param = toSnakeCase(v.Field)
			validationError.Detail = "Field validation for " + toSnakeCase(v.Field) + " failed on the " + v.Tag + " tag."
			if v.Tag == "required" {
				validationError.Message = "This field is required"
			}
			if v.Tag == "email" {
				validationError.Message = "Invalid email address. Valid email can contain only letters, numbers, '@' and '.'"
			}
			if v.Tag == "url" {
				validationError.Message = "Invalid URL address. Valid URL start with http:// or https://"
			}
			if v.Tag == "min" {
				validationError.Message = v.Name + " need to be at least " + v.Param + " characters long."
			}
			if v.Tag == "max" {
				validationError.Message = v.Name + " must be less than " + v.Param + " characters long."
			}
			apiErrors.Errors = append(apiErrors.Errors, validationError)
		}
		return http.StatusBadRequest, apiErrors
	default:
		apiErrors.Errors = append(apiErrors.Errors, Error{
			Detail: err.Error(),
		})
		return http.StatusBadRequest, apiErrors
	}
}

// BuildResponseError apply the right status to the http response and build the error JSON object
func BuildResponseError(err *ErrorDescription) (status int, apiErrors ApiErrors) {
	if err.Code >= 400 {
		if err.Detail == nil {
			err.Detail = errors.New(err.Message)
		}
		apiErrors.Errors = append(apiErrors.Errors, Error{
			Detail:  err.Detail.Error(),
			Message: err.Message,
			Param:   err.Param,
		})
	}
	return int(err.Code), apiErrors
}

// GetBoolQueryParam allow to retrieve a boolean query parameter.
// It takes the gin context as param to build error if queryParam is not formatted correctly and a default value to set value if the parameter is optional and not set.
func GetBoolQueryParam(c *gin.Context, value *bool, queryParam string, defaultValue bool) bool {
	var err error
	if c.Query(queryParam) != "" {
		if *value, err = strconv.ParseBool(c.Query(queryParam)); err != nil {
			c.JSON(BuildRequestError(
				errors.New("query parameter '" + queryParam + "' value should be true or false (omit the key and default value '" + strconv.FormatBool(defaultValue) + "' will be applied)")),
			)
			return false
		}
	} else {
		*value = defaultValue
	}
	return true
}

// GetStringQueryParam allow to retrieve a string query parameter.
// It takes the gin context as param to build error if queryParam is not formatted correctly and a default value to set value if the parameter is optional and not set.
// the accepted values list will define if the value is valid or not
func GetStringQueryParam(c *gin.Context, queryParam string, acceptedValues []string, defaultValue string) (string, bool) {
	param := c.Query(queryParam)
	if param != "" {
		sort.Strings(acceptedValues)
		i := sort.SearchStrings(acceptedValues, param)
		if i < len(acceptedValues) && acceptedValues[i] == param {
			return param, true
		}
		c.JSON(BuildRequestError(
			errors.New("query parameter '" + queryParam + "' value should be a date with the format " + strings.Join(acceptedValues, ", ") + "(omit the key and default value '" + defaultValue + "' will be applied)"),
		))
		return "", false
	}
	return defaultValue, true
}

// GetStringQueryParam allow to retrieve a string query parameter.
// It takes the gin context as param to build error if queryParam is not formatted correctly and a default value to set value if the parameter is optional and not set.
// the accepted values list will define if the value is valid or not
func GetMandatoryIntQueryParam(c *gin.Context, queryParam string) (int, bool) {
	param := c.Query(queryParam)
	if param == "" {
		return 0, false
	}
	if value, err := strconv.Atoi(param); err == nil {
		return value, true
	}
	c.JSON(BuildRequestError(
		errors.New("query parameter '" + queryParam + "' value should be an integer (this value is mandatory)"),
	))
	return 0, false
}

// GetStringQueryParam allow to retrieve a string query parameter.
// It takes the gin context as param to build error if queryParam is not formatted correctly and a default value to set value if the parameter is optional and not set.
// the accepted values list will define if the value is valid or not
func GetMandatoryUintQueryParam(c *gin.Context, queryParam string) (uint, bool) {
	param := c.Query(queryParam)
	if param == "" {
		return 0, false
	}
	if value, err := strconv.Atoi(param); err == nil {
		return uint(value), true
	}
	c.JSON(BuildRequestError(
		errors.New("query parameter '" + queryParam + "' value should be an integer (this value is mandatory)"),
	))
	return 0, false
}

// GetDateQueryParam allow to retrieve a date (as a string) query parameter.
// It takes the gin context as param to build error if queryParam is not formatted correctly and a default value to set value if the parameter is optional and not set.
// The dateLayout will define if the date is valid or not
func GetDateQueryParam(c *gin.Context, queryParam string, dateLayout string, defaultValue string) (string, bool) {
	date := c.Query(queryParam)
	if date != "" {
		if _, err := time.Parse(dateLayout, date); err != nil {
			c.JSON(BuildRequestError(
				errors.New("query parameter '" + queryParam + "' value should be a date with the format " + dateLayout + "(omit the key and default value '" + defaultValue + "' will be applied)"),
			))
			return "", false
		}
		return date, true
	}
	return defaultValue, true
}

// toSnakeCase change a string to it's snake case version
func toSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
