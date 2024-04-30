-include .env
export

CURRENT_DIR=$(shell pwd)
APP=content_service
CMD_DIR=./cmd

.DEFAULT_GOAL = build

# build for current os
.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# build for linux amd64
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# run service
.PHONY: run
run:
	go run ${CMD_DIR}/main.go

# migrate
.PHONY: migrate
migrate-up:
	migrate -source file://migrations -database postgresql://postgres:mubina2007@localhost:5432/companydb?sslmode=disable up

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_jobs_table

migrate-down:
	migrate -source file://migrations -database postgresql://postgres:mubina2007@localhost:5432/companydb?sslmode=disable down

migrate-force:
	migrate -source file://migrations -database postgresql://postgres:mubina2007@localhost:5432/companydb?sslmode=disable force 1

# proto
.PHONY: proto-gen
proto-gen:
	chmod +x ./scripts/gen-proto.sh
	./scripts/gen-proto.sh

# git submodule init 	
.PHONY: pull-proto
pull-proto:
	git submodule update --init --recursive

# go generate	
.PHONY: go-gen
go-gen:
	go generate ./...

# run test
test:
	go test -v -cover -race ./internal/...

# -------------- for deploy --------------
build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

.PHONY: consumer
consumer:
	go run cmd/main.go consumer job_create_consumer
