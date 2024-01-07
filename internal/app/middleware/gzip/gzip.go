package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	encoding              = "gzip"
	headerContentEncoding = "Content-Encoding"
)

func Compress() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rw := w
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, encoding)
			if supportsGzip {
				cw := newCompressWriter(w)
				rw = cw
				defer func() {
					_ = cw.Close()
				}()
			}
			contentEncoding := r.Header.Get(headerContentEncoding)
			sendsGzip := strings.Contains(contentEncoding, encoding)
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer func() {
					_ = cr.Close()
				}()
			}
			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p) //nolint:wrapcheck // self-explanatory
}

func (c *compressWriter) WriteHeader(statusCode int) {
	// check if we have only successful status code
	if statusCode < http.StatusMultipleChoices {
		c.w.Header().Set(headerContentEncoding, encoding)
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close() //nolint:wrapcheck // self-explanatory
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("error while creating new gzip reader: %w", err)
	}
	return &compressReader{r: r, zr: zr}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p) //nolint:wrapcheck // self-explanatory
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("error while closing compress reader: %w", err)
	}
	return c.zr.Close() //nolint:wrapcheck // self-explanatory
}
