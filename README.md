# go-shop

Учебный мини-проект: два микросервиса на Go (стандартная библиотека).

## Services

- `catalog` — каталог товаров (позже: CRUD)
- `order` — оформление заказа (позже: создание заказа, статусы)

Пока что оба сервиса умеют:
- `GET /healthz`
- `GET /version`

## Quickstart

Запуск в двух терминалах:

```bash
make run-catalog
```

```bash
make run-order
```

Проверка:

```bash
curl -sS localhost:8081/healthz
curl -sS localhost:8082/healthz
curl -sS localhost:8081/version
curl -sS localhost:8082/version
```

## Workflow (GitHub-style)

- `main` — только через PR
- feature branches — `feat/<short-name>`
- fixes — `fix/<short-name>`

## Roadmap (итеративно)

- [ ] базовые модели (Product, Order)
- [ ] простое in-memory хранилище
- [ ] эндпоинты каталога: list/get/create
- [ ] эндпоинты заказа: create/get/status
- [ ] интеграция order -> catalog (проверка товара)
- [ ] docker-compose
- [ ] OpenAPI спецификация
- [ ] минимальные тесты хендлеров
