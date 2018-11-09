NAME     := dayone-to-quiver
VERSION  := v0.0.1
ifeq "$(shell git status | grep -E 'nothing to commit, working (tree|directory) clean' | wc -l | sed -e 's/ //g')" "0"
  REVISION := $(shell git rev-parse --short HEAD)-dev
else
  REVISION := $(shell git rev-parse --short HEAD)
endif
BUILD_DATE  := $(shell date -u +%Y%m%d.%H%M%S)

GO := /usr/local/go/bin/go
GOPATH = $(shell git rev-parse --show-toplevel)
#SRCS    := $(shell find -E . -type f -regex '.*(go|html)')
SRCS    := $(shell find -E ${GOPATH} -type f -regex '.*(go|html)')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
GOBIN := ../../../../bin/go-assets-builder


$(NAME): $(SRCS)
#	$(GOBIN) --output=resources.go resources
	GOPATH=$(GOPATH) $(GO) build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o $(NAME)


.PHONY: clean
clean:
	rm -f $(NAME)*

.PHONY: cross-build
#;cross-build: deps
cross-build:
#	$(GOBIN) --output=resources.go resources
	for os in darwin linux; do \
	    for arch in amd64; do \
	        GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 GOPATH=$(GOPATH) $(GO) build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
	    done; \
	done