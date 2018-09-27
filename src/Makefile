# defining environment variables

TARGETNAME=microservice-adapter-mqtt

LINTERCOMMAND=gometalinter.v2

# Packages to be tested
packages =  ./adapter \
			./config \
			./logger \
			./mqtt

# global commands
.PHONY: all
all: dependencies test code-quality build

.PHONY: dependencies
dependencies:
	@echo "Installing dependencies"
	@glide cc
	@glide install

.PHONY: code-quality
code-quality:
	$(LINTERCOMMAND) --vendor --tests --skip=mocks \
		--exclude='vendor/*' --exclude='_test.go' --exclude='mocks/*' \
		--disable=megacheck --disable=ineffassign \
		--disable=gotype --disable=errcheck --disable=gas --disable=dupl --disable=aligncheck \
		--deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o  $(TARGETNAME)

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o $(TARGETNAME)

.PHONY: test
test:
	@$(foreach package,$(packages), \
		set -e; \
		go test -coverprofile $(package)/cover.out -covermode=count $(package);)

.PHONY: example
example:
	@echo "Run Example"
	@cd ../ && \
	docker-compose build && \
	docker-compose up
