# Распределённый вычислитель арифметических выражений (ФИНАЛЬНЫЙ ПРОЕКТ ЯНДЕКС ЛИЦЕЯ) :(

![version](https://shields.microej.com/github/go-mod/go-version/golkity/Calc?style=for-the-badge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)
[![GitHub Repo stars](https://img.shields.io/github/stars/USERNAME/REPOSITORY?style=social)](https://github.com/golkity/Calc_2.0)

![intro](source/final_logoYL.png)

## Структура проекта

<pre>
Calc_2.0/
├── infostructure/               
│   ├── docker/                    
│   │   └── create-postgres-databases.sh
│   ├── kafka-data/                
│   ├── kui-data/                
│   ├── postgres-data/ …         
│   ├── redis-data/               
│   ├── zk-data/                  
│   ├── zk-txn-logs/             
│   ├── .env                      
│   ├── .example.env            
│   ├── .gitignore                 
│   └── docker-compose.yml                   
│
├── internal/                        
│   ├── custom_errors/
│   │   └── custom_errors.go
│   ├── middleware/
│   │   └── middleware.go
│   ├── store/
│   │   └── store.go
│   └── task/
│       └── manager.go
│
├── pkg/                             
│   ├── calc/
│   │   ├── calc.go
│   │   ├── calc_test.go
│   │   └── calc_test.yaml
│   ├── logger/
│   │   └── logger.go
│   └── tokens/
│       └── manager.go
│
├── service/                          
│   ├── api-gateway/
│   │   ├── Dockerfile
│   │   ├── cmd/main.go
│   │   ├── go.mod | go.sum
│   │   ├── handler/
│   │   │   ├── calculate.go
│   │   │   ├── get_exp.go
│   │   │   └── healthz.go
│   │   ├── kafka/producer.go
│   │   └── middleware/auth.go
│   │
│   ├── auth/
│   │   ├── Dockerfile
│   │   ├── cmd/main.go
│   │   ├── go.mod | go.sum
│   │   ├── internal/
│   │   │   ├── config/config.go
│   │   │   ├── db/postgres.go
│   │   │   ├── http/
│   │   │   │   ├── handler/{handler.go,http_mid.go}
│   │   │   │   ├── routes/routes.go
│   │   │   │   └── server.go
│   │   │   ├── kafka/producer.go
│   │   │   ├── models/user.go
│   │   │   └── tokens/{jwt.go,redis_store.go}
│   │   └── migrations/001_init_schema.sql
│   │
│   ├── calc-orchenstrator/
│   │   ├── Dockerfile
│   │   ├── cmd/main.go
│   │   ├── go.mod | go.sum
│   │   ├── kafka/{consumer.go,producer.go}
│   │   ├── orchenstrator/orchenstrator.go
│   │   └── repository/{expression_pg.go,task_pg.go}
│   │
│   └── calc-worker/
│       ├── Dockerfile
│       ├── cmd/main.go
│       ├── go.mod | go.sum
│       ├── kafka/{consumer.go,producer.go}
│       └── worker/work.go
│
├── source/                          
│   ├── final_logoYL.png
│   ├── img.png
│   ├── logger.png
│   ├── logo_int.png
│   ├── logo_out.png
│   ├── outro.png
│   └── worker.jpg
│
├── LICENSE
└── Readme.md
</pre>

Ниже показана топология микросервисов, поднимаемых через `docker-compose`:

```mermaid
graph TD
    subgraph docker-compose
        PG[(Postgres)]
        ZK[(ZooKeeper)]
        KF[(Kafka)]
        GW[api-gateway]
        AU[auth]
        ORC[calc-orchestrator]
        WRK[calc-worker]
    end
```

## О приложение

>[!IMPORTANT]
> ### **Auth-Service**
> 	- POST /register – регистрация нового пользователя
>	- POST /login – выдача JWT (access + refresh)
>	- POST /refresh – обновление access-токена
>	- GET /healthz – проверка «здоровья» сервиса
>	- Хранит пользователей в PostgreSQL и рефреш-токены в Redis
> ### **API-Gateway**
>	- GET /healthz – проверяет доступность Kafka
>	- POST /api/v1/calculate – принимает выражение, валидирует его, сохраняет в БД, публикует задачу в Kafka и возвращает { "id": … }
>	- При синтаксической ошибке в выражении – 500 Internal Server Error
>	- GET /api/v1/expressions – возвращает список всех выражений пользователя
>	- GET /api/v1/expressions/{id} – возвращает конкретное выражение по ID
> ### **Calc-Orchestrator**
>	- Консьюмер Kafka-топика calc.tasks – получает { user_id, expression }, создаёт запись в Postgres (status=pending)
>	- При необходимости дробит длинные выражения на подзадачи, публикует их обратно в Kafka
>	- Консьюмер Kafka-топика calc.results – получает { expr_id, result }, помечает задачу как выполненную и сохраняет result
> ### **Calc-Worker (Agent)**
>	- Запускается с параллельными воркерами по COMPUTING_POWER
>	- Каждый воркер в цикле:
>	  	- GET /internal/task у Orchestrator
>		- Эмулирует «сложное» вычисление (time.Sleep(operation_time))
>		- Вычисляет результат через pkg/calc.Calc(...)
>		- POST /internal/task с итогом
>		- Повторяет
> ### **Calc (pkg/calc)**
>	- Модуль парсинга и вычисления арифметических выражений
>	- Экспортируемая функция:
>		```go
>		 func Calc(expr string) (float64, error)
>		```

```mermaid
graph LR
    subgraph Client Side
        U(User)
    end
    subgraph Edge
        G(API-Gateway)
    end
    subgraph Core
        OR(Calc-Orchestrator)
        W(Calc-Worker)
    end
    subgraph Infrastructure
        K(Kafka)
        P(Postgres)
        Z(ZooKeeper)
    end

    U -->|HTTP| G
    G -->|verify JWT| A(Auth)
    A -.-> P
    G -->|produce Task| K
    OR <-->|consume Task| K
    OR -->|enqueue Work| K
    W <-->|consume Work| K
    W -->|produce Result| K
    OR <-->|consume Result| K
    OR --> P
    G <-->|result| OR
```

## Шаг 1. Регистрация нового пользователя
Клиент отправляет `POST /register` к Auth-Service:

```bash
curl -i -X POST http://localhost:8080/api/v1/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@mail.com","password":"MEGASECRETPASSWORDYL"}'
```

```mermaid
flowchart LR
    U[Пользователь] -->|"POST /api/v1/register"| A[Auth-Service]
    A -->|"201 Created"| U
```

## Шаг 2. Логин

Пользователь аутентифицируется в Auth-Service, получая `access_token` и `refresh_token`, которые понадобятся для дальнейших запросов.

```bash
# Выполняем логин и сохраняем access_token
ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
        "email":    "demo@mail.com",
        "password": "Pa$$w0rd"
      }' \
  | jq -r .access_token)

echo "Access token:" $ACCESS_TOKEN
```

```mermaid
flowchart LR
    U[Пользователь] -->|"POST /api/v1/login \n { email, password }"| A[Auth-Service]
    A -->|"200 OK \n { access_token, refresh_token }"| U
```

## Шаг 3. Отправка выражения на вычисление

С помощью полученного access_token пользователь отправляет арифметическое выражение через API-Gateway. Если выражение валидно – возвращается идентификатор задачи.

```bash
curl -i -X POST http://localhost:8090/api/v1/calculate \
     -H "Authorization: Bearer $ACCESS_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "expression": "2+2*3"
         }'
```

```mermaid
flowchart LR
    U[Пользователь] -->|"POST /api/v1/calculate\n{ expression }"| G[API-Gateway]
    G -->|"200 OK\n{ id }"| U
```

Агент: отправляет запрос GET /internal/task.
Оркестратор: отдаёт задачу (200 OK) в JSON (или 404, если задач нет).

## Шаг 4. Оркестратор сохраняет и распределяет задачи

Что происходит: Calc-Orchestrator читает из Kafka, сохраняет выражение в БД, разбивает на подзадачи и публикует их.

```mermaid
flowchart LR
    G[API-Gateway] -->|"produce calc.tasks"| K[Kafka]
    K -->|"consume calc.tasks"| O[Calc-Orchestrator]
    O -->|"INSERT INTO expressions(status=pending)"| DB[(Postgres)]
    O -->|"produce subtasks (если нужно)"| K
```

## Шаг 5. Пользователь получает результат

Что происходит: Calc-Worker опрашивает оркестратор, вычисляет, отправляет результат.

```mermaid
flowchart LR
    W[Calc-Worker] -->|"GET /internal/task"| O[Orchestrator]
    O -->|"200 OK\nid,arg,operation_time"| W

    %% Разделяем вычисление и ожидание
    W -->|"compute Calc(arg)"| W
    W -->|"sleep(operation_time ms)"| W

    %% Отправка результата
    W -->|"POST /internal/task\nid,result"| O
    O -->|"update status & result"| DB[(Postgres)]
```


## Шаг 6. Получение списка выражений

Что происходит: клиент получает все свои выражения через API-Gateway.

```bash
curl -i http://localhost:8090/api/v1/expressions \
     -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
```

```mermaid
flowchart LR
    U[User] -->|"GET /api/v1/expressions"| G[API-Gateway]
    G -->|"200 OK\n{expressions:[…]}"| U
```

## Шаг 7. Получение конкретного выражения по ID

Что происходит: клиент запрашивает статус и результат одного выражения.

```bash
curl -i http://localhost:8090/api/v1/expressions/1 \
     -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
```

```mermaid
flowchart LR
    U[User] -->|"GET /api/v1/expressions/1"| G[API-Gateway]
    G -->|"200 OK\n{expression:{…}}"| U
```


>[!IMPORTANT]
> Что такое воркеры и как они работают?
>![worker](source/worker.jpg)
>
>В коде агента (`calc-worker`) реализован механизм параллельных «воркеров», который позволяет обрабатывать несколько задач одновременно:
>	1. При запуске функции `Start()` агент считывает из переменной окружения `COMPUTING_POWER` число воркеров (cp).  
>	2. Каждый воркер стартует в собственной горутине:
>   ```go
>      for i := 0; i < cp; i++ {
>           go workerLoop(i)
>       }
>   ```
>	3.	В цикле каждый воркер выполняет:
>	    -	getTask() — отправляет запрос GET /internal/task к Оркестратору, пытаясь получить задачу.
>	    -	Если задачи нет (404 Not Found), воркер делает time.Sleep(2 * time.Second) и снова пытается получить задачу.
>	    -	Если задача есть, воркер:
>	        1.	Вычисляет выражение с помощью функции calc.Calc(...).
>	        2.	Ждёт время, указанное в OperationTime (эмуляция «сложности» вычислений или иных затрат).
>	        3.	Отправляет результат обратно (метод sendResult(...)).
>	        4.	Спит 1 секунду и возвращается в начало цикла.
>	        4.	Все воркеры работают независимо и параллельно, позволяя агенту обрабатывать несколько задач одновременно.
>    
> Все воркеры работают независимо и параллельно, обеспечивая высокую пропускную способность и устойчивость к пиковым нагрузкам.

>[!NOTE]
> Пользователь вызывает функцию Calc("2+2*2").
> -	Calc() сначала вызывает rmvspc(), чтобы удалить пробелы (здесь строка остаётся такой же).
> -	Затем Calc() вызывает parsexp(), которая начинает разбирать выражение.
> -	Внутри parsexp() вызывается parsetrm(), которая, в свою очередь, вызывает parsefct() и parsnum(), чтобы извлечь число 2.
> -	Результат 2 возвращается обратно по цепочке до parsexp(), где обнаруживается оператор "+".
> -	Далее для правой части выражения вызывается parsexp("2*2"), которая обрабатывается через parsetrm(), parsefct() и parsnum() для получения результата 4.
> -	В итоге Calc() суммирует промежуточные результаты (2 + 4) и возвращает 6.

```mermaid
sequenceDiagram
    participant U as User
    participant G as API-Gateway
    participant A as Auth
    participant K as Kafka
    participant O as Orchestrator
    participant W as Worker
    participant DB as Postgres

    U->>G: POST /calculate (expr)
    G->>A: JWT introspect
    A-->>G: 200 OK
    G->>K: Produce Task(topic=tasks)
    O-->>K: Consume Task
    O->>K: Produce Work(topic=work)
    W-->>K: Consume Work
    W->>K: Produce Result(topic=results)
    O-->>K: Consume Result
    O->>DB: Save result
    O-->>G: result
    G-->>U: 200 OK + JSON
```

## Какие бывают запросы?? :trollface:

>[!TIP]
> 201 (Created) <- пользователь зарегестрирован в базе данных
> ```shell
> ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/registration \
>  -H "Content-Type: application/json" \
>  -d '{
>        "email":    "demo@mail.com",
>        "password": "Pa$$w0rd"
>      }' \
>  | jq -r .access_token)
>
>echo "Access token:" $ACCESS_TOKEN
> ```
> 
> Ответ:
> ![img](./source/3.png)
>
> 200 (OK) <- успешный вход
> ```shell
> ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
>  -H "Content-Type: application/json" \
>  -d '{
>        "email":    "demo@mail.com",
>        "password": "Pa$$w0rd"
>      }' \
>  | jq -r .access_token)
>
>echo "Access token:" $ACCESS_TOKEN
> ```
> 
> Ответ:
> ![img](./source/4.png)
> 
> 200 (ОК) <- успешно получено выражение
> ```shell
> curl -i -X POST http://localhost:8090/api/v1/calculate \
>    -H "Authorization: Bearer $ACCESS_TOKEN" \
>     -H "Content-Type: application/json" \
>     -d '{
>           "expression": "9-36"
>         }'
> ```
> 
> Ответ:
> ![img](./source/5.png)
> 
> 200 (OK)
> 
> ```shell
> curl -i http://localhost:8090/api/v1/expressions/{id} \
>     -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
> ```
> 
> Ответ:
>![img](./source/6.png)


>[!CAUTION]
> 409 (Conflict) <- пользователь с такой почтой уже есть в БД
> ```shell
> ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/registration \
>  -H "Content-Type: application/json" \
>  -d '{
>        "email":    "demo@mail.com",
>        "password": "Pa$$w0rd"
>      }' \
>  | jq -r .access_token)
>
>echo "Access token:" $ACCESS_TOKEN
> ```
>![img](./source/1_auth.png)
>
>
>401 (Unauthorized) ошибка в регистрации
> ```shell
> ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
>  -H "Content-Type: application/json" \
>  -d '{
>        "email":    "deo@mail.com",
>        "password": "Pa$$w0rd"
>      }' \
>  | jq -r .access_token)
>
>echo "Access token:" $ACCESS_TOKEN
> ```
> 
> Ответ:
> ![img](./source/2_auth.png)
> 
> 500 (ОК) <- invalid expression: expected number
> ```shell
> curl -i -X POST http://localhost:8090/api/v1/calculate \
>    -H "Authorization: Bearer $ACCESS_TOKEN" \
>     -H "Content-Type: application/json" \
>     -d '{
>           "expression": "2+y"
>         }'
> ```
> 
> Ответ:
> ![img](./source/1_calc.png)
>
>404 (Not Found) когда мы указываем id,которого нету
> ```shell
> curl -i http://localhost:8090/api/v1/expressions/{id} \
>     -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
> ```
> 
> Ответ:
>![img](./source/2_calc.png)



## ТЕСТЫ??? НОУУ ВЭЭЙ 

>[!IMPORTANT]
> Как их запускать и зачем их есть?<br>
> **Тесты для Handler.go**
> ```shell
> cd internal/http/handler
> go test -v
> ```

>[!IMPORTANT]
> **Тесты для Calc.go**
> ```shell
> cd pkg/calc
> go test -v
> ```
>
>**Тесты для agent.go**
> ```shell
> cd internal/agent
> go test -v
> ```
>
>**Тесты для orchenstrator.go**
> ```shell
> cd internal/orchenstrator
> go test -v
> ```
>**Тесты для config.go**
> ```shell
> cd config
> go test -v
> ```

## Kafka-UI

>[!IMPORTANT]
>
> После запуска контейнера Kafka-UI, перейдите по ссылке:
> ![kafka](./source/ui.png)
> Потом заполните конфигурацию, как на скриншоте ниже:
>
> ![kafka](./source/exp.png)
>
> ![gif](./source/kf.gif)


## Запуск

>[!IMPORTANT]
> **Запуск через Docker 🐳:**
> ```shell
> cd infostructure
> docker-compose up -d --build
> ```
> **Миграция**
> Auth: users
> ```shell
>   docker exec -i infostructure-postgres-1 psql \
>       -U root \
>       -d postgres \
>       < ../service/auth/migrations/001_init_schema.sql
> ```
>
>Calc: exp & tasks
>
>```shell
>   docker exec -i infostructure-postgres-1 psql \
>       -U root \
>       -d postgres \
>       < ../service/auth/migrations/002_create_expressions.sql
>```
>
> **Проверка полей**
>```shell
>export PGPASSWORD="your pass in .env for pass_pg"                                                 
>psql -h localhost -p 5433 -U root -d postgres
>```
>
>**Вывод пользователей без хэша**
>```shell
> SELECT id, email, created_at FROM users ORDER BY id;
>```
>![img](./source/1.png)
>
>**Вывод пользователей c хэша**
>
>```shell
>SELECT * FROM users;
>```
>![img](./source/2.png)
> **Пример .env** :
> ```env
># ===== VERSION =====
>    ZK_VERSION=3.9.3
>    CONFLUENT_VERSION=7.7.2
>    PG_VERSION=alpine3.20
>    REDIS_VERSION=alpine3.20
>    ES_VERSION=7.17.26
>    KIBANA_VERSION=7.17.26
>    LOGSTASH_VERSION=7.17.26
>
>    # ===== PORTS =====
>    ZK_PORT=2181
>    KAFKA_PORT=9092
>    KAFKA_LOCALHOST_PORT=9093
>    PG_PORT=5433
>    PG_ADMIN_PORT=5050
>    REDIS_PORT=6379
>    KAFKA_UI_PORT=1488
>    AUTH_HTTP_PORT=8080
>
>    # =====DATA BASE========== <3
>    POSTGRES_HOST=localhost
>    POSTGRES_USER=root
>    POSTGRES_PASSWORD=yan2028yan
>    POSTGRES_DB=postgres
>
>    MS_ACCOUNTS_USER = accounts
>    MS_ACCOUNTS_DB = accounts
>    MS_ACCOUNTS_PASSWORD = yan2028yanZ
>
>    REDIS_PASSWORD=yan2027yan
>    REDIS_USER_PASSWORD = yan2027yanZ
>    REDIS_USER = redis
>    REDIS_HOST=localhost
>    REDIS_PORT=6379
>    REDIS_DB=0
>
>
>    # ======= Mirco-serivce=========
>
>    POSTGRES_DSN=postgres://root:yan2028yan@postgres:5432/postgres?sslmode=disable
>
>    REDIS_ADDR=redis:6379
>    KAFKA_BROKERS=kafka:9092
>
>    #================JWT=================
>    JWT_SECRET= ваш jwt-secret
>    ACCESS_TTL=1800        
>   REFRESH_TTL=604800     
>
>
>#======BY FINNIK====
>
> ```
>
> **Генерация jwt-tokens**
>
> ```shell
>  openssl rand -hex 32
> ```
>
>**Git clone**
> ```shell
> git clone https://github.com/golkity/Calc_2.0.git
> ```


![logo_out](source/logo_out.png)
