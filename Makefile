.PHONY: all
all: jump check

.PHONY: check
check:
	go test ./...

.PHONY: test
test: check

.PHONY: jump
jump:
	go build

.PHONY: shellcheck
shellcheck:
	shellcheck jump.sh
