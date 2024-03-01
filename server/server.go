package server

import (
	"errors"
	"strconv"

	model "calc/internal/base_model"
)

// перевод строки в число
func Get_string_number(value string, pos *int) (string, error) {
	var str_number = ""
	for ; *pos < len(value); *pos++ {
		if _, err := strconv.Atoi(string(value[*pos])); err == nil {
			str_number += string(value[*pos])
		} else {
			*pos--
			break
		}
	}

	if len(str_number) > 1 && str_number[0] == '0' {
		return "", errors.New("zero")
	}

	return str_number, nil
}

// Создать отдельный стек операций
// Использовать его для определения приоритета и соответственно вых. выражения
func Parse_expression(value string) (string, error) {
	var output string
	var stek []rune
	stek = make([]rune, 0)

	var counter = 0
	for counter < len(value) {
		if value[counter] == '(' {
			stek = append(stek, rune(value[counter]))
		} else if value[counter] == ')' {
			for len(stek) > 0 && stek[len(stek)-1] != '(' {
				output += string(stek[len(stek)-1]) + " "
				stek = stek[:len(stek)-1]
			}

			if len(stek) == 0 {
				return "", errors.New("bracket")
			}

			stek = stek[:len(stek)-1]

		} else if _, err := strconv.Atoi(string(value[counter])); err == nil {
			num, err := Get_string_number(value, &counter)

			if err != nil {
				if err.Error() == "zero" {
					return "", errors.New("zero")
				}
			}

			output += num + " "
		} else {
			for len(stek) > 0 && model.Priority(rune(value[counter])) <= model.Priority(stek[len(stek)-1]) {
				output += string(stek[len(stek)-1]) + " "
				stek = stek[:len(stek)-1]
			}

			stek = append(stek, rune(value[counter]))
		}

		counter++
	}
	// если что-то осталось, то переводим это в выходную строку
	for len(stek) > 0 {
		if stek[len(stek)-1] == '(' {
			return "", errors.New("bracket")
		}

		if stek[len(stek)-1] == '+' || stek[len(stek)-1] == '-' || stek[len(stek)-1] == '*' || stek[len(stek)-1] == '/' {
			return "", errors.New("symbol")
		}

		output += string(stek[len(stek)-1]) + " "
		stek = stek[:len(stek)-1]
	}

	return output, nil
}
