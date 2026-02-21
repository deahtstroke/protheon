package input

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

type Extractor interface {
	Open() (io.ReadCloser, error)
}

type FileExtractor struct {
	Path     string
	Compress Compress
}

type combinedCloser struct {
	Reader io.Reader
	Closer io.Closer
}

type Compress int

const (
	Zstd Compress = iota
	Gzip
)

type FileExtractorOpt func(*FileExtractor)

func NewFileExtractor(filePath string, opts ...FileExtractorOpt) (*FileExtractor, error) {
	fe := &FileExtractor{
		Path: filePath,
	}

	if _, err := os.Stat(filePath); err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("Resource with path %s does not exist", filePath)
	}

	for _, opt := range opts {
		if opt != nil {
			opt(fe)
		}
	}

	return fe, nil
}

func (e *FileExtractor) Open() (io.ReadCloser, error) {
	file, err := os.Open(e.Path)
	if err != nil {
		return nil, err
	}

	var r io.Reader = file

	switch e.Compress {
	case Zstd:
		r, err = zstd.NewReader(r)
		if err != nil {
			return nil, err
		}
		return &combinedCloser{Reader: r, Closer: file}, nil
	case Gzip:
		r, err = gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		return &combinedCloser{Reader: r, Closer: file}, nil
	default:
	}

	return file, nil
}

func (c *combinedCloser) Read(p []byte) (int, error) {
	return c.Reader.Read(p)
}

func (c *combinedCloser) Close() error {
	return c.Closer.Close()
}

func WithGzip() FileExtractorOpt {
	return func(fe *FileExtractor) {
		fe.Compress = Gzip
	}
}

func WithZstd() FileExtractorOpt {
	return func(fe *FileExtractor) {
		fe.Compress = Zstd
	}
}
