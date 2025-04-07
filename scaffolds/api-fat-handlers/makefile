# ── Database ────────────────────────────────────────────────────────────────────

.PHONY: db_up
db_up:
	docker-compose up postgres

.PHONY: db_up_d
db_up_d:
	docker-compose up postgres -d

.PHONY: db_down
db_down:
	docker-compose down postgres

# ── API ─────────────────────────────────────────────────────────────────────────

.PHONY: mockery
mockery:

.PHONY: swagger
swagger:
	swag init \
		--generalInfo "./../../cmd/api/main.go" \
		--dir "./internal/handlers" \
		--output "./internal/swagger/docs" \
		--parseInternal

.PHONY: api
api: swagger
	docker-compose up
