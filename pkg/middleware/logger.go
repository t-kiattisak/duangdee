package middleware

import (
	"bytes"

	"duangdee/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// responseBodyWriter wraps Fiber response to capture HTTP status and body
type responseBodyWriter struct {
	fiber.Ctx
	body *bytes.Buffer
}

// Logger returns a Fiber middleware that logs detailed HTTP request/response metrics to Kibana
func Logger(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Generate or extract Trace ID / X-Request-ID for distributed tracing
		traceID := c.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("X-Request-ID", traceID)

		// 2. Extract Request Headers and Payload (Limit max 4KB to save disk space)
		reqHeaders := make(map[string]string)
		c.Request().Header.VisitAll(func(k, v []byte) {
			reqHeaders[string(k)] = string(v)
		})

		reqBody := string(c.Body())
		if len(reqBody) > 4096 {
			reqBody = reqBody[:4096] + "...[truncated]"
		}

		// 3. Process the actual HTTP Request
		err := c.Next()

		// 4. Capture Response Status and Body
		resBody := string(c.Response().Body())
		if len(resBody) > 4096 {
			resBody = resBody[:4096] + "...[truncated]"
		}

		status := c.Response().StatusCode()

		// 5. Construct Structured Metadata for Kibana Indexing
		meta := map[string]interface{}{
			"trace_id":         traceID,
			"http_method":      c.Method(),
			"http_path":        c.Path(),
			"http_status":      status,
			"client_ip":        c.IP(),
			"user_agent":       c.Get("User-Agent"),
			"request_headers":  reqHeaders,
			"request_payload":  reqBody,
			"response_payload": resBody,
			"query_params":     c.Queries(),
		}

		if status >= 500 {
			log.Error(err, "HTTP Server Error Response", meta)
		} else if status >= 400 {
			log.Debug("HTTP Client Error Response", meta)
		} else {
			log.Info("HTTP Request Processed", meta)
		}

		return err
	}
}
