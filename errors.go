package tei

import "errors"

var (
	// ErrEmptyInputs is returned when the inputs are empty
	ErrEmptyInputs = errors.New("inputs cannot be empty")

	ErrValidation = errors.New("validation error")
	ErrTokenizer  = errors.New("tokenizer error")
	ErrBackend    = errors.New("backend error")
	ErrOverloaded = errors.New("server overloaded")
)
