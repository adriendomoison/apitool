package apitool

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// WaitForServerToStart return true only when API is ready
func WaitForServerToStart(url string) bool {
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		client := &http.Client{}
		if _, err := client.Do(req); err == nil {
			return true
		}
		time.Sleep(2500)
	}
	return false
}

// HttpRequestHandlerForUnitTesting is the same as HttpRequestHandler but add logging for tests
func HttpRequestHandlerForUnitTesting(t *testing.T, requestHeader RequestHeader, requestBody interface{}, responseDTO interface{}) (*http.Response, ApiErrors) {

	// Log header
	t.Log(requestHeader)

	// Create request
	req, err := http.NewRequest(requestHeader.Method, requestHeader.URL, encodeRequestBodyAndLog(t, requestBody))
	req.Header.Set("Content-Type", requestHeader.ContentType)
	req.Header.Set("Authorization", requestHeader.Authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read response
	body, _ := ioutil.ReadAll(resp.Body)

	apiErrors := ApiErrors{}
	json.Unmarshal(body, &responseDTO)
	json.Unmarshal(body, &apiErrors)

	// Log response
	t.Log(responseDTO)
	t.Log(apiErrors)
	return resp, apiErrors
}

func encodeRequestBodyAndLog(t *testing.T, reqBody interface{}) io.Reader {
	t.Log("testing with following parameters:")
	t.Log(reqBody)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	return b
}
