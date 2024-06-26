# Тестовое задание для LaModa 
## Задание
### API для работы с товарами на складе
Необходимо спроектировать и реализовать API методы для работы с товарами на
одном складе. Учесть, что вызов API может быть одновременно из разных систем и они могут работать с одинаковыми товарами. Методы API можно расширять доп. параметрами на своё усмотрение.
#### Спроектировать и реализовать БД для хранения следующих сущностей
-  Склад

• название

• признак доступности
- Товар

• название

• размер

• уникальный код

• количество

#### Реализовать методы API:
- резервирование товара на складе для доставки

• на вход принимает:

• массив уникальных кодов товара

- освобождение резерва товаров

• на вход принимает

• массив уникальных кодов товара

- получение кол-ва оставшихся товаров на складе

• на вход принимает:

• идентификатор склада

#### Будет плюсом
✅ Реализация логики работы с товарами, которые одновременно могут
находиться на нескольких складах

## Реализация
#### Конфигурация сервиса

Для конфигурирования сервиса ипользуется файл, лежащий в deploy/.env.

- SERVER_HOST // адрес сервера 
- SERVER_PORT // порт сервера
- DB_HOST // адрес БД с портом
- DB_USER // имя пользователя для подключения к БД
- DB_PASS // пароль пользователя для подключения к БД
- DB_NAME // имя БД
- DB_DISABLE_TLS // флаг использования TLS для подключения к БД
- LOG_LEVEL // уровень логгирования
- GOODS // количество генерируемых товаров
- STOCKS // количество генерируемых стоков на складах
- WAREHOUSES // количество генерируемых складов
- SWAGGER_PORT // порт, на котором запущен SwaggerUI

#### Запуск сервиса

Для запуска сервиса необходимо ввести команду:

```cmd
make up
```

#### Остановка сервиса

Для остановки сервиса необходимо ввести команду:

```cmd
make down
```
#### Просмотр логов сервиса

Для вывода логов сервиса необходимо ввести команду:

```cmd
make logs
```

#### Запуск тестов
Тестами покрыт слой стора.
Для запуска юнит-тестов необходимо ввести команду:
```cmd
make tests 
```
Для запуска анализа покрытия тестами необходимо ввести команду:
```cmd
make coverage 
```
#### Дополнительные команды
Для форматирования кода можно использовать команду:
```cmd
make fmt 
```
Для генерации swagger-документации можно использовать команду:
```cmd
make swagger-generate 
```
Для генерации функции, конфигурирующей всех зависимостей сервера, можно использовать команду:
```cmd
make wire 
```
## Дополнительная информация

### Используемые технологии
Go: 
- urfave/cli, wire, echo
- sqlx, go-playground/validator
- testcontainers-go, slog, swaggo

DB: 
- PostgresSQL

Swagger: 
- OpenAPI v2

Deploy: 
- Docker, docker-compose

### Swagger
С дефолтным env, после запуска сервиса, появляется доступ к SwaggerUI по адресу:
http://localhost:9000/doc/index.html

### Postman
Postman-коллекция доступна по ссылке: [тык](https://github.com/ArinLin/LaModa_testTask/blob/main/laModa_Inventory%20Hub.postman_collection.json "тык")
