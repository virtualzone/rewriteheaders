package rewriteheaders

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// Rewrite holds one rewrite body configuration.
type Rewrite struct {
	Header      string `json:"header,omitempty"`
	Regex       string `json:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	Rewrites []Rewrite `json:"rewrites,omitempty"`
}

type ResponseWriter struct {
	Writer   http.ResponseWriter
	Rewrites []rewrite
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type rewrite struct {
	header      string
	regex       *regexp.Regexp
	replacement []byte
}

type rewriteHeader struct {
	name     string
	next     http.Handler
	rewrites []rewrite
}

// New creates and returns a new rewrite body plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	rewrites := make([]rewrite, len(config.Rewrites))
	for i, rewriteConfig := range config.Rewrites {
		regex, err := regexp.Compile(rewriteConfig.Regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", rewriteConfig.Regex, err)
		}
		rewrites[i] = rewrite{
			header:      rewriteConfig.Header,
			regex:       regex,
			replacement: []byte(rewriteConfig.Replacement),
		}
	}

	return &rewriteHeader{
		name:     name,
		next:     next,
		rewrites: rewrites,
	}, nil
}

func (r *rewriteHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedWriter := &ResponseWriter{
		Writer:   rw,
		Rewrites: r.rewrites,
	}
	r.next.ServeHTTP(wrappedWriter, req)
}

func (r *ResponseWriter) Header() http.Header {
	return r.Writer.Header()
}

func (r *ResponseWriter) Write(bytes []byte) (int, error) {
	return r.Writer.Write(bytes)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	for _, rewrite := range r.Rewrites {
		headerValues := r.Writer.Header().Values(rewrite.header)

		if len(headerValues) == 0 {
			continue
		}

		r.Writer.Header().Del(rewrite.header)

		for _, header := range headerValues {
			value := rewrite.regex.ReplaceAll([]byte(header), rewrite.replacement)
			r.Writer.Header().Add(rewrite.header, string(value))
		}
	}

	r.Writer.WriteHeader(statusCode)
}
