package custom_errors

import "errors"

var (

	//For calc
	ErrExpectedNumber            = errors.New("expected number")
	ErrInvalidNumber             = errors.New("invalid number")
	ErrDivisionByZero            = errors.New("division by zero")
	ErrMissingClosingParenthesis = errors.New("missing closing parenthesis")
	ErrUnexpectedEndOfExpression = errors.New("unexpected end of expression")

	//For Handler
	ErrInvalidExpression = errors.New("invalid expression")

	//For Json
	ErrLoadConfiguration = errors.New("failed to load configuration")
)
