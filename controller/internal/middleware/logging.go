package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type requestLog struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Query     string `json:"query,omitempty"`
	Status    int    `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
	Error     string `json:"error,omitempty"`
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()

		entry := requestLog{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			RequestID: requestID,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Query:     c.Request.URL.RawQuery,
			Status:    c.Writer.Status(),
			LatencyMS: time.Since(start).Milliseconds(),
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}

		if len(c.Errors) > 0 {
			entry.Error = c.Errors.String()
		}

		data, err := json.Marshal(entry)
		if err != nil {
			log.Printf(`{"level":"error","message":"failed to marshal request log","error":"%v"}`, err)
			return
		}

		log.Print(string(data))
	}
}

func newRequestID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}
