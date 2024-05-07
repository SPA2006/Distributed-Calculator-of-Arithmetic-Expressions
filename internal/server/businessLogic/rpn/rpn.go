package rpn

import (
	"errors"
	"slices"
	"strings"
)

var (
	priority = map[string]int{
		"^": 4,
		"*": 3,
		"/": 3,
		"+": 2,
		"-": 2,
		"(": 1,
	}
)

func Get_rpn(value string) (string, error) {
	stack := []string{}
	output := ""

	value = strings.ReplaceAll(value, " ", "")
	sep_values, err := split_string(value)

	if err != nil {
		return "", err
	}

	for _, token := range sep_values {
		switch token {
		// if the token (operand) is the operation of
		// multiplication, division, addition, substraction
		case "^", "*", "/", "+", "-":
			for len(stack) > 0 && priority[stack[len(stack)-1]] >= priority[token] {
				ln := len(stack)
				output += stack[ln-1] + " "
				stack = slices.Delete(stack, ln-1, ln)
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)

		case ")":
			// move all the operands to output until there's the opening brakcet
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				ln := len(stack)
				output += stack[ln-1] + " "
				stack = slices.Delete(stack, ln-1, ln)
			}
			ln := len(stack)
			stack = slices.Delete(stack, ln-1, ln)
		default:
			if token[0] >= '0' && token[0] <= '9' {
				output += token + " "
			} else {
				return "", errors.New("extra symbols in the expression")
			}
		}
	}

	// if there's any symbols in the stack stayed after pasrsing the values
	// from the original expression we push them to output
	// reversed polish notation string which we use for our further
	// computing on agent

	for len(stack) > 0 {
		ln := len(stack)

		if stack[ln-1] == "(" {
			return "", errors.New("too much brackets")
		}

		output += stack[ln-1] + " "
		stack = slices.Delete(stack, ln-1, ln)
	}

	ln := len(output)
	output = output[:ln-1]

	return output, nil
}

func split_string(value string) ([]string, error) {
	res := []string{}
	count_bracket := 0
	num := ""

	// split the string into the separate values
	for _, token := range value {
		if token >= '0' && token <= '9' {
			num += string(token)
		} else {
			if num != "" {
				res = append(res, num)
			}
			num = ""

			// if there's too excessive brackets, we'll raise the error
			if token == '(' {
				count_bracket++
			}
			if token == ')' {
				count_bracket--
			}
			res = append(res, string(token))
		}
	}
	// if something wasn't passed, we pass it to the output slice
	// of separated values
	if num != "" {
		res = append(res, num)
	}

	if count_bracket != 0 {
		return []string{}, errors.New("too much brackets")
	}

	return res, nil
}
