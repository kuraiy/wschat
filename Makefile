# ── wschat Makefile ───────────────────────────────────────────────
# Все значения берутся из .env (в гит не коммитится).
# Сам Makefile секретов не содержит — его коммитить безопасно.

# Подтягиваем переменные из .env и экспортируем в окружение под-команд.
include .env
export

# DB_URL собирается из компонентов .env.
# PGHOST=localhost + проброшенный порт => goose с ХОСТА бьёт в докерную БД.
DB_URL := postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(DB_NAME)?sslmode=disable

# goose гоняем через go run — не нужно ставить его глобально.
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "$(DB_URL)"

.DEFAULT_GOAL := help

# ── повседневный цикл: инфра в докере, app через go run ────────────

## infra: поднять postgres + redis в фоне (app НЕ трогаем)
infra:
	docker compose up -d postgres redis

## infra-down: остановить postgres + redis (данные в volume целы)
infra-down:
	docker compose stop postgres redis

## run: запустить app локально (ходит в контейнеры по localhost)
run:
	go run ./cmd

## dev: всё для локальной разработки одной командой — инфра + миграции + app
dev: infra migrate-up run

# ── миграции (goose, с хоста в докерную БД) ────────────────────────

## migrate-up: накатить все миграции
migrate-up:
	$(GOOSE) up

## migrate-down: откатить последнюю миграцию
migrate-down:
	$(GOOSE) down

## migrate-status: показать состояние миграций
migrate-status:
	$(GOOSE) status

## migrate-reset: откатить ВСЕ миграции (осторожно)
migrate-reset:
	$(GOOSE) reset

# ── полный прогон в докере (app тоже в контейнере) ─────────────────

## up: собрать и поднять ВЕСЬ стек в докере
up:
	docker compose up -d --build

## down: погасить стек, volume с данными ОСТАВИТЬ
down:
	docker compose down

## clean: погасить стек и УДАЛИТЬ данные БД (сброс volume) — осторожно!
clean:
	docker compose down -v

# ── наблюдение / отладка ───────────────────────────────────────────

## ps: статус контейнеров
ps:
	docker compose ps

## logs: живые логи app
logs:
	docker compose logs -f --tail=50 app

## logs-all: живые логи всех сервисов
logs-all:
	docker compose logs -f --tail=50

## psql: открыть psql внутри контейнера postgres
psql:
	docker compose exec postgres psql -U $(PGUSER) -d $(DB_NAME)

## help: список команд
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | awk -F': ' '{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: infra infra-down run dev migrate-up migrate-down migrate-status \
	migrate-reset up down clean ps logs logs-all psql help
