package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

// compress Writer implements ResponseWriter
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// newCompress Writer returns compressWriter
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header a method of the compress Writer structure that writes data to the Header
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write a method of the compressWriter structure that writes data to the http.ResponseWriter
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader method of the compressWriter structure, for writing statusCode
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close method of the compressWriter structure, closes the Writer
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader implements ReadCloser
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader returns compressReader
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

// Read a method of the newCompressReader structure that reads the request body
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close method of the newCompressReader structure, closes the Close() structure
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
