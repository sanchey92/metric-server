// Package middleware provides HTTP middleware handlers for common web server functionality.
package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(p []byte) (int, error) {
	return g.Writer.Write(p)
}

// GzipMiddleware is an HTTP middleware that handles gzip compression
// for both requests and responses.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request", http.StatusBadRequest)
				return
			}
			defer func() {
				if err = gr.Close(); err != nil {
					fmt.Println("failed to close gzip reader: %w", err)
				}
			}()
			r.Body = gr
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gzWriter := gzip.NewWriter(w)
		defer func() {
			if err := gzWriter.Close(); err != nil {
				fmt.Println("failed to close gzip writer: %w", err)
			}
		}()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gzWriter}, r)
	})
}
