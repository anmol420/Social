include .env
MIGRATION_PATH=./cmd/migrate/migrations

.PHONY: migrate-version
migrate-version:
	@migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_ADDR) version

.PHONY: migrate-force
migrate-force:
	@bash -c 'read -p "Enter version to force to: " version; \
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_ADDR) force $$version'

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATION_PATH) -database=$(DATABASE_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATION_PATH) -database=$(DATABASE_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

%:
	@:
