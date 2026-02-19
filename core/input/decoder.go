package input

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

type Decoder interface {
	Next() (map[string]any, error)
}

type JSONLDecoder struct {
	bufReader    *bufio.Reader
	readerCloser io.ReadCloser
}

func NewJSONLDecoder(readCloser io.ReadCloser) *JSONLDecoder {
	return &JSONLDecoder{
		bufReader:    bufio.NewReader(readCloser),
		readerCloser: readCloser,
	}
}

func (d *JSONLDecoder) Next() (map[string]any, error) {
	line, err := d.bufReader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	if len(line) == 0 {
		if err == io.EOF {
			return nil, io.EOF
		}
	}

	var record map[string]any
	err = json.Unmarshal(bytes.TrimSpace(line), &record)
	return record, err
}
