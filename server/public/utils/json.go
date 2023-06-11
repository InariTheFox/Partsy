package utils

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

type HumanizedJSONError struct {
	Err       error
	Line      int
	Character int
}

func (e *HumanizedJSONError) Error() string {
	return e.Err.Error()
}

func HumanizeJSONError(err error, data []byte) error {
	if syntaxError, ok := err.(*json.SyntaxError); ok {
		return NewHumanizedJSONError(syntaxError, data, syntaxError.Offset)
	} else if unmarshalError, ok := err.(*json.UnmarshalTypeError); ok {
		return NewHumanizedJSONError(unmarshalError, data, unmarshalError.Offset)
	} else {
		return err
	}
}

func NewHumanizedJSONError(err error, data []byte, offset int64) *HumanizedJSONError {
	if err == nil {
		return nil
	}

	if offset < 0 || offset > int64(len(data)) {
		return &HumanizedJSONError{
			Err: errors.Wrapf(err, "Invalid offset %d.", offset),
		}
	}

	lineSep := []byte{'\n'}

	line := bytes.Count(data[:offset], lineSep) + 1
	lastLineOffset := bytes.LastIndex(data[:offset], lineSep)
	character := int(offset) - (lastLineOffset + 1) + 1

	return &HumanizedJSONError{
		Line:      line,
		Character: character,
		Err:       errors.Wrapf(err, "parsing error at line %d, character %d", line, character),
	}
}
