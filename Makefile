TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=cobbler
COBBLER_SERVER_URL=http://localhost:8081/cobbler_api

default: build

build: fmtcheck
	go install

release-test:
	goreleaser release --snapshot --clean

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

#testacc: fmtcheck
#	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
# Run acceptance tests
.PHONY: testacc
testacc:
	@COBBLER_VERSION=v3.3.0 sh -c "'./docker/start.sh' $(COBBLER_SERVER_URL)"
	TF_LOG=TRACE TF_ACC_LOG=TRACE TF_LOG_PATH_MASK="test-%s.log" TF_ACC=1 TF_ACC_PROVIDER_NAMESPACE=cobbler COBBLER_URL=$(COBBLER_SERVER_URL) COBBLER_USERNAME=cobbler COBBLER_PASSWORD=cobbler go test -v -p 1 -coverprofile="coverage.out" -covermode="atomic" './...'

.PHONY: docs
docs:
	@tfplugindocs

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build release-test test testacc vet fmt fmtcheck errcheck test-compile

