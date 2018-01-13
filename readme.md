# Заготовка для сервиса автодополнения

Для управления зависимостями использован [Glide](https://github.com/Masterminds/glide)

В качестве вспомогательного фреймворка [Go-Kit](https://github.com/go-kit/kit)

Для логирования используется [Logrus](https://github.com/sirupsen/logrus)

## Сборка

```sh
glide install
go build
```

## Использование

Запуск через `./city-autocomplete <адрес сервера для получения списка городов>`

После запуска сервис начинает слушать на порту 8080.

Досутупные настройки логирования можно посмотреть запустив сервис без параметров или с ключем `--help`.

Запрос

```sh
curl --request GET --url 'http://localhost:8080/query?query=a'
```
