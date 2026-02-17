package input

import (
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

type Extractor interface {
	Open() (io.ReadCloser, error)
}

type FileExtractor struct {
	Path     string
	Compress string
}

type combinedCloser struct {
	r io.Reader
	z io.Closer
}

func (e *FileExtractor) Open() (io.ReadCloser, error) {
	file, err := os.Open(e.Path)
	if err != nil {
		return nil, err
	}

	var r io.Reader = file
	if e.Compress == "zstd" {
		r, err = zstd.NewReader(r)
		if err != nil {
			return nil, err
		}

		return &combinedCloser{r: r, z: file}, nil
	}

	return file, nil
}

func (c *combinedCloser) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

func (c *combinedCloser) Close() error {
	return c.z.Close()
}
