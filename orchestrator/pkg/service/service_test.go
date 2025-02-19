package service

import (
	"slices"
	"testing"
)

func TestPostfixExpression(t *testing.T) {
	testCases := []struct {
		expression string
		expect     []string
		error      bool
	}{
		{
			"2+2*2",
			[]string{"2", "2", "2", "*", "+"},
			false,
		},
		{
			"(2+2)*2",
			[]string{"2", "2", "+", "2", "*"},
			false,
		},
		{
			"2*0",
			[]string{"2", "0", "*"},
			false,
		},
		{
			"2*(2+2",
			[]string{},
			true,
		},
	}

	for i, tc := range testCases {
		ex, err := PostfixExpression(tc.expression)
		if err != nil {
			if tc.error {
				continue
			} else {
				t.Errorf("test %d failed: unexpected error (%s)", i+1, err)
			}
		}

		if slices.Compare(ex, tc.expect) != 0 {
			t.Errorf("test %d failed", i+1)
		}
	}
}