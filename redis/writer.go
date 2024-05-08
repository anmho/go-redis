package redis

import (
	"fmt"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	_, err := w.writer.Write(v.Marshal())
	if err != nil {
		return fmt.Errorf("writing input: %w", err)
	}
	return nil
}
