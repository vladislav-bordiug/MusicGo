# Реализация онлайн библиотеки песен

## Переменные окружения
+ ```DATABASE_URL``` - URL к базе данных PostgreSQL
+ ```API_URL``` - URL к API, куда будет отправлять GET запрос по маршруту /info
+ ```SERVER_IP``` - IP сервера в ListenAndServe
+ ```PORT``` - порт сервера в ListenAndServe

## Документация
Документация в [docs.go](docs/docs.go), [swagger.json](docs/swagger.json) и [swagger.yaml](docs/swagger.yaml).

## Методы

+ getdata - получение данных библиотеки с фильтрацией по всем полям и пагинацией
+ getsongtext - получение текста песни с пагинацией по куплетам
+ deletesong - удаление песни
+ editsong - изменение данных песни
+ addsong - добавление новой песни

## База данных
База данных PostgreSQL была поднята с docker-compose [docker-compose.yml](docker-compose.yml).