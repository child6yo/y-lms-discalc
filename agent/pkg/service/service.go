package service

import (
	"errors"
	"strconv"
)

func EvaluatePostfix(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		if value, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, value)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, a/b)
			default:
				return 0, errors.New("unknown operator")
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}