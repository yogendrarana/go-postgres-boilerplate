include .env
export

migrate-up:
	goose up

migrate-down:
	goose down

migrate-status:
	goose status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=your_migration_name"; \
	else \
		goose create $(name) sql; \
	fi
