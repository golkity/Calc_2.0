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
	ErrOpenLogFile       = errors.New("failed to open log file")

	//For LoadConfig
	ErrOpenConfiguration  = errors.New("unable to open config file")
	ErrParceConfiguration = errors.New("unable to parse config file")

	//For Orchenstrator
	ErrInitOrchestrator = errors.New("unable to init orchestrator")

	//For Agent
	ErrInitAgent  = errors.New("unable to init agent")
	ErrStartAgent = errors.New("unable to start agent")

	//404

	ErrNotFound = errors.New("not found")
)
