package input

import (
	"bufio"
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

	var record map[string]any
	err = json.Unmarshal(line, &record)
	return record, err
}
