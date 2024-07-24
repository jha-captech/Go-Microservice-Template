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

.PHONY: lambda_local_api
lambda_local_api: db_up lambda_build
	sam local start-api -p 8080 --env-vars env.local.json
	make db_down

.PHONY: lambda_local_list_users
lambda_local_list_users: db_up_d lambda_build
	sam local invoke --event ./events/list_users.json --env-vars env.local.json
	make db_down

.PHONY: lambda_local_update_users
lambda_local_update_users: db_up_d lambda_build
	sam local invoke --event ./events/update_user.json --env-vars env.local.json
	make db_down
