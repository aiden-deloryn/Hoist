package util

import (
	"bufio"
)

type ProgressReader struct {
	bufio.Reader
	bytesCopied      int64
	ProgressCallback func(bytesCopied int64)
}

func (this *ProgressReader) Read(p []byte) (int, error) {
	byteCount, err := this.Reader.Read(p)

	if err == nil {
		this.bytesCopied += int64(byteCount)
		this.ProgressCallback(this.bytesCopied)
	}

	return byteCount, err
}
