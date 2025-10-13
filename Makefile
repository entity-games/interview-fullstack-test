.PHONY: deploy_app deploy_server run_stack

LATEST_TAG := $(shell git tag --sort=-v:refname | grep -E '^[0-9]+\.[0-9]+\.[0-9]+' | head -n 1)
GIT_HASH   := $(shell git show --summary | grep commit | head -n 1 | cut -c 8-15)
DOCKER_TAG := $(shell cat version.txt 2>/dev/null || echo "ci-${LATEST_TAG}-${GIT_HASH}")
DOCKER_ECR := 654654477305.dkr.ecr.eu-west-1.amazonaws.com/interview
PWD	       := $(shell pwd)

start_stack:
	docker compose -f deploy/dev-stack.yml up --remove-orphans -d

start_server_go:
	air -c .air.conf

start_app_web:
	npm run watch

init_db:
	redis-cli -h 127.0.0.1 FLUSHALL
	PGPASSWORD=example psql -h 127.0.0.1 -U postgres -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'interview' AND pid <> pg_backend_pid()"
	PGPASSWORD=example psql -h 127.0.0.1 -U postgres -d postgres -c "DROP DATABASE IF EXISTS interview;"
	PGPASSWORD=example psql -h 127.0.0.1 -U postgres -d postgres -c "CREATE DATABASE interview;"

migrate_db:
	migrate -database postgres://postgres:example@127.0.0.1:5432/interview?sslmode=disable -path ${PWD}/cmd/server/migrations up

migrate: init_db migrate_db

run_test:
	go test -json -v ./cmd/tests/... -race -count=1 | gotestfmt
	go tool cover -func=coverage.out

test: migrate run_test

start_all:
	make start_stack & make start_server_go & make start_app_web & wait
