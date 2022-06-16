REVISION := $(shell git describe --always)
LDFLAGS	 := -ldflags="-X \"main.Revision=$(REVISION)\""

.PHONY: help

name		:= famili-api
linux_name	:= $(name)-linux-amd64
darwin_name	:= $(name)-darwin-amd64
darwin_arm_name	:= $(name)-darwin-arm64

go_version := $(shell cat $(realpath .go-version))
go_bindir  := ~/.go_binary/$(go_version)
goproxy    := direct
goprivate  := $(PWD)
gosumdb    := off
goroot     := $(go_bindir)/go
go         := $(go_bindir)/go/bin/go
arch       := $(shell arch)

db_name      := famili
db_user_name := famili-api

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

dist: build/docker ## create .tar.gz linux & darwin to /bin
	cd bin && tar zcvf $(linux_name).tar.gz $(linux_name) && rm -f $(linux_name)
	cd bin && tar zcvf $(darwin_name).tar.gz $(darwin_name) && rm -f $(darwin_name)
	cd bin && tar zcvf $(darwin_arch_name).tar.gz $(darwin_arch_name) && rm -f $(darwin_arch_name)

clean: docker_compose/down_all ## このMakefileで利用したファイルをクリアにする
	rm -rf $(go_bindir)
	rm -f $(name)
	rm -rf ./db-data

build: go/install ## build
	GOROOT=$(goroot) $(go) build -o bin/$(name) cmd/$(name)/*.go

build/cross: go/install ## create to build for linux & darwin to bin/
	GOOS=linux GOARCH=amd64 $(go) build -o bin/$(linux_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=linux GOARCH=arm64 $(go) build -o bin/$(linux_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=darwin GOARCH=amd64 $(go) build -o bin/$(darwin_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=darwin GOARCH=arm64 $(go) build -o bin/$(darwin_arm_name) $(LDFLAGS) cmd/$(name)/*.go

run: ## go run
	GOROOT=$(goroot) $(go) run cmd/$(name)/main.go -c examples/config.toml

run/binary: ## run binary
	./bin/$(name) -c examples/config.toml

### go operation
go/install: file         = go.tar.gz
go/install: download_url = https://golang.org/dl/go$(go_version).darwin-$(arch).tar.gz
go/install:
# If you have a different version, delete it.
	@if [ -f $(go) ]; then \
		$(go) version | grep -q "$(go_version)" || rm -rf $(go_bindir)/go; \
	fi

# If the file is not there, download it.
	@if [ ! -f $(go) ]; then \
  		mkdir -p $(go_bindir) && \
		curl -L -fsS --retry 2 -o $(file) $(download_url) && \
		tar zxvf $(file) -C $(go_bindir) && rm -f $(file); \
	fi

go/get: require_package ## 特定のパッケージをgo getする
	GOPROXY=$(goproxy) GOSUMDB=$(gosumdb) GOPRIVATE=$(goprivate) GOROOT=$(goroot) $(go) get -u $(PACKAGE)

go/get_u: ## パッケージを全部更新する
	GOPROXY=$(goproxy) GOSUMDB=$(gosumdb) GOPRIVATE=$(goprivate) GOROOT=$(goroot) $(go) get -u

go/mod_tidy: go/install ## 不要なgo packageを削除する
	GOPROXY=$(goproxy) GOSUMDB=$(gosumdb) GOPRIVATE=$(goprivate) GOROOT=$(goroot) $(go) mod tidy -v

go/test: go/install ## go test
	GOROOT=$(goroot) $(go) test ./... -parallel 10

### docker
docker/build: ## docker build
	DOCKER_BUILDKIT=1 docker build -f docker/api/Dockerfile . -t $(name) --secret id=gitconfig,src=$(HOME)/.gitconfig

### docker_compose
docker_compose/up: ## compose起動
	COMPOSE_DOCKER_CLI_BUILD=1 docker-compose up --build

docker_compose/down: ## compose停止
	COMPOSE_DOCKER_CLI_BUILD=1 docker-compose down

docker_compose/down_f: ## composeではないけど、dockerで強制停止する. conflict対策
	docker ps -a | grep "$(name)" | awk '{print $$1}' | xargs docker rm -f

docker_compose/down_all: ## compose停止 + 全てを初期化する
	COMPOSE_DOCKER_CLI_BUILD=1 docker-compose down --rmi all --volumes --remove-orphans

docker_compose/rebuild: ## appだけbuildし直す
	COMPOSE_DOCKER_CLI_BUILD=1 docker-compose build app
	COMPOSE_DOCKER_CLI_BUILD=1 docker-compose up -d app

### migrate
migrate := -path "/migrations" -database "mysql://$(db_user_name):password@tcp(db:3306)/$(db_name)"

migrate/up: ## migration. docker compose up後に実行できる
	docker compose run migrate $(migrate) up

migrate/down: ## migrationのrollback. docker compose up後に実行できる
	docker compose run migrate $(migrate) down

migrate/create: require_migrate_file ## migrationファイル作成. migrations/ にup/downが作成される
	docker compose run migrate create -ext mysql -dir migrations $(MIGRATE_FILE)

openapi/gen: ## openapiのファイルをgenerateする
	oapi-codegen -generate chi-server -o openapi/openapi.gen.go -package openapi openapi.yaml


### require
### 要求するターゲットをここで設定
require_package:
ifeq ($(PACKAGE),)
	@echo -e "$(RED)you must set a argument PACKAGE=xxx$(NC)"
	@exit 1
endif


require_migrate_file:
ifeq ($(MIGRATE_FILE),)
	@echo -e "$(RED)you must set a argument MIGRATE_FILE=xxx$(NC)"
	@exit 1
endif
