package main

import (
	"fmt"
	"net/http"
	"strconv"

	agent "calc/agent"
	model "calc/internal/base_model"
	test "calc/internal/test"
)

func main() {
	Duration := make(map[model.OperationClass]int)
	model.SetDuration(Duration)

	fmt.Println("Building REST APIs in GO")
	//  создаём новый мультиплексор
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "--- SLOWED CALCULATOR ---")
	})
	// блок кода для обработки обращения на http-сервер через
	// curl -X GET или POST http://localhost:8080/
	mux.HandleFunc("GET /duration", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "addition:       %d\n", Duration["addition"])
		fmt.Fprintf(w, "subtraction:    %d\n", Duration["subtraction"])
		fmt.Fprintf(w, "multiplication: %d\n", Duration["multiplication"])
		fmt.Fprintf(w, "division:       %d\n", Duration["division"])
	})

	mux.HandleFunc("POST /duration/addition/{duration}", func(w http.ResponseWriter, r *http.Request) {
		val, _ := strconv.Atoi(r.PathValue("duration"))
		readDuration(model.Addition, val, Duration)
	})

	mux.HandleFunc("POST /duration/substraction/{duration}", func(w http.ResponseWriter, r *http.Request) {
		val, _ := strconv.Atoi(r.PathValue("duration"))
		readDuration(model.Subtraction, val, Duration)
	})

	mux.HandleFunc("POST /duration/multiplication/{duration}", func(w http.ResponseWriter, r *http.Request) {
		val, _ := strconv.Atoi(r.PathValue("duration"))
		readDuration(model.Multiplication, val, Duration)
	})

	mux.HandleFunc("POST /duration/division/{duration}", func(w http.ResponseWriter, r *http.Request) {
		val, _ := strconv.Atoi(r.PathValue("duration"))
		readDuration(model.Division, val, Duration)
	})

	// реализация запроса на сервер
	// получаем от сервера результат выражения, введённого пользоавтелем

	mux.HandleFunc("GET /{expression}", func(w http.ResponseWriter, r *http.Request) {
		expression := r.PathValue("expression")
		result, err := agent.Reading_notation(test.Testing_cases[expression], Duration)
		err_check(w, result, err)
	})

	// получаем от сервера результат тестового выражения id (название), которого ввёл пользователь

	mux.HandleFunc("GET /test/{test_case}", func(w http.ResponseWriter, r *http.Request) {
		case_name := r.PathValue("test_case")

		_, ok := test.Testing_cases[case_name]

		if ok {
			result, err := agent.Reading_notation(test.Testing_cases[case_name], Duration)
			err_check(w, result, err)
		}
	})

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err.Error())
	}
}

func readDuration(operation model.OperationClass, val int, Duration map[model.OperationClass]int) {
	Duration[operation] = val
}

// проверяет наличие ошибки
func err_check(w http.ResponseWriter, result int, err error) {
	if err != nil {
		switch err.Error() {
		case "zero":
			fmt.Fprintf(w, "В выражении есть незначащий ноль")
		case "letter":
			fmt.Fprintf(w, "В выражении есть буква(-ы)")
		case "symbol":
			fmt.Fprintf(w, "В выражении есть лишний символ")
		case "number":
			fmt.Fprintf(w, "В выражении находятся только числа")
		case "bracket":
			fmt.Fprintf(w, "В выражении есть лишняя скобка")
		default:
			fmt.Fprintf(w, "Неизвестная ошибка: %s", err.Error())
		}
	} else {
		fmt.Fprintf(w, "")
	}
	fmt.Fprintf(w, "Результат выражения: %d", result)
}
