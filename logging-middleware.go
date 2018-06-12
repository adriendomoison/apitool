package apitool

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/adriendomoison/theholzbrothers/tools/gentool"
	"net/http"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func bodyReader(req *http.Request) *bytes.Buffer {
	if req.Body == nil {
		return nil
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil
	}
	buff := bytes.NewBuffer(body)
	req.Body = ioutil.NopCloser(buff)
	return bytes.NewBuffer(body)
}

func LoggingMiddleware(c *gin.Context) {

	c.Set("requestID", gentool.GenerateRandomString(24))

	// Decode Request Body
	m := make(map[string]interface{})
	json.NewDecoder(bodyReader(c.Request)).Decode(&m)

	// Log Request
	Logger.WithFields(logrus.Fields{
		"RequestID": c.GetString("requestID"),
		"Route":     c.Request.RequestURI,
		"Method":    c.Request.Method,
		"Body":      m,
	}).Info("Request")

	// Prepare Response Body
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	// Log Response
	Logger.WithFields(logrus.Fields{
		"RequestID":  c.GetString("requestID"),
		"StatusCode": c.Writer.Status(),
		"Body":       blw.body.String(),
	}).Info("Response")
}

