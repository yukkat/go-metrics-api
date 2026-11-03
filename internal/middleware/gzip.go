package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"mime"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w           http.ResponseWriter
	zw          *gzip.Writer
	wroteHeader bool
	compress    bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
	}
}

func isCompressible(contentType string) bool {
    mediaType, _, err := mime.ParseMediaType(contentType)
    if err != nil {
        return false
    }

    return mediaType == "application/json" ||
        mediaType == "text/html"
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
    if !c.wroteHeader {
        if c.w.Header().Get("Content-Type") == "" {
            c.w.Header().Set("Content-Type", http.DetectContentType(p))
        }

        c.WriteHeader(http.StatusOK)
    }

    if !c.compress {
        return c.w.Write(p)
    }

    return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if c.wroteHeader {
		return
	}
	c.wroteHeader = true

	noBody := statusCode < http.StatusOK ||
		statusCode == http.StatusNoContent ||
		statusCode == http.StatusResetContent ||
		statusCode == http.StatusNotModified

	contentType := c.w.Header().Get("Content-Type")

	c.compress = !noBody && isCompressible(contentType)
	if c.compress {
        c.zw = gzip.NewWriter(c.w)
        c.w.Header().Set("Content-Encoding", "gzip")
    }

	c.w.WriteHeader(statusCode)
}

// Flush сбрасывает буфер gzip и вызывает Flush у нижележащего ResponseWriter.
func (c *compressWriter) Flush() {
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}

	if c.compress {
		_ = c.zw.Flush()
	}

	if f, ok := c.w.(http.Flusher); ok {
		f.Flush()
	}
}

// ReadFrom позволяет сжимать данные при io.Copy без лишних промежуточных буферов.
func (c *compressWriter) ReadFrom(r io.Reader) (int64, error) {
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}
	if !c.compress {
		return io.Copy(c.w, r)
	}
	return io.Copy(c.zw, r)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	if !c.compress {
		return nil
	}
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

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

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	})
}
