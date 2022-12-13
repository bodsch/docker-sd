export CGO_ENABLED:=0
export GOARCH:=amd64
export PATH:=$(PATH):$(PWD)

LDFLAGS=-X main.Version=$(shell $(CURDIR)/build/git-version.sh) -s
GOFLAGS="-a"

# TODO: Fix this on windows.
ALL_SRC := $(shell find . \
		-name '*.go' \
		-not -path './vendor/*' \
    	-not -path '*/gen-go/*' \
    	-type f | sort)
ALL_PKGS := $(shell go list $(sort $(dir $(ALL_SRC))))

GOFMT=gofmt
GOLINT=golint
GOVET=go vet
GO_BUILD=go build

README_FILES := $(shell find . -name '*README.md' | sort | tr '\n' ' ')

all-pkgs:
	@echo $(ALL_PKGS) | tr ' ' '\n' | sort

.PHONY: fmt
fmt:
	@FMTOUT=`$(GOFMT) -s -l $(ALL_SRC) 2>&1`; \
	if [ "$$FMTOUT" ]; then \
	        echo "$(GOFMT) FAILED => gofmt the following files:\n"; \
	        echo "$$FMTOUT\n"; \
	        exit 1; \
	else \
	    echo "Fmt finished successfully"; \
	fi

.PHONY: lint
lint:
	@LINTOUT=`$(GOLINT) $(ALL_PKGS) | grep -v $(TRACE_ID_LINT_EXCEPTION) | grep -v $(TRACE_OPTION_LINT_EXCEPTION) 2>&1`; \
	if [ "$$LINTOUT" ]; then \
	        echo "$(GOLINT) FAILED => clean the following lint errors:\n"; \
	        echo "$$LINTOUT\n"; \
	        exit 1; \
	else \
	    echo "Lint finished successfully"; \
	fi

.PHONY: build
build:
	@BUILD_OUT=`$(GO_BUILD) $(GOFLAGS) -ldflags "$(LDFLAGS)" -o docker_sd .` ; \
	if [ "$$BUILD_OUT" ]; then \
		echo -e "$(GO_BUILD) FAILED =>\n"; \
		echo -e "$$BUILD_OUT\n" ;\
		exit 1 ;\
	else \
	  	echo "build finished successfully"; \
        exit 0 ;\
	fi
