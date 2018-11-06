package parser

import (
	"io"
	"strings"
)

type SheetType int

const (
	XLS SheetType = iota
	XLSX
)

type ApplicationError struct {
	message string
}

func (e *ApplicationError) Error() string {
	return e.message
}

type DataError struct {
	message string
}

func (e *DataError) Error() string {
	return e.message
}

type ParsingErrors struct {
	aplicationErrors []ApplicationError
	dataErrors       []DataError
}

func (pe *ParsingErrors) Error() string {
	var errorMessages []string

	for _, ae := range pe.aplicationErrors {
		errorMessages = append(errorMessages, ae.Error())
	}

	for _, de := range pe.dataErrors {
		errorMessages = append(errorMessages, de.Error())
	}

	return strings.Join(errorMessages, ",")
}

// Parses the XLS(X) passed as parameters and returns the CSV contents, the request errors and other errors.
func Parse (r io.Reader, sheetType SheetType) ([][]string, error) {
	return nil, &ParsingErrors{[]ApplicationError{{"Not Implemented yet"}}, []DataError{}}
}