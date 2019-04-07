.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## installs dependencies
	go get github.com/cortesi/modd/cmd/modd
	go get gopkg.in/alecthomas/gometalinter.v1

.PHONY: install
install: ## installs tmplgen
	go install github.com/unders/tmplgen

.PHONY: start
start: ## starts development environment for tmplserv
	modd -f modd.conf