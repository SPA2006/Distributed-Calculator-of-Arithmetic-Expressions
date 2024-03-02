# Distributed Calculator of Arithmetic Expressions
### Yandex Go Course Project
Медленный калькулятор (Slow Calc)

###### примеры curl запросов идут далее

##### Это мой первый проект, который выполняет:
 - Запрос времени исполнения операций
 - Вычисление выражения, который вводит пользователь (User)
 - Переводит выражение в обратную польскую нотацию
 - Выдаёт ответ выисления

##### Что не получилось:
- Написать фронт, базу данных (Postgres, SQLlite)
- Реализовать решение через Docker или Docker Compose, поскольку команда
```
nvm install touch-cli -g
```
- Выдаёт ошибку
```
Error retrieving "https://nodejs.org/dist/latest-touch-cli/SHASUMS256.txt": HTTP Status 404
```

##### Какие цели стоят:
- Написать фронт
- Прописать базу данных и хранить выражение(-я) и промежуточные вычисления на ней
- Реализовать контроль воркеров

Попробовал сделать на минимальный балл

## HOW TO INSTALL

##### Copy git repo for using if you had done that, otherwise:
``` bash
git clone https://github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions.git
```
##### and run Distributed-Calculator-of-Arithmetic-Expressions/main/main.go
```
go run Distributed-Calculator-of-Arithmetic-Expressions/main/main.go
```



## curl запросы

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
