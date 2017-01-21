package jsonapi

import (
	"encoding/json"
	"fmt"
	"io"
)

// MarshalErrors writes a JSON API response using the given `[]error`.
// This function works by passing each `error` in the given slice through a routine which attempts
// to cast the error to each of this package's `Error<FieldName>Compatible` interfaces. Though
// it works for errors which do not implement any of these interfaces as well. For more information
// on JSON API error payloads, see the spec here: http://jsonapi.org/format/#document-top-level
// and here: http://jsonapi.org/format/#error-objects
//
// To use this function in your code effectively, make your package's error types implement some
// or all of the error compatibility interfaces (`ErrorTitleCompatible`, `ErrorDetailCompatible`,
// ...). Doing so gives you direct control over what will appear in the serialied error object's
// fields. Without implementing these interfaces, the `title` & `detail` fields will be populated
// with derived data about the error.
func MarshalErrors(w io.Writer, errs []error) error {
	// Serialize the given errors.
	var formattedErrors []ErrorObject
	for _, err := range errs {
		e := marshalError(err)
		formattedErrors = append(formattedErrors, e)
	}

	// Write out the serialize errors payload.
	if err := json.NewEncoder(w).Encode(&ErrorsPayload{Errors: formattedErrors}); err != nil {
		return err
	}
	return nil
}

// marshalError will serialize the given error as best as possible according to this
// package's `Error<field>Compatible` interfaces.
func marshalError(err error) ErrorObject {
	errorObject := ErrorObject{}
	if e, ok := err.(ErrorIDCompatible); ok {
		errorObject.ID = e.GetID()
	}

	if e, ok := err.(ErrorTitleCompatible); ok {
		errorObject.Title = e.GetTitle()
	} else {
		errorObject.Title = fmt.Sprintf("Encountered error of type: %T", err)
	}

	if e, ok := err.(ErrorDetailCompatible); ok {
		errorObject.Detail = e.GetDetail()
	} else {
		errorObject.Detail = err.Error()
	}

	if e, ok := err.(ErrorStatusCompatible); ok {
		errorObject.Status = e.GetStatus()
	}

	if e, ok := err.(ErrorCodeCompatible); ok {
		errorObject.Code = e.GetCode()
	}

	if e, ok := err.(ErrorMetaCompatible); ok {
		errorObject.Meta = e.GetMeta()
	}

	return errorObject
}

// ErrorsPayload is a serializer struct for representing a valid JSON API errors payload.
type ErrorsPayload struct {
	Errors []ErrorObject `json:"errors"`
}

// ErrorObject is an `Error` implementation as well as an implementation of the JSON API error object.
// The main idea behind this struct is that you
// For more information on Golang errors, see: https://golang.org/pkg/errors/
// For more information on the JSON API spec's error objects, see: http://jsonapi.org/format/#error-objects
type ErrorObject struct {
	// ID is a unique identifier for this particular occurrence of a problem.
	ID string `json:"id,omitempty"`

	// Title is a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization.
	Title string `json:"title"`

	// Detail is a human-readable explanation specific to this occurrence of the problem. Like title, this field’s value can be localized.
	Detail string `json:"detail"`

	// Status is the HTTP status code applicable to this problem, expressed as a string value.
	Status string `json:"status,omitempty"`

	// Code is an application-specific error code, expressed as a string value.
	Code string `json:"code,omitempty"`

	// TODO: (thedodd): add this when we have an internal model to use.
	// Links is an array of link objects containing hyper-links to further details about
	// this particular occurrence of the problem.
	// Links []*Link `json:"links,omitempty"`

	// TODO: (thedodd): add this when we have an internal model to use.
	// Source is an object containing references to the source of the error.
	// Source *Source `json:"source,omitempty"`

	// Meta is an object containing non-standard meta-information about the error.
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

// Error implements the `Error` interface.
func (e *ErrorObject) Error() string {
	return fmt.Sprintf("Error: %s %s\n", e.Title, e.Detail)
}

// GetID implements the `ErrorIDCompatible` interface.
func (e *ErrorObject) GetID() string { return e.ID }

// GetTitle implements the `ErrorTitleCompatible` interface.
func (e *ErrorObject) GetTitle() string { return e.Title }

// GetDetail implements the `ErrorDetailCompatible` interface.
func (e *ErrorObject) GetDetail() string { return e.Detail }

// GetStatus implements the `ErrorStatusCompatible` interface.
func (e *ErrorObject) GetStatus() string { return e.Status }

// GetCode implements the `ErrorCodeCompatible` interface.
func (e *ErrorObject) GetCode() string { return e.Code }

// GetMeta implements the `ErrorMetaCompatible` interface.
func (e *ErrorObject) GetMeta() *map[string]interface{} { return e.Meta }

/////////////////////////////////////////////
// JSON API Error Compatibility Interfaces //
/////////////////////////////////////////////

// ErrorIDCompatible is the interface needed for exposing the `id` field of a JSON API compatible error.
type ErrorIDCompatible interface {
	GetID() string
}

// ErrorTitleCompatible is the interface needed for exposing the `title` field of a JSON API compatible error.
type ErrorTitleCompatible interface {
	GetTitle() string
}

// ErrorDetailCompatible is the interface needed for exposing the `detail` field of a JSON API compatible error.
type ErrorDetailCompatible interface {
	GetDetail() string
}

// ErrorStatusCompatible is the interface needed for exposing the `status` field of a JSON API compatible error.
type ErrorStatusCompatible interface {
	GetStatus() string
}

// ErrorCodeCompatible is the interface needed for exposing the `code` field of a JSON API compatible error.
type ErrorCodeCompatible interface {
	GetCode() string
}

// ErrorMetaCompatible is the interface needed for exposing the `meta` field of a JSON API compatible error.
type ErrorMetaCompatible interface {
	GetMeta() *map[string]interface{}
}
