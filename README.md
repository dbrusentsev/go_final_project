## Описание проекта

Веб-приложение для планирования задач. Позволяет создавать, редактировать, удалять задачи и отмечать их выполненными. Поддерживается функционал повторения задач по правилам: ежедневно, еженедельно, ежемесячно, ежегодно.

## Запуск локально

```bash
go build -o scheduler .
./scheduler
```

Адрес в браузере: http://localhost:7540

Переменные окружения:
`TODO_PORT` - порт сервера (по умолчанию 7540)
`TODO_DBFILE` - путь к файлу БД (по умолчанию scheduler.db)
`TODO_PASSWORD` - пароль для авторизации (если не задан, авторизация отключена)

Пример запуска с паролем:
```bash
TODO_PASSWORD=mypass ./scheduler
```

## Запуск тестов

Параметры в `tests/settings.go`:
```go
var Port = 7540
var DBFile = "../scheduler.db"
var Token = ``
```

Если включена аутентификация, получите токен и укажите его в Token:
```bash
curl -X POST http://localhost:7540/api/signin -d '{"password":"mypass"}'
```

Запуск тестов (сервер должен быть запущен):
```bash
go test ./tests
```

## Docker

Сборка:
```bash
docker build -t scheduler .
```

Запуск:
```bash
docker run -p 7540:7540 -v $(pwd)/data:/data scheduler
```

С паролем:
```bash
docker run -p 7540:7540 -v $(pwd)/data:/data -e TODO_PASSWORD=mypass scheduler
```

Адрес в браузере: http://localhost:7540
