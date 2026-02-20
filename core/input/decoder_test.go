package input

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

type mockReadCloser struct {
	jsonResult string
	off        int
	shouldErr  bool
	shouldEof  bool
}

func (m *mockReadCloser) Close() error {
	if m.shouldErr {
		return errors.New("Error closing")
	}

	if m.shouldEof {
		return io.EOF
	}
	return nil
}

func (m *mockReadCloser) Read(b []byte) (n int, err error) {
	if m.shouldErr {
		return 0, errors.New("Something happened")
	}

	n = copy(b, []byte(m.jsonResult))
	m.off += n
	return n, nil
}

func TestNext_Success(t *testing.T) {
	mock := mockReadCloser{
		jsonResult: "{ \"hello\": \"world!\" }\n",
		shouldErr:  false,
		shouldEof:  false,
	}

	sut := NewJSONLDecoder(&mock)
	res, err := sut.Next()
	if err != nil {
		t.Fatal(err)
	}

	if res["hello"] != "world!" {
		t.Fatalf("Error getting the world of the hello :(: %v", res)
	}
}

func TestNext_ErrorReadingBytes(t *testing.T) {
	mock := mockReadCloser{
		shouldErr: true,
	}

	sut := NewJSONLDecoder(&mock)
	_, err := sut.Next()
	if err == nil {
		t.Fatal("Expecting error, found none")
	}
}

func TestNext_NoErrorOnNewLinesOrEmptyLines(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: " \n",
		},
		{
			input: " ",
		},
		{
			input: "",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test for input %s", test.input), func(t *testing.T) {
			mock := mockReadCloser{
				jsonResult: " \n",
			}

			sut := NewJSONLDecoder(&mock)
			obj, err := sut.Next()
			if err != nil {
				t.Fatalf("Not expecting error, found: %v", err)
			}

			if len(obj) != 0 {
				t.Fatal("Resulting object should have a length of zero")
			}
		})
	}
}
