BINARY_NAME=opnsense-exporter-local

.PHONY: default
default:
	go build \
	-tags osusergo,netgo \
	-v -o ${BINARY_NAME}

sync-vendor:
	go mod tidy
	go mod vendor

local-run: default
	./${BINARY_NAME} --log.level="debug" \
		--log.format="logfmt" \
		--web.telemetry-path="/metrics" \
		--web.listen-address=":$(or $(OPS_EXPORTER_PORT), 8080)" \
		--runtime.gomaxprocs=4 \
		--exporter.instance-label="$(or $(OPS_INSTANCE), opnsense-local1)" \
		--opnsense.protocol="https" \
		--opnsense.address="${OPS_ADDRESS}" \
		--opnsense.api-key="${OPS_API_KEY}" \
		--opnsense.api-secret="${OPS_API_SECRET}" \
		--web.disable-exporter-metrics \
		$(if $(OPS_ADDITIONAL_ARGS),"${OPS_ADDITIONAL_ARGS}")

test:
	go test ./...

clean:
	gofmt -s -w $(shell find . -type f -name '*.go'| grep -v "/vendor/\|/.git/")
	go clean
	rm ./${BINARY_NAME}

lint:
	gofmt -s -w $(shell find . -type f -name '*.go'| grep -v "/vendor/\|/.git/")
	golangci-lint run --fix
