package redis

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}
	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go func() {
		for {
			aof.mu.Lock()
			err := aof.file.Sync()
			if err != nil {
				aof.mu.Unlock()
				return
			}
			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()
	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())

	if err != nil {
		return err
	}
	return nil
}

func (aof *Aof) Read() error {
	aof.file.Seek(0, io.SeekStart)

	reader := NewResp(aof.file)
	for {
		req, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("reading aof: %w", err)
		}

		// we should execute it so we can load it into memory
		_, err = Handle(req)
		if err != nil {
			return fmt.Errorf("loading from append-only file: %w", err)
		}
	}
	return nil
}
