# ------- Constants -------
# Директория для исполняемых файлов
LOCAL_BIN := $(CURDIR)/bin

# Директория для сгенерированных прото-файлов
OUT_PATH := $(CURDIR)/pkg/protogen

GOBIN := $(LOCAL_BIN)
export GOBIN

# для установки на другие OS нужно заменить "win64" см. https://github.com/protocolbuffers/protobuf/releases/
PROTOC_VERSION := protoc-31.1-win64
export PROTOC_VERSION

PROTO_FILES := $(wildcard api/**/*.proto)

# ------- gRPC-deps -------
.dependencies/protoc: ## Установить бинарные зависимости и protobuf-плагины (win64)
	echo "Installing bin dependencies"
	curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v31.1/$(PROTOC_VERSION).zip
	unzip -o $(PROTOC_VERSION).zip -d $(LOCAL_BIN)
	rm -f $(PROTOC_VERSION).zip

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@mkdir -p .dependencies
	@touch $@
	@echo "All binaries installed"

# Скачивание протобаф плагинов
.dependencies/vendor-proto/validate:
	echo "Installing validate"
	git clone -b main --single-branch --depth=2 --filter=tree:0 https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp
	cd vendor.protogen/tmp && git sparse-checkout init --no-cone && git sparse-checkout set validate && git checkout
	mv vendor.protogen/tmp/validate vendor.protogen/
	rm -rf "vendor.protogen/tmp"
	@mkdir -p .dependencies/vendor-proto
	@touch $@

.dependencies/vendor-proto/google/api:
	echo "Installing google/api"
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 https://github.com/googleapis/googleapis vendor.protogen/googleapis
	cd vendor.protogen/googleapis && git sparse-checkout init --no-cone && git sparse-checkout set google/api && git checkout
	mv vendor.protogen/googleapis/google vendor.protogen/google
	rm -rf "vendor.protogen/googleapis"
	@mkdir -p .dependencies/vendor-proto/google
	@touch $@

.dependencies/vendor-proto/protoc-gen-openapiv2/options:
	echo "Installing protoc-gen-openapiv2/options"
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem
	cd vendor.protogen/grpc-ecosystem && git sparse-checkout init --no-cone && git sparse-checkout set protoc-gen-openapiv2/options && git checkout
	mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2 vendor.protogen/protoc-gen-openapiv2
	rm -rf  "vendor.protogen/grpc-ecosystem"
	@mkdir -p .dependencies/vendor-proto/protoc-gen-openapiv2
	@touch $@

vendor-deps: .dependencies/vendor-proto/validate .dependencies/vendor-proto/google/api .dependencies/vendor-proto/protoc-gen-openapiv2/options
	@echo "All vendor dependencies installed"

grpc-deps: vendor-deps .dependencies/protoc

# ------- Proto compilation -------
# Генерирование гошных имплементаций прото файлов
proto-generate: $(PROTO_FILES) ## Сгенерировать protobuf
	mkdir -p $(OUT_PATH)
	$(LOCAL_BIN)/bin/protoc --proto_path=api --proto_path=vendor.protogen \
		--go_out=$(OUT_PATH) --go_opt=paths=source_relative --plugin protoc-gen-go="${GOBIN}/protoc-gen-go.exe" \
		--go-grpc_out=$(OUT_PATH) --go-grpc_opt=paths=source_relative --plugin protoc-gen-go-grpc="${GOBIN}/protoc-gen-go-grpc.exe" \
		--validate_out="lang=go,paths=source_relative:$(OUT_PATH)" --plugin protoc-gen-validate="$(LOCAL_BIN)/protoc-gen-validate.exe" \
		--grpc-gateway_out=$(OUT_PATH) --grpc-gateway_opt=paths=source_relative --plugin protoc-gen-grpc-gateway="$(LOCAL_BIN)/protoc-gen-grpc-gateway.exe" \
		--openapiv2_out=$(OUT_PATH) --plugin=protoc-gen-openapiv2="$(LOCAL_BIN)/protoc-gen-openapiv2.exe" \
		api/*/*.proto
	go mod tidy

.PHONY: up
up:
	docker-compose up