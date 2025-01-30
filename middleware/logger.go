package cMiddleware

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type CustomResponseWriter struct {
	http.ResponseWriter               // Embed original ResponseWriter
	body                *bytes.Buffer // Ready to capture body and status
	status              int
}

func (c *CustomResponseWriter) Write(body []byte) (int, error) {
	c.body.Write(body)
	return c.ResponseWriter.Write(body)
}

func (c *CustomResponseWriter) WriteHeader(statusCode int) {
	c.status = statusCode // We take and we write like normal WriteHeader
	c.ResponseWriter.WriteHeader(statusCode)
}

func ResponseLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		crw := &CustomResponseWriter{
			ResponseWriter: c.Response().Writer,
			body:           bytes.NewBufferString(""),
			status:         http.StatusOK,
		}

		c.Response().Writer = crw

		err := next(c)

		clientIP := c.RealIP()
		method := c.Request().Method
		path := c.Path()
		status := crw.status

		responseBody := crw.body.String()

		if isTemplate(responseBody) {
			responseBody = "<!DOCTYPE html>"
		}

		log.Printf("Client IP: %s | %s %d %s | Response: %s\n", clientIP, method, status, path, responseBody)

		return err
	}
}

func isTemplate(respose string) bool {
	label := "<!DOCTYPE html>"
	ptrLabel := 0
	resposeLen := len(respose)
	if resposeLen == 0 {
		return false
	}

	for i := 0; i < resposeLen/10; i++ {
		if respose[i] == label[ptrLabel] {
			ptrLabel++
			if ptrLabel == len(label) {
				return true
			}
			continue
		}
		ptrLabel = 0
	}
	return false
}
