# Сервис сокращения URL (СЕРВЕР)

### Задание по треку «Сервис сокращения URL»
Чтобы написать сервис, который будет сжимать длинные URL до нескольких символов, для начала вам нужно разработать сервер.
Сервер должен быть доступен по адресу `http://localhost:8080` и предоставлять два эндпоинта:

* Эндпоинт с методом `POST` и путём `/`. Сервер принимает в теле запроса строку URL как `text/plain` и возвращает ответ с кодом `201` и сокращённым URL как `text/plain`.

Пример запроса к серверу:
```
POST / HTTP/1.1
Host: localhost:8080
Content-Type: text/plain

https://practicum.yandex.ru/
```

Пример ответа от сервера:
```
HTTP/1.1 201 Created
Content-Type: text/plain
Content-Length: 30

http://localhost:8080/EwHXdJfB
```

* Эндпоинт с методом `GET` и путём `/{id}`, где `id` — идентификатор сокращённого URL (например, `/EwHXdJfB`). В случае успешной обработки запроса сервер возвращает ответ с кодом `307` и оригинальным URL в HTTP-заголовке `Location`.

Пример запроса к серверу:
```
GET /EwHXdJfB HTTP/1.1
Host: localhost:8080
Content-Type: text/plain
```

Пример ответа от сервера:
```
HTTP/1.1 307 Temporary Redirect
Location: https://practicum.yandex.ru/
```

* На любой некорректный запрос сервер должен возвращать ответ с кодом `400`.

### Запуск сервера:
Запуск из директории `/cmd/shortener`
Сервер может быть запущен с параметрами
```
-a value
        Server run address host:port (default localhost:8080)
-b value
        Base address for requests protokol://host:port
```

### Примеры запросов для тестирования:

curl -H 'Content-Type: text/plain' -d 'https://lenta.ru/' -X POST http://localhost:8080

curl -H 'Content-Type: text/plain' -d 'https://ya.ru/' -X POST http://localhost:8080

curl -H 'Content-Type: text/plain' -d 'https://market.yandex.ru/product--3526-2/749543053?sku=101098668745&uniqueId=847854&do-waremd5=SOnVO6hF_HDgQHQhCDQsyw&utm_term=70505730%7C749543053&clid=1601&utm_source=yandex&utm_medium=search&utm_campaign=ymp_offer_dp_komputer_model_mrkscr_bko_dyb_search_rus&utm_content=cid%3A113941570%7Cgid%3A5565834525%7Caid%3A1873011291195968827%7Cph%3A205565834525%7Cpt%3Apremium%7Cpn%3A5%7Csrc%3Anone%7Cst%3Asearch%7Crid%3A205565834525%7Ccgcid%3A20728017&yclid=11196615616457080831' -X POST http://localhost:8080

# Клиент для тестирования сервиса сокращенных URL

Клиент позволяет делать запросы к сервису и получать ответы.

### Запуск клиента
Запуск из директории `cmd/client`
Клиент может быть запущен с параметрами:
```
-a value
        Address to send requests protocol:://host:port (default http://localhost:8080/)
```
