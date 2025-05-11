package calc

import (
	"errors"
	"strings"
	"testing"

	"pkg/internal/custom_errors"
)

func TestCalc(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name            string
		expr            string
		want            float64
		wantErr         error
		wantErrContains string
	}

	tests := []testCase{
		{"simple addition", "1+2", 3, nil, ""},
		{"with spaces", " 1 + 2 ", 3, nil, ""},
		{"subtraction", "5-3", 2, nil, ""},
		{"multiplication", "4*2", 8, nil, ""},
		{"division", "8/4", 2, nil, ""},
		{"precedence", "2+3*4", 14, nil, ""},
		{"parentheses", "(2+3)*4", 20, nil, ""},
		{"nested parentheses", "((2+3))*2", 10, nil, ""},
		{"unary minus", "-3+5", 2, nil, ""},
		{"double unary", "--3", 3, nil, ""},
		{"floats", "3.5+2.1", 5.6, nil, ""},

		{"empty expression", "", 0, custom_errors.ErrUnexpectedEndOfExpression, ""},
		{"expected number", "a+1", 0, custom_errors.ErrExpectedNumber, ""},
		{"invalid number", ".a", 0, custom_errors.ErrInvalidNumber, ""},
		{"missing closing paren", "(1+2", 0, custom_errors.ErrMissingClosingParenthesis, ""},
		{"division by zero", "1/0", 0, custom_errors.ErrDivisionByZero, ""},
		{"unexpected characters", "1+2x", 0, nil, "unexpected characters"},
		{"unterminated factor", "-", 0, custom_errors.ErrUnexpectedEndOfExpression, ""},
		{"complex expression", "-(1+2)*-3", 9, nil, ""},
		{"operator at end", "(1+2)*", 0, custom_errors.ErrUnexpectedEndOfExpression, ""},
		{"deep nesting", "(((4)))", 4, nil, ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := Calc(tc.expr)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("%s: error = %v; want %v", tc.name, err, tc.wantErr)
				}
				return
			}
			if tc.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Errorf("%s: error = %v; want to contain %q", tc.name, err, tc.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tc.name, err)
				return
			}
			if got != tc.want {
				t.Errorf("%s: got = %v; want %v", tc.name, got, tc.want)
			}
		})
	}
}
