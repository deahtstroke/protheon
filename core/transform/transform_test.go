package transform

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestTransform_SuccessSimpleTypes(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"a":   1,
		"b":   2,
	}
	transformer := NewLuaTransformer("./testdata/simple_fields.lua")
	output, err := transformer.Transform(input)
	if err != nil {
		t.Fatalf("Error while applying transformations to input: %s", err)
	}
	m, ok := output.(map[string]any)
	if !ok {
		t.Fatalf("Unable to parse the output as a generic map")
	}
	assert.Check(t, m["sum"] == float64(3))
	assert.Check(t, m["foo"] == "bar Hello World")
}

func TestTransform_SuccessStructsAndMaps(t *testing.T) {
	type sum struct {
		A int
		B int
	}
	input := map[string]any{
		"sum": sum{
			A: 5,
			B: 10,
		},
		"hello world": map[string]any{
			"foo": "hello",
			"bar": "world!",
		},
		"sums": []sum{
			{
				A: 3,
				B: 3,
			},
			{
				A: 5,
				B: 5,
			},
		},
	}

	transformer := NewLuaTransformer("./testdata/complex_fields.lua")
	output, err := transformer.Transform(input)
	if err != nil {
		t.Fatalf("Error while applying transformations to input: %v", err)
	}

	m, ok := output.(map[string]any)
	if !ok {
		t.Fatalf("Unable to parse the output as a generic map")
	}

	assert.Check(t, m["sum"] == float64(15))
	assert.Check(t, m["hello world"] == "hello world!")
	assert.Check(t, m["total"] == float64(16))
}
