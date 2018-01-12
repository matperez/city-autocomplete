# Заготовка для сервиса автодополнения

Для управления зависимостями использован [Glide](https://github.com/Masterminds/glide)

В качестве вспомогательного фреймворка [Go-Kit](https://github.com/go-kit/kit)

## Сборка

```sh
glide install
go build
```

## Использование

Запуск через `./city-autocomplete`

Никаких параметров нет. Логирование производится в stdOut в виде json строк. После запуска сервис начинает слушать на порту 8080.

Запрос

```sh
curl --request GET --url 'http://localhost:8080/query?query=a'
```
