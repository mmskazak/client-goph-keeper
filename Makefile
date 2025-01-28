.PHONY: lint
lint:
	# | jq > ./golangci-lint/report.json
	golangci-lint run --fix -c .golangci.yml > golangci-lint/report-unformatted.json
	goimports -local mmskazak -w .

.PHONY: lint-clean
lint-clean:
	sudo rm -rf ./golangci-lint

.PHONY: test
test:
	go test ./...

