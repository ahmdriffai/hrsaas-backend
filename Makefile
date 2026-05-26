APP_NAME	= hrsaas-backend
BUILD_DIR	= ./bin

DB_USER		= postgres
DB_PASS		= postgres
DB_HOST		= localhost
DB_PORT     = 5432
DB_NAME     = hr_saas
DB_SSLMODE  = disable

DB_URL = "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

MIGRATE = migrate -path db/migrations -database $(DB_URL)

.PHONY: run build dev tidy seed migrate-up migrate-down migrate-up-n migrate-down-n migrate-force migrate-version migrate-create

## ─── APP ────────────────────────────────────────────────────────────────────

## Jalankan aplikasi langsung
run:
	go run cmd/web/main.go

## Build binary ke ./bin/
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

## Hot-reload dengan air (otomatis restart saat kode berubah)
dev:
	$(shell go env GOPATH)/bin/air

## Tidy dependencies
tidy:
	go mod tidy

## Jalankan seeder
seed:
	go run db/seeder/main.go

## ─── MIGRATION ──────────────────────────────────────────────────────────────

## Jalankan semua migrasi
migrate-up:
	$(MIGRATE) up

## Rollback semua migrasi
migrate-down:
	$(MIGRATE) down

## Jalankan N migrasi ke atas   → make migrate-up-n N=1
migrate-up-n:
	$(MIGRATE) up $(N)

## Rollback N migrasi ke bawah  → make migrate-down-n N=1
migrate-down-n:
	$(MIGRATE) down $(N)

## Paksa set versi migrasi (untuk fix dirty state) → make migrate-force V=20260514023530
migrate-force:
	$(MIGRATE) force $(V)

migrate-drop:
	$(MIGRATE) drop

## Lihat versi migrasi saat ini
migrate-version:
	$(MIGRATE) version

## Buat file migrasi baru → make migrate-create NAME=create_table_xxx
migrate-create:
	migrate create -ext sql -dir db/migrations -format "20060102150405" $(NAME)