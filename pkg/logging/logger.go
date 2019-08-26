package logging

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var Logger = newLogger()

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func newLogger() *StandardLogger {
	var baseLogger = logrus.New()

	var standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}

	standardLogger.SetOutput(os.Stdout)

	// Info or above
	standardLogger.SetLevel(logrus.InfoLevel)

	return standardLogger
}

// Create a logger entry and add the fields method, traceId and requestId from the http request object
func GetLoggerWithFields(r *http.Request) *logrus.Entry {
	logger := Logger.WithFields(logrus.Fields{
		"method": r.Method,
		//"traceId":   r.Header.Get("x-trace-id"),
		//"requestId": r.Header.Get("x-request-id"),
	})
	return logger
}

// accessWriter is a simple wrapper that helps us capture the http response status and content-length
type accessWriter struct {
	http.ResponseWriter
	responseStatus int
	contentLength  int
}

// Write status header and capture value
func (w *accessWriter) WriteHeader(status int) {
	w.responseStatus = status
	w.ResponseWriter.WriteHeader(status)
}

// Write response and capture content length
func (w *accessWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)

	w.contentLength += n
	return n, err
}

// Define a middleware to log all the requests handled by the service
func AccessLoggingMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		path := r.URL.Path

		responseWriter := &accessWriter{
			ResponseWriter: w,
		}

		defer func() {
			latency := float64(time.Now().UTC().Sub(start).Nanoseconds()) / 1e6

			logger := GetLoggerWithFields(r)

			info := logrus.Fields{
				"path":           path,
				"ip":             r.RemoteAddr,
				"duration":       latency,
				"user_agent":     r.Header.Get("User-Agent"),
				"content_length": responseWriter.contentLength,
				"status":         responseWriter.responseStatus,
			}
			logger = logger.WithFields(info)

			logger.Info("access")
		}()

		nextHandler.ServeHTTP(responseWriter, r)
	})
}
