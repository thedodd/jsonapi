package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
)

type errorInterfaceTester struct{}

func (e *errorInterfaceTester) Error() string               { return "Test Error." }
func (e *errorInterfaceTester) GetID() string               { return "Test ID." }
func (e *errorInterfaceTester) GetTitle() string            { return "Test Title." }
func (e *errorInterfaceTester) GetDetail() string           { return "Test Detail." }
func (e *errorInterfaceTester) GetStatus() string           { return "400" }
func (e *errorInterfaceTester) GetCode() string             { return "E1100" }
func (e *errorInterfaceTester) GetMeta() *map[string]string { return &(map[string]string{"key": "val"}) }

func TestErrorObjectWritesExpectedErrorMessage(t *testing.T) {
	err := &ErrorObject{Title: "Title test.", Detail: "Detail test."}
	var input error = err

	output := input.Error()

	if output != fmt.Sprintf("Error: %s %s\n", err.Title, err.Detail) {
		t.Fatal("Unexpected output.")
	}
}

var marshalErrorsTableTasts = []struct {
	In  []error
	Out map[string]interface{}
}{
	{ // This tests the insertion of reasonable placeholders for incompatible errors, and other fields are omitted.
		In: []error{errors.New("Testing.")},
		Out: map[string]interface{}{"errors": []interface{}{
			map[string]interface{}{"title": "Encountered error of type: *errors.errorString", "detail": "Testing."},
		}},
	},
	{ // This tests that given fields are turned into the appropriate JSON representation.
		In: []error{&ErrorObject{ID: "0", Title: "Test title.", Detail: "Test detail", Status: "400", Code: "E1100"}},
		Out: map[string]interface{}{"errors": []interface{}{
			map[string]interface{}{"id": "0", "title": "Test title.", "detail": "Test detail", "status": "400", "code": "E1100"},
		}},
	},
	{ // This tests that the `Meta` field is serialized properly.
		In: []error{&ErrorObject{Title: "Test title.", Detail: "Test detail", Meta: &map[string]string{"key": "val"}}},
		Out: map[string]interface{}{"errors": []interface{}{
			map[string]interface{}{"title": "Test title.", "detail": "Test detail", "meta": map[string]interface{}{"key": "val"}},
		}},
	},
}

func TestMarshalErrorsWritesTheExpectedPayload(t *testing.T) {
	for _, testRow := range marshalErrorsTableTasts {
		buffer, output := bytes.NewBuffer(nil), map[string]interface{}{}
		var writer io.Writer = buffer

		_ = MarshalErrors(writer, testRow.In)
		json.Unmarshal(buffer.Bytes(), &output)

		if !reflect.DeepEqual(output, testRow.Out) {
			t.Fatalf("Expected: \n%#v \nto equal: \n%#v", output, testRow.Out)
		}
	}
}

func TestMarshalErrorSerializesErrorAccordingToInterfaces(t *testing.T) {
	var err error = &errorInterfaceTester{}

	output := MarshalError(err)
	meta := *output.Meta
	val, ok := meta["key"]

	if output.ID != "Test ID." {
		t.Fatal("Unexpected value for error field: ID")
	}
	if output.Title != "Test Title." {
		t.Fatal("Unexpected value for error field: Title")
	}
	if output.Detail != "Test Detail." {
		t.Fatal("Unexpected value for error field: Detail")
	}
	if output.Status != "400" {
		t.Fatal("Unexpected value for error field: Status")
	}
	if output.Code != "E1100" {
		t.Fatal("Unexpected value for error field: Code")
	}
	if len(meta) != 1 || ok != true || val != "val" {
		t.Fatal("Unexpected value for error field: Meta")
	}
}

func TestMarshalErrorSerializesUsingFallbackApproachForIncompatibleErrors(t *testing.T) {
	err := errors.New("Testing fallback.")

	output := MarshalError(err)

	if output.ID != "" {
		t.Fatal("Unexpected value for error field: ID")
	}
	if output.Title != fmt.Sprintf("Encountered error of type: %T", err) {
		t.Fatal("Unexpected value for error field: Title")
	}
	if output.Detail != "Testing fallback." {
		t.Fatal("Unexpected value for error field: Detail")
	}
	if output.Status != "" {
		t.Fatal("Unexpected value for error field: Status")
	}
	if output.Code != "" {
		t.Fatal("Unexpected value for error field: Code")
	}
	if output.Meta != nil {
		t.Fatal("Unexpected value for error field: Meta")
	}
}
