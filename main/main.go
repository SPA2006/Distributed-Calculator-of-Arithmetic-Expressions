package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	agent "calc/agent"
	model "calc/internal/base_model"
	test "calc/internal/test"
)

func main() {
	Duration := make(map[model.OperationClass]int)
	model.SetDuration(Duration)

	var durrr []int
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Будете вводить время выполнения операций?")
		fmt.Println("--- по умолчанию время выполнения каждой операции -- 50 сек ---")
		fmt.Println()
		fmt.Println("Введите true, если да, false -- если нет")
		str_ans, _ := reader.ReadString('\n')
		str_ans = strings.TrimSpace(str_ans)
		ans, err := strconv.ParseBool(str_ans)

		if err != nil {
			for {
				fmt.Println("Вы ввели неправильно значение")
				fmt.Println("Введите true, если да, false -- если нет")
				str_ans, _ = reader.ReadString('\n')
				ans, err = strconv.ParseBool(str_ans)

				if err == nil {
					break
				}
			}
		}

		if ans {
			fmt.Println("Введите время выполнения операций:")
			durrr = readDuration("Сложение", durrr)
			durrr = readDuration("Вычитание", durrr)
			durrr = readDuration("Умножение", durrr)
			durrr = readDuration("Деление", durrr)
			model.ResetDuration(durrr, Duration)
		}

		fmt.Println("1. Проверить тестовые выражения")
		fmt.Println("2. Ввести собственное выражение")
		fmt.Println("3. Выход")
		str_ans, _ = reader.ReadString('\n')
		str_ans = strings.TrimSpace(str_ans)

		switch str_ans {
		case "1":
			fmt.Println("Введите название теста -- слово перед двоеточием:")
			fmt.Println("standard:   1+1")
			fmt.Println("suboptimal: 01+1")
			fmt.Println("letter:     1+a")
			fmt.Println("symbol:     +-/*")
			fmt.Println("number:     12345")
			fmt.Println("bracket:    (1+2")

			type_ans, _ := reader.ReadString('\n')
			type_ans = strings.TrimSpace(type_ans)

			result, err := agent.Reading_notation(test.Testing_cases[type_ans], Duration)
			err_check(result, err)

		case "2":
			fmt.Println("Введите собственное выражение:")
			type_ans, _ := reader.ReadString('\n')
			type_ans = strings.TrimSpace(type_ans)

			result, err := agent.Reading_notation(type_ans, Duration)
			err_check(result, err)

		case "3":
			break
		default:
			fmt.Println("Введите правильный ответ")
		}
	}
}

func readDuration(operation string, durrr []int) []int {
	var val int
	fmt.Printf("%s: ", operation)
	fmt.Fscan(os.Stdin, &val)
	durrr = append(durrr, val)
	return durrr
}

// проверяет наличие ошибки
func err_check(result int, err error) {
	if err != nil {
		switch err.Error() {
		case "zero":
			fmt.Println("В выражении есть незначащий ноль")
		case "letter":
			fmt.Println("В выражении есть буква(-ы)")
		case "symbol":
			fmt.Println("В выражении есть лишний символ")
		case "number":
			fmt.Println("В выражении находятся только числа")
		case "bracket":
			fmt.Println("В выражении есть лишняя скобка")
		default:
			fmt.Println("Неизвестная ошибка:", err)
		}
	} else {
		fmt.Println("Результат выражения")
	}
	fmt.Println("answer: ", result)
}
