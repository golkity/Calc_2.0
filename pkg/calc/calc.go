package calc

import (
	"fmt"
	"strconv"
	"unicode"

	"pkg/internal/custom_errors"
)

func rmvspc(expression string) string {
	res := ""
	for _, char := range expression {
		if char != ' ' {
			res += string(char)
		}
	}
	return res
}

func parsnum(expression string, i *int) (float64, error) {
	var newstr string
	for *i < len(expression) && (unicode.IsDigit(rune(expression[*i])) || expression[*i] == '.') {
		newstr += string(expression[*i])
		*i++
	}

	if newstr == "" {
		return 0, custom_errors.ErrExpectedNumber
	}
	num, err := strconv.ParseFloat(newstr, 64)
	if err != nil {
		return 0, custom_errors.ErrInvalidNumber
	}
	return num, nil
}

func parsexp(expression string, i *int) (float64, error) {
	res, err := parsetrm(expression, i)
	if err != nil {
		return 0, err
	}
	for *i < len(expression) {
		char := expression[*i]
		if char == '+' || char == '-' {
			*i++
			nextT, err := parsetrm(expression, i)
			if err != nil {
				return 0, err
			}
			if char == '+' {
				res += nextT
			} else {
				res -= nextT
			}
		} else {
			break
		}
	}
	return res, nil
}

func parsetrm(expression string, i *int) (float64, error) {
	res, err := parsefct(expression, i)
	if err != nil {
		return 0, err
	}
	for *i < len(expression) {
		char := expression[*i]
		if char == '*' || char == '/' {
			*i++
			nextF, err := parsefct(expression, i)
			if err != nil {
				return 0, err
			}
			if char == '*' {
				res *= nextF
			} else {
				if nextF == 0 {
					return 0, custom_errors.ErrDivisionByZero
				}
				res /= nextF
			}
		} else {
			break
		}
	}
	return res, nil
}

func parsefct(expression string, i *int) (float64, error) {
	if *i >= len(expression) {
		return 0, custom_errors.ErrUnexpectedEndOfExpression
	}
	char := expression[*i]
	if char == '-' {
		*i++
		factor, err := parsefct(expression, i)
		if err != nil {
			return 0, err
		}
		return -factor, nil
	}

	if char == '(' {
		*i++
		res, err := parsexp(expression, i)
		if err != nil {
			return 0, err
		}
		if *i >= len(expression) || expression[*i] != ')' {
			return 0, custom_errors.ErrMissingClosingParenthesis
		}
		*i++
		return res, nil
	}

	return parsnum(expression, i)
}

func Calc(expression string) (float64, error) {
	expression = rmvspc(expression)
	i := new(int)
	result, err := parsexp(expression, i)
	if err != nil {
		return 0, err
	}
	if *i < len(expression) {
		return 0, fmt.Errorf("unexpected characters at position %d: %s", *i, expression[*i:])
	}
	return result, nil
}
