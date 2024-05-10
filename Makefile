GO		:= go

HAS_GOMODULES := $(shell go help mod why 2> /dev/null)

ifdef HAS_GOMODULES
export GO111MODULE=on
export GOFLAGS = -mod=vendor
else
$(warn vendor import can only be done on  go 1.11+ which supports go modules)
endif

.PHONY: all build install clean test format vet golint lint errcheck vendor


# Tools
#
# In module mode, 'go get' has a side-effect of updating the go.mod
# file.  We do not want to update go.mod when installing tools.
# As a workaround, when installing a tool, cd to /tmp and turn off
# module mode.  This should be solved in:
#   https://github.com/golang/go/issues/30515
#   https://github.com/golang/go/issues/24250

fetch-tools:
	mkdir -p tools
	(cd tools && $(GO) get -u golang.org/x/lint/golint)
	(cd tools && $(GO) get -v github.com/kisielk/errcheck)
	(cd tools && $(GO) get -u github.com/vbatts/git-validation)
	(cd tools && $(GO) get k8s.io/client-go/rest@v11.0.1-0.20190409021438-1a26190bd76a+incompatible)
	(cd tools && $(GO) get get github.com/libopenstorage/openstorage/api@v8.0.0+incompatible)
	(cd tools && $(GO) get get golang.org/x/oauth2/internal@v0.0.0-20190604053449-0f29369cfe45)
	(cd tools && $(GO) get get google.golang.org/grpc@v1.26.0)
	(cd tools && $(GO) get get google.golang.org/grpc/internal/transport@v1.26.0)

# Deliverables

fmt:
	$(GO) fmt ./... | wc -l | grep 0

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

git-validation:
	git-validation -run DCO,short-subject

lint:
	golint `go list ./...`

errcheck:
	errcheck -verbose -blank ./...

vendor:
	@echo "Updating vendor tree"
	go mod vendor

vendor-tidy:
	@echo "Removing unused files in vendor tree"
	go mod tidy
