package logger

import (
	"io"
	"time"
)

type logWriter struct {
	out io.Writer
}

func (lw *logWriter) Write(b []byte) (int, error) {
	_, err := io.WriteString(lw.out, time.Now().UTC().Format("2006-01-02T15:04:05.999Z")+" ")
	if err != nil {
		return 0, err
	}
	return lw.out.Write(b)
}
