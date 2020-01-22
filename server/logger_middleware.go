package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type loggerMiddleware struct {
	logger *log.Logger
	next   http.Handler
}

func NewLoggerMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return &loggerMiddleware{
		logger: logger,
		next:   next,
	}
}

func (m *loggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string]string, len(r.Header))
	for name, values := range r.Header {
		headers[name] = strings.Join(values, "; ")
	}

	body := m.drainBody(r)

	entry := logEntry{
		Method:  r.Method,
		URI:     r.URL.RequestURI(),
		Headers: headers,
		Body:    body.String(),
	}
	line, err := json.Marshal(entry)
	if err != nil {
		m.logger.Fatalln(err)
	}
	m.logger.Println(string(line))

	m.next.ServeHTTP(w, r)
}

func (m *loggerMiddleware) drainBody(r *http.Request) *bytes.Buffer {
	var body bytes.Buffer
	if _, err := body.ReadFrom(r.Body); err != nil {
		m.logger.Fatalln(err)
	}
	if err := r.Body.Close(); err != nil {
		m.logger.Fatalln(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

	return &body
}
