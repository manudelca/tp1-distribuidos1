package protocol

import "fmt"

type InvalidEventTypeError struct {
	eventType uint8
}

func (e InvalidEventTypeError) Error() string {
	return fmt.Sprintf("Unrecognized event type: %d", e.eventType)
}

type InvalidFieldTypeError struct {
	fieldType uint8
}

func (e InvalidFieldTypeError) Error() string {
	return fmt.Sprintf("Unrecognized field type: %d", e.fieldType)
}

type InvalidMessageFormatError struct {
	errorMsg string
}

func (e InvalidMessageFormatError) Error() string {
	return e.errorMsg
}

type InvalidMetricIdFieldError struct {
	errorMsg string
}

func (e InvalidMetricIdFieldError) Error() string {
	return e.errorMsg
}

type InvalidFloatFieldError struct {
	errorMsg string
}

func (e InvalidFloatFieldError) Error() string {
	return e.errorMsg
}

type InvalidAggregationFieldError struct {
	errorMsg string
}

func (e InvalidAggregationFieldError) Error() string {
	return e.errorMsg
}

type InvalidDateFieldError struct {
	errorMsg string
}

func (e InvalidDateFieldError) Error() string {
	return e.errorMsg
}

type MissingFieldError struct {
	errorMsg string
}

func (e MissingFieldError) Error() string {
	return e.errorMsg
}
