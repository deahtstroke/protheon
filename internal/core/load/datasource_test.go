package load

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLoad_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	loader := SqlLoader{
		URL:   "someurl",
		Db:    db,
		Table: "test_table",
	}

	input := map[string]any{
		"hello":   "world!",
		"foo":     "bar",
		"howMany": 1,
	}

	mock.ExpectExec("INSERT INTO \"test_table\"").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = loader.Load(input)
	if err != nil {
		t.Fatalf("Error loading input to DB: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations for the mock were not met: %v", err)
	}
}

func TestLoad_ErrorInferrableTypes(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "Nil struct",
			input: nil,
		},
		{
			name: "Struct with nil fields",
			input: map[string]any{
				"Hello": nil,
				"World": nil,
			},
		},
		{
			name: "Struct with partial nil fields",
			input: map[string]any{
				"Foo": "bar",
				"":    nil,
			},
		},
		{
			name: "Fields with zero-string values",
			input: map[string]any{
				"": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, err := sqlmock.New()
			loader := SqlLoader{
				URL:   "someurl",
				Db:    db,
				Table: "test_table",
			}
			err = loader.Load(tt.input)
			if err == nil {
				t.Fatal("Expecting error, found none")
			}
		})
	}
}

func TestLoad_ErrorWhenCallingExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	table := "test_table"
	loader := SqlLoader{
		URL:   "someurl",
		Db:    db,
		Table: table,
	}

	mock.ExpectExec(fmt.Sprintf("INSERT INTO \"%s\"", table)).
		WithArgs("bar").
		WillReturnError(errors.New("Oopsie, doopsie. DB just exploded"))

	err = loader.Load(map[string]any{
		"foo": "bar",
	})

	if err == nil {
		t.Fatal("Expecting error, found none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %s", err)
	}
}
