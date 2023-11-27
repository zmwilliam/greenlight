help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

run: ## run the cmd/api application
	go run ./cmd/api

migrations-up: confirm ## apply all up datadase migrations
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN} up

migrations-down: confirm ## apply all down datadase migrations
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${DB_DSN} down

migrations-new: ## name=$1: create a new database migration
	@echo 'Creating migration files for ${names}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: help confirm run migrations-up migrations-down migrations-new
