package input

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestOpen_NoOpts(t *testing.T) {
	tests := []struct {
		filePath string
		opt      FileExtractorOpt
	}{
		{
			filePath: "testdata/example.json",
		},
		{
			filePath: "testdata/example.jsonl",
		},
		{
			filePath: "testdata/example.jsonl.gzip",
			opt:      WithGzip(),
		},
		{
			filePath: "testdata/example.jsonl.zst",
			opt:      WithZstd(),
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing: %s", filepath.Ext(tt.filePath)), func(t *testing.T) {
			sut, err := NewFileExtractor(tt.filePath, tt.opt)
			if err != nil {
				t.Fatalf("Error while initializing extractor: %v", err)
			}

			r, err := sut.Open()
			if err != nil {
				t.Fatalf("Error while calling open: %s", err)
			}

			if _, ok := r.(*combinedCloser); !ok {
				t.Fatalf("Expecting combined closer")
			}
		})
	}
}

func TestOpen_ErrorFileNotFound(t *testing.T) {
	path := "./testdata/this_doesnt_exist.json"
	_, err := NewFileExtractor(path)
	if err == nil {
		t.Fatal("Expecting error, found none")
	}
}
