// В этом файле вы найдёте:
// Тестовые случаи для проверки корректности работы микросервиса

package test

var Testing_cases = map[string]string{
	"standard": "1+1",
	"zero":     "01+1",
	"letter":   "1+a",
	"symbol":   "+-/*",
	"number":   "12345",
	"bracket":  "(1+2",
}
