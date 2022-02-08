# Тестовое задание от Hezzl.com на позицию Go (GoLang) Developer
## Условие задания:
1. Описать proto файл с сервисом из 3 методов: добавить пользователя, удалить пользователя, список пользователей  
2. Реализовать gRPC сервис на основе proto файла на Go  
3. Для хранения данных использовать PostgreSQL  
4. На запрос получения списка пользователей данные будут кешироваться в redis на минуту и браться из редиса  
5. При добавлении пользователя делать лог в clickHouse   
6. Добавление логов в clickHouse делать через очередь Kafka

## Архитектура проекта:
### Схема данных:
* Основная модель - структура `User`:
```go
    type User struct {
	Id int `db:"id"`

	Username     string     `db:"username"`
	Password     []byte     `db:"password_hash"`
	Gender       string     `db:"gender"`
	Age          int        `db:"age"`
	Description  string     `db:"description"`
	
	Hash         []byte     `db:"user_hash"`
	Date         *time.Time `db:"start_date"`
}
```
* Таблица в PostgreSQL - users со схемой:
```go
    CREATE TABLE users
(
    id            serial       not null unique,
    username      varchar(255) not null unique,
    password_hash bytea,
    gender        varchar(10)  not null,
    age           int          not null,
    start_date    date         default now(),
    description   varchar(255) not null,
    user_hash     bytea
);  
``` 
* Схема Clickhouse-таблиц с настройкой представления:
```go
CREATE TABLE clicklogs (
    postgres_id Int32 Codec(DoubleDelta, LZ4),
    time DateTime Codec(DoubleDelta, LZ4),
    start_date ALIAS toDate(time),
    age Int32,
    username String,
    gender String,
    description String
) Engine = MergeTree
PARTITION BY toYYYYMM(time)
ORDER BY (postgres_id, time); 

CREATE TABLE clicklogs_queue (
    postgres_id Int32,
    time DateTime Codec(DoubleDelta, LZ4),
    age Int32,
    username String,
    gender String,
    description String
)
ENGINE = Kafka
SETTINGS kafka_broker_list = 'kafka-1:19092, kafka-2:29092, kafka-3:39092',
       kafka_topic_list = 'logs',
       kafka_group_name = 'readings_consumer_group1',
       kafka_format = 'JSONEachRow',
       kafka_max_block_size = 1048576; 

CREATE MATERIALIZED VIEW clicklogs_queue_mv TO clicklogs AS
SELECT postgres_id, time, username, age, gender, description
FROM clicklogs_queue;  

SHOW
TABLES();

SELECT *
FROM clicklogs;
``` 

### gRPC сервер, реализованный на основе proto файла и UserActionsService с 4 методами:
* __NewUser()__ - добавляет пользователя в таблицу PostgreSQL users и делает лог в Kafka-Clickhouse
* __DeleteUser()__ - удаялет пользователя из базы по `username`
* __GetUserList()__ - возвращает массив указателей на структуру `User` и кеширует данные в Redis на минуту
* __GetUserStringedList()__ - возвращает json-сериализованную строку данных, закешированных в Redis

Построен на принципе чистой архитектуры и обратном внедрении зависимостей.
По очереди иницииализирует слои:
1. `Repository` - [PosgresSQL]
2. `Cache` - [Redis]
3. `Clicklogs` - [Kafka-Clickhouse]
4. `Service` - [gRPC-methods]
Запускает `tcp-сервер` и регистрирует на нем `gRPC-службу`, слушая входящие соединения.

### Клиентское REST API, слушающее gRPC сервер и вызывающие его методы по собственным эндпоинтам:
* [__api/test/new__] - принимает json, соответствующий структуре `User` - возвращает json-`NewUserResponse`
* [__api/test/delete/:username__] - парсит username из `url` запроса - возвращает json-`DeleteUserResponse`
* [__api/test/list/list__] - возвращает json-`GetUserListResponse` в виде массива json-`User`
* [__api/test/list/string__] - возвращает json-`GetUserStringedListResponse` в виде json-строки из Redis\
*Структуры ...Response определены в proto файле - /grpc/proto каталоге*

 
### Конфигурация и настройки:
* За конфигурацию слоев внутри go отвечает `ConfigManager`, подтягивающий конфиги из `/configs/settings`
* `/configs/settings` содержит файлы, которые инициализируют Kafka, Redis и Postgres драйверы
* Также поддтягиваются переменные окружения из файла `.env`
* Инициализация конфигов происходит при каждом запуске client/server, данные с конфигами - персистентны

### Тестовое развертывание и запуск приложения:
* Развертывание сервисов вокруг `client/server` происходит внутри папки `/deploytment`
* Настройки сервисов и зависимости прописаны в файле `docker-compose.yml` и каталоге `/schema`
* Kafka разворачивается в виде кластера из трех `Kafka-brocker` и `ZooKeeper`
* Сервисы запускаются на `default-docker` сети, для чего мы указываем {MY_IP} перед выполнением docker-compose up\
*{MY_iP} можно посмотреть, выполнив: `ifconfig -a` в терминале*
* Перед первым запусокм grpc-сервера нужно поднять миграции, что делаеся через небольшое go-приложение, которое подключается к базе и выполняет последовательно `migrate-down` && `migrate-up`
* Команды для создания топика внутри Kafka-кластера, чтения из него в консольном клиенте, удаления кеша и мета-данных и `/shared` и остановки сервисов - прописаны в `make-файле`
* Подключение к Clickhouse-серверу делается через консольный клиент, через него же инициализируется схема таблиц и материализованное представление, которые лежат в `/schema` каталоге
