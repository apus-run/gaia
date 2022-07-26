user	:=	$(shell whoami)
rev 	:= 	$(shell git rev-parse --short HEAD)

# GOBIN > GOPATH > INSTALLDIR
GOBIN	:=	$(shell echo ${GOBIN} | cut -d':' -f1)
GOPATH	:=	$(shell echo $(GOPATH) | cut -d':' -f1)
BIN		:= 	""

# golangci-lint
LINTER := bin/golangci-lint

# check GOBIN
ifneq ($(GOBIN),)
	BIN=$(GOBIN)
else
	# check GOPATH
	ifneq ($(GOPATH),)
		BIN=$(GOPATH)/bin
	endif
endif

$(LINTER):
	curl -SL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

all:
	#@cd cmd/gaia && go build && cd - &> /dev/null
	@cd cmd/protoc-gen-go-gin && go build && cd - &> /dev/null

.PHONY: install
install: all
ifeq ($(user),root)
#root, install for all user
	#@cp ./cmd/gaia/gaia /usr/bin
	@cp ./cmd/protoc-gen-go-gin/protoc-gen-go-gin /usr/bin
else
#!root, install for current user
	$(shell if [ -z $(BIN) ]; then read -p "Please select installdir: " REPLY; mkdir -p $${REPLY};\
    cp ./cmd/protoc-gen-go-gin/protoc-gen-go-gin $${REPLY}/;\
	cp ./cmd/protoc-gen-go-gin/protoc-gen-go-gin $(BIN);fi)
endif
	@which protoc-gen-go &> /dev/null || go get google.golang.org/protobuf/cmd/protoc-gen-go
	@which protoc-gen-go-grpc &> /dev/null || go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@which protoc-gen-validate  &> /dev/null || go get github.com/envoyproxy/protoc-gen-validate
	@echo "install finished"

.PHONY: uninstall
uninstall:
	$(shell for i in `which -a protoc-gen-go-grpc | grep -v '/usr/bin/protoc-gen-go-gin' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	$(shell for i in `which -a protoc-gen-validate | grep -v '/usr/bin/protoc-gen-go-gin' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	@echo "uninstall finished"
