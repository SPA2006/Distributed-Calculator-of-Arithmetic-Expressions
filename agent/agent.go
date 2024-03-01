package agent

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	model "calc/internal/base_model"
	server "calc/server"
)

// Считывает выражение, переведенное на ОПН
// Проверяет на наличие ошибок (zero, bracket, symbol -- calc/test/test.go)
func Reading_notation(value string, Duration map[model.OperationClass]int) (int, error) {
	res, err := server.Parse_expression(value)
	fmt.Println("!!!   ", res)

	if err != nil {
		if err.Error() == "zero" || err.Error() == "bracket" || err.Error() == "symbol" {
			return 0, err
		}
	}

	res_separated := strings.Split(res, " ")
	for i := 0; i < len(res_separated); i++ {
		if _, err := strconv.Atoi(res_separated[i]); err != nil {
			break
		}

		if i == len(res_separated)-1 {
			return 0, errors.New("number")
		}
	}

	var mem1 int
	var mem2 int

	var stek []string
	var counter = 0

	mem3 := make(chan int, 2)
	// вычисление каждой операции в горутинах (каждая операция определяется в model.go)
	for counter < len(res_separated) {
		if res_separated[counter] == "+" || res_separated[counter] == "-" || res_separated[counter] == "*" || res_separated[counter] == "/" {
			mem1, _ = strconv.Atoi(stek[len(stek)-1])
			mem2, _ = strconv.Atoi(stek[len(stek)-2])
			res, _ := model.DefineOperation(rune(res_separated[counter][0]))

			if res == model.Addition {
				go func() {
					time.Sleep(time.Second * time.Duration(Duration["addition"]))
					mem3 <- mem2 + mem1
				}()
			} else if res == model.Subtraction {
				go func() {
					time.Sleep(time.Second * time.Duration(Duration["substraction"]))
					mem3 <- mem2 - mem1
				}()
			} else if res == model.Multiplication {
				go func() {
					time.Sleep(time.Second * time.Duration(Duration["multiplication"]))
					mem3 <- mem2 * mem1
				}()
			} else {
				go func() {
					time.Sleep(time.Second * time.Duration(Duration["division"]))
					mem3 <- mem2 / mem1
				}()
			}
			stek = stek[:len(stek)-2]
			stek = append(stek, strconv.Itoa(<-mem3))
		} else if _, err := strconv.Atoi(res_separated[counter]); err == nil {
			stek = append(stek, string(res[counter]))
		} else {
			return 0, errors.New("letter")
		}
	}

	ans, _ := strconv.Atoi(stek[0])
	return ans, nil
}
