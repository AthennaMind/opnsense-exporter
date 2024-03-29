BINARY_NAME=opnsense-exporter-local

.PHONY: default
default: run-test

sync-vendor:
	go mod tidy
	go mod vendor

local-run:
	go build \
	-tags osusergo,netgo \
	-ldflags '-w -extldflags "-static" -X main.version=local-test' \
	-v -o ${BINARY_NAME}

	./${BINARY_NAME} --log.level="debug" \
		--log.format="logfmt" \
		--web.telemetry-path="/metrics" \
		--web.listen-address=":8080" \
		--runtime.gomaxprocs=4 \
		--exporter.instance-label="opnsense-local1" \
		--exporter.disable-arp-table \
		--exporter.disable-cron-table \
		--opnsense.protocol="https" \
		--opnsense.address="${OPS_ADDRESS}" \
		--opnsense.api-key="${OPS_API_KEY}" \
		--opnsense.api-secret="${OPS_API_SECRET}" \
		--web.disable-exporter-metrics \
		
test:
	go test ./...

clean:
	gofmt -s -w $(shell find . -type f -name '*.go'| grep -v "/vendor/\|/.git/")
	go clean
	rm ./${BINARY_NAME}

lint:
	gofmt -s -w $(shell find . -type f -name '*.go'| grep -v "/vendor/\|/.git/")
	golangci-lint run --fix
