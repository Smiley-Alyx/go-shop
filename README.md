# go-shop

Мини-проект: два микросервиса на Go (стандартная библиотека).

## Требования

- Go 1.22+
- Docker (опционально, для сборки образов и `docker compose`)

## Services

- `catalog` — каталог товаров
- `order` — оформление заказа

## Эндпоинты

### catalog (порт по умолчанию: 8081)

- `GET /healthz` -> `ok`
- `GET /version` -> версия из env `VERSION`
- `GET /products` -> список товаров (JSON)
- `GET /products/{id}` -> товар (JSON) или `404`
- `POST /products` -> создать товар
  - body: `{ "name": "...", "price": 123 }`
  - response: созданный товар `{ "id": 1, "name": "...", "price": 123 }`

### order (порт по умолчанию: 8082)

- `GET /healthz` -> `ok`
- `GET /version` -> версия из env `VERSION`
- `POST /orders` -> создать заказ
  - body: `{ "items": [ { "product_id": 1, "qty": 2 } ] }`
  - response: созданный заказ `{ "id": 1, "status": "new", "items": [...], "total": 200 }`
- `GET /orders/{id}` -> заказ (JSON) или `404`
- `GET /orders/{id}/status` -> `{ "status": "new" }`

## Конфигурация (env)

### catalog

- `PORT` (по умолчанию `8081`)
- `VERSION` (по умолчанию `0.0.0`)

### order

- `PORT` (по умолчанию `8082`)
- `VERSION` (по умолчанию `0.0.0`)
- `CATALOG_URL` (по умолчанию `http://localhost:8081`) — базовый URL сервиса `catalog`

## Quickstart

Запуск в двух терминалах:

```bash
make run-catalog
```

```bash
make run-order
```

Проверка базовых ручек:

```bash
curl -sS localhost:8081/healthz
curl -sS localhost:8082/healthz
curl -sS localhost:8081/version
curl -sS localhost:8082/version
```

Пример:

```bash
curl -sS -X POST localhost:8081/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"tea","price":100}'

curl -sS -X POST localhost:8082/orders \
  -H 'Content-Type: application/json' \
  -d '{"items":[{"product_id":1,"qty":2}]}'

curl -sS localhost:8082/orders/1/status
```

## Docker

Сборка образов:

```bash
docker build -t go-shop-catalog:dev services/catalog
docker build -t go-shop-order:dev services/order
```

Запуск через compose:

```bash
docker compose up --build
```

## OpenAPI

- `openapi/catalog.yaml`
- `openapi/order.yaml`

## Make

```bash
make tidy
make test
make build
```

## Тесты

Тесты запускаются через:

```bash
make test
```

Сейчас покрыто:

- `catalog`: `/healthz`, `/version`, `/products` (list/get/create)
- `order`: `/healthz`, `/version`, `/orders` (create/get/status)

## Workflow (GitHub-style)

- `main` — только через PR
- feature branches — `feat/<short-name>`
- fixes — `fix/<short-name>`

## Roadmap (итеративно)

- [x] базовые модели (Product, Order)
- [x] простое in-memory хранилище
- [x] эндпоинты каталога: list/get/create
- [x] эндпоинты заказа: create/get/status
- [x] интеграция order -> catalog (проверка товара)
- [x] docker-compose
- [x] OpenAPI спецификация
- [x] минимальные тесты хендлеров
- [ ] статусы заказа: переходы (paid/cancelled)
- [ ] простая валидация и ошибки
- [ ] нормальная конфигурация `order -> catalog` без хардкода
