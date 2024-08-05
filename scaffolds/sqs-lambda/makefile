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

# ── Lambda ──────────────────────────────────────────────────────────────────────

.PHONY: lambda_build
lambda_build:
	sam build --no-cached


.PHONY: lambda_local_create_users
lambda_local_create_users: db_up_d lambda_build
	sam local invoke --event ./events/create_users.json --env-vars env.local.json
	docker-compose down
