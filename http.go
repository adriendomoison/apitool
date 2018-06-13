package apitool

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/mitchellh/mapstructure"
	"errors"
)

// RequestHeader is the object to send to the HttpRequestHandlers
type RequestHeader struct {
	URL           string
	Method        string
	ContentType   string
	Authorization string
}

// BuildErrorDescriptionFromApiError rebuild an ErrorDescription object out of an ApiError object
func BuildErrorDescriptionFromApiError(apiErrors ApiErrors, code int) (*ErrorDescription) {
	if code < 400 {
		return nil
	} else if len(apiErrors.Errors) > 0 {
		var err Error
		mapstructure.Decode(apiErrors.Errors[0], &err)
		return &ErrorDescription{
			Detail:  errors.New(err.Detail),
			Message: err.Message,
			Param:   err.Param,
			Code:    Code(code),
		}
	}
	if code == 500 {
		return &ErrorDescription{
			Detail: errors.New("the service is currently unavailable"),
			Message: "Something wrong happened on our end. Please try again later.",
			Code:   Code(code),
		}
	}
	return &ErrorDescription{
		Detail: errors.New("unknown error, the error was not explained by the service"),
		Message: "Something wrong happened on our end. Please try again later.",
		Code:   Code(code),
	}
}

// HttpRequestHandler is a handler for easy http requests
func HttpRequestHandler(requestHeader RequestHeader, requestBody interface{}, responseDTO interface{}) *ErrorDescription {
	req, err := http.NewRequest(requestHeader.Method, requestHeader.URL, encodeRequestBody(requestBody))
	req.Header.Set("Content-Type", requestHeader.ContentType)
	req.Header.Set("Authorization", requestHeader.Authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	var apiErrors ApiErrors
	json.Unmarshal(body, &responseDTO)
	json.Unmarshal(body, &apiErrors)

	return BuildErrorDescriptionFromApiError(apiErrors, resp.StatusCode)
}

func encodeRequestBody(reqBody interface{}) io.Reader {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	return b
}
