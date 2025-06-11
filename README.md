# Балансировщик нагрузки (Golang)

## Описание

Проект представляет собой простой балансировщик нагрузки на Go. В работе использовался алгоритм Round Robin.

## Требования для запуска

- Docker 20.10+
- Docker Compose 1.29
- Go 1.20+ — для локальной разработки и запуска тестов (unit-тесты находятся в каталогах `internal/health` и `pkg/lb`)

---

## Запуск проекта

Собрать и запустить контейнеры из папки проекта:

```bash
docker-compose up --build
```

Доступ:

- Балансировщик: http://localhost:8080
- Бэкенд 1: http://localhost:8081
- Бэкенд 2: http://localhost:8082
- Бэкенд 3: http://localhost:8083

---

## Конфигурация

Файл: `cmd/config.json`:

```json
{
  "listen_port": "8080",
  "health_path": "/health",
  "backends": [
    "http://backend1:8081",
    "http://backend2:8082",
    "http://backend3:8083"
  ]
}
```

Параметры health check задаются в `internal/health/checker.go`:

- `timeoutSec` — таймаут запроса
- `intervalSec` — интервал проверки

---

## API Endpoints

### Балансировщик

- `GET /` — основной эндпоинт (`{"status":"ok"}`)
- Все остальные запросы проксируются на бэкенды

### Бэкенды

- `GET /` — возвращает текст с номером порта
- `GET /health` — проверка работоспособности (возвращает HTTP 200)

---

## Логирование

В логах балансировщика отображаются:

- Конфигурация
- Результаты health check
- Ошибки запросов

**Пример логов:**

```
loadbalancer-1  | Configured backend: http://backend1:8081
loadbalancer-1  | Backend http://backend1:8081 is healthy
loadbalancer-1  | Unhealthy: http://backend2:8082 (Connection refused)
```

## Остановка

```bash
docker-compose down
```
