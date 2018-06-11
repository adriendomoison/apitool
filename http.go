package apitool

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// RequestHeader is the object to send to the HttpRequestHandlers
type RequestHeader struct {
	URL           string
	Method        string
	ContentType   string
	Authorization string
}

// HttpRequestHandler is a handler for easy http requests
func HttpRequestHandler(requestHeader RequestHeader, requestBody interface{}, responseDTO interface{}) (*http.Response, ApiErrors) {
	req, err := http.NewRequest(requestHeader.Method, requestHeader.URL, encodeRequestBody(requestBody))
	req.Header.Set("Content-Type", requestHeader.ContentType)
	req.Header.Set("Authorization", requestHeader.Authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	var apiErrors ApiErrors
	json.Unmarshal(body, &responseDTO)
	json.Unmarshal(body, &apiErrors)
	return resp, apiErrors
}

func encodeRequestBody(reqBody interface{}) io.Reader {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	return b
}
