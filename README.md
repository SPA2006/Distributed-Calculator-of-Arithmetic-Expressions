# Distributed Calculator of Arithmetic Expressions
### Yandex Go Course Project
Медленный калькулятор (Slow Calc)

По всем вопросам:
###### tg: @Sergey_Poltorak

###### примеры curl запросов идут далее


##### // ————————————————————————

##### Это мой первый проект, который выполняет:
 - Запрос времени исполнения операций
 - Вычисление выражения, который вводит пользователь (User)
 - Переводит выражение в обратную польскую нотацию
 - Выдаёт ответ выисления

Попробовал сделать на минимальный бал

##### Для запуска микросервиса скачайте файлы с директории, либо импортируйте как модуль
##### https://github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions/
##### и запустите Distributed-Calculator-of-Arithmetic-Expressions/main/main.go
##### go run Distributed-Calculator-of-Arithmetic-Expressions/main/main.go



### curl запросы

#### Output duration of operations
``` bash
$ curl -X GET http://localhost:8080/duration
```

#### Input duration of operations
###### Addition
``` bash
$ curl -X POST http://localhost:8080/addition/{duration}
```

###### Substraction
``` bash
$ curl -X POST http://localhost:8080/substraction/{duration}
```

###### Multiplication
``` bash
$ curl -X POST http://localhost:8080/multiplication/{duration}
```

###### Division
``` bash
$ curl -X POST http://localhost:8080/division/{duration}
```



#### Input expression
``` bash
$ curl -X GET http://localhost:8080/{expression}
```

#### Input test expression (one of the test cases)
###### - standard
###### - zero
###### - letter
###### - symbol
###### - number
###### - bracket
``` bash
$ curl -X GET http://localhost:8080/test/{test_case}
```
