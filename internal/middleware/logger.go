package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp            string        `json:"timestamp"`
	RequestMethod        string        `json:"http.request.method"`
	RequestPath          string        `json:"http.route"`
	RequestHost          string        `json:"http.request.host"`
	RequestRemoteAddr    string        `json:"http.request.remote_addr"`
	ResponseStatus       int           `json:"http.response.status_code"`
	ResponseDuration     int64         `json:"http.server.request.duration"`
	LogLevel             string        `json:"http.log.level"`
	Message              string        `json:"http.request.message"`
	UserAgent            string        `json:"http.user_agent,omitempty"`
	RequestID            string        `json:"request_id,omitempty"`
	Error                string        `json:"error,omitempty"`
}

func JSONLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = c.GetHeader("X-Request-ID")
		}

		c.Next()

		duration := time.Since(startTime).Milliseconds()

		logLevel := "info"
		if c.Writer.Status() >= 400 {
			logLevel = "error"
		} else if c.Writer.Status() >= 300 {
			logLevel = "warning"
		}

		logEntry := LogEntry{
			Timestamp:         time.Now().Format(time.RFC3339Nano),
			RequestMethod:     c.Request.Method,
			RequestPath:       c.Request.URL.Path,
			RequestHost:       c.Request.Host,
			RequestRemoteAddr: c.ClientIP(),
			ResponseStatus:    c.Writer.Status(),
			ResponseDuration:  duration,
			LogLevel:          logLevel,
			Message:           "Incoming request:",
			UserAgent:         c.Request.UserAgent(),
			RequestID:         requestID,
		}

		if len(c.Errors) > 0 {
			logEntry.Error = c.Errors.String()
		}

		logBytes, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("failed to marshal log entry: %v", err)
			return
		}

		log.Printf("%s", string(logBytes))
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}