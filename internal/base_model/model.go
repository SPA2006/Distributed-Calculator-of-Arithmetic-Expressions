package model

type OperationClass string

const (
	Addition       OperationClass = "addition"
	Subtraction    OperationClass = "subtraction"
	Multiplication OperationClass = "multiplication"
	Division       OperationClass = "division"
)

func SetDuration(Duration map[OperationClass]int) {
	Duration["addition"] = 50
	Duration["subtraction"] = 50
	Duration["multiplication"] = 50
	Duration["division"] = 50
}

func ResetDuration(durrr []int, Duration map[OperationClass]int) {
	Duration["addition"] = durrr[0]
	Duration["subtraction"] = durrr[1]
	Duration["multiplication"] = durrr[2]
	Duration["division"] = durrr[3]
}

// Стоит первой, поскольку одна из главных задач данной программы --
// определение операций, которые ввёл пользователь (тестировщик)
func DefineOperation(s rune) (OperationClass, bool) {
	switch s {
	case '+':
		return Addition, true
	case '-':
		return Subtraction, true
	case '*':
		return Multiplication, true
	case '/':
		return Division, true

	default:
		return "", false
	}
}

// Проверяет, является ли операция, введённая пользователем, валидной
func IsAllowedOperation(operationClass OperationClass) bool {
	return operationClass == Addition || operationClass == Subtraction ||
		operationClass == Multiplication || operationClass == Division
}

// Функция возвращает приоритет операции. Чем выше приоритет операции при математических вычислениях,
// тем больше возвращаемое значение
func Priority(op rune) int {
	res, _ := DefineOperation(op)
	switch res {
	case Multiplication, Division:
		return 2
	case Addition, Subtraction:
		return 1
	default:
		return 0
	}
}
