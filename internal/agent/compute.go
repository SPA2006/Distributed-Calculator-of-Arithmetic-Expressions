package agent

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"strings"
)

// why do we need this function?
// it provides the main business logic of our app
// that's we all need from that
// calculate :)
// every single investor expects it!
func ComputeRPN(postfix string) (int, error) {
	sepPostfix := strings.Split(postfix, " ")
	var stack []int

	for _, value := range sepPostfix {
		if num, err := strconv.Atoi(value); err == nil {
			stack = append(stack, num)
		} else if value == "^" || value == "*" || value == "/" ||
			value == "+" || value == "-" {
			if len(stack) < 2 {
				return 0, errors.New("not enough operands for calculation")
			}

			ln := len(stack)
			val1 := stack[ln-2]
			val2 := stack[ln-1]
			stack = slices.Delete(stack, ln-2, ln)

			switch value {
			case "+":
				stack = append(stack, val1+val2)
			case "-":
				stack = append(stack, val1-val2)
			case "*":
				stack = append(stack, val1*val2)
			case "/":
				if val2 == 0 {
					return 0, errors.New("dividing by zero")
				}

				stack = append(stack, val1/val2)
			case "^":
				stack = append(stack, int(math.Pow(float64(val1), float64(val2))))
			}
		} else {
			return 0, errors.New("extra symbols")
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("incorrect number of operations")
	}

	return stack[0], nil
}

func InitAgent(user string) {

}
