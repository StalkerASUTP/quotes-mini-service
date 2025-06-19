# Мини-сервис “Цитатник”
Реализовать REST API-сервис на Go для хранения и управления цитатами.
Необходимо реализовать следующее:
### 1. Выставить rest методы
1. Добавление новой цитаты (POST /quotes)
```json
{
    "author":"Confucius",
    "quote":"Life is simple, but we insist on making it complicated."
}
```
2. Получение всех цитат (GET /quotes)
3. Получение случайной цитаты (GET /quotes/random)
4. Фильтрация по автору (GET /quotes?author=Confucius)
5. Удаление цитаты по ID (DELETE /quotes/{id})

### 2. Проверочные команды (сurl)
```
curl -X POST http://localhost:8080/quotes \ 
-H "Content-Type: application/json" \ 
-d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```
```
curl http://localhost:8080/quotes
```
```
curl http://localhost:8080/quotes/random
```
```
curl http://localhost:8080/quotes?author=Confucius
```
```
curl -X DELETE http://localhost:8080/quotes/1
```
### Для запуска программы использовать команду
```
go run cmd/main.go
```    
### Для запуска тестов
```
go test -v ./...
```    
### В самом тестовом было ограничение на использование на сторнних библиотек, по этой причене в некоторых моментах использованы "костыли". Две библиотеки, которые пришлось использовать это: [GoDotEnv](https://github.com/joho/godotenv), [go-sqlite3](https://github.com/mattn/go-sqlite3).
### `GoDotEnv` использовался для загрузки локального `.env` файла, в котором хранится конфигурация, а `go-sqlite3` для использования драйвера для работы с БД.
### P.S. Для упрощения работы с валидации `json` стоило бы использовать [validator](https://github.com/go-playground/validator), для генерации `mock` [mockery]( https://github.com/vektra/mockery) и использовать какой-нибудь `router` со встроенными `middleware` и другими плюшками


 