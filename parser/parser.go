package parser

import (
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


