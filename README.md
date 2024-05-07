### Yandex Go Course Project
Медленный калькулятор (Slow Calc), работающий на технологиях gRPC, SQLite и JWT

###### для контактов:
###### tg: @Sergey_Poltorak
###### email: sergey.poltorak.main@gmail.com



##### Функциональность приложения
- Приложение выполняет регистрацию пользователя и предоставляет ему доступ к сервису по тем данным, которые он ввёл
- Далее пользователь вводит выражение, указывая время исполнения каждой операции
- Именно: возведение в степень, умножение, деление, сложение и вычитание
- Свой ID — UserID, а также ID выражения 

##### Какие цели стоят:
- Написать фронт
- Реализовать контроль воркеров

Попробовал сделать на 50 баллов — реализовал через gRPC взаимодействие сервера — оркестратора — и агента, сохранение выражения и данных пользователя в базе данных (SQLite) и аутентификацию пользователя при помощи JWT — JSON Web Token

## HOW TO INSTALL

##### Copy git repo for using if you had done that, otherwise:
``` bash
git clone https://github.com/SPA2006/Distributed-Calculator-of-Arithmetic-Expressions.git
```
##### and run in Go compiler or cmd: go run cmd/app/main.go
##### if you call go run in cmd go to the directory where you want to put files of the project and using "cd" command go to folder app and print in cmd: go run main.go
```
go run cmd/app/main.go
```

## NOTICE
Check availibility of ":8081", ":8079" ports for proper working of the project, specificially gRPC and HTTP ports

## curl запросы

#### Register user
```
curl -X POST http://localhost:8081/register
```
#### Request body:
```
{
  "id": number,
  "email": "your_email",
  "name": "your_name",
  "password": "your_password"
}
```

### Login user
```
curl -X POST http://localhost:8081/login
```
#### Request body:
```
{
  "id": unique_int,
  "email": "your_email",
  "name": "your_name",
  "password": "your_password"
}
```

### Input an expression:
```
curl -X POST http://localhost:8081/expression
```
#### Request body:
```
{
  "id": unique_int,
  "expression": "infix_expression",
  "user_id": your_id,
  "time_power": int_number,
  "time_mult": int_number,
  "time_div": int_number,
  "time_add": int_number,
  "time_sub": int_number
}
```

### Get expression from database by expression ID
```
curl -X GET http://localhost:8081/get_exp?id
```

### Get all expressions from database in one userID
```
curl -X GET http://localhost:8081/get_user_exp?user_id
```

### Get all expressions from database
```
curl -X GET http://localhost:8081/get_all_exp
```

### Delete expression from database
```
curl -X DELETE http://localhost:8081/delete_exp?id
```
