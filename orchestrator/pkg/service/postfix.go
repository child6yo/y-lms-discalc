package service

import (
	"errors"
	"strings"
	"unicode"
)

func (s *Service) PostfixExpression(expression string) ([]string, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}

	return infixToPostfix(tokens)
}

func tokenize(expr string) ([]string, error) {
	var tokens []string
	var number strings.Builder

	for _, r := range expr {
		if unicode.IsSpace(r) {
			continue
		}
		if unicode.IsDigit(r) || r == '.' {
			number.WriteRune(r)
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			if strings.ContainsRune("+-*/()", r) {
				tokens = append(tokens, string(r))
			} else {
				return nil, errors.New("invalid character in expression")
			}
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}
	return tokens, nil
}

func infixToPostfix(tokens []string) ([]string, error) {
	var postfix []string
	var stack []string

	precedence := map[string]int{
		"+": 1, "-": 1,
		"*": 2, "/": 2,
	}

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1] != "(" {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		default:
			postfix = append(postfix, token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return postfix, nil
}
