package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

// compressWriter реализует ResponseWriter
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// newCompressWriter возвращает compressWriter
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header метод структуры compressWriter, который записывает данные в Header
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write метод структуры compressWriter, который записывает данные в тело http.ResponseWriter
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader метод структуры compressWriter, для записи StatusCode
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close метод структуры compressWriter, закрывает Writer
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует ReadCloser
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader возвращает compressReader
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read метод структуры newCompressReader, который читает тело запроса
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close метод структуры newCompressReader, закрывает структуру Close()
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
