
INSTALL_DIR ?= ${HOME}
PROVIDER_INSTALL_DIR = $(INSTALL_DIR)/.terraform.d/plugins/example.com/tfp-example/example

ACC_TEST_VARS = TF_ACC=1

BINARY = $(shell go list)
EXT = $(shell go version | cut -d' ' -f4 | sed -E 's/\//_/')
VERSION ?= 0.0.1
DIST_BIN = $(BINARY)_v$(VERSION)
ZIP_FILE = $(BINARY)_$(VERSION)_$(EXT).zip

provider:
	go build

test:	provider
	go test ./example

acc: provider
	$(ACC_TEST_VARS) go test -v ./example

docs:
	go generate

dist: provider test
	mv $(BINARY) $(DIST_BIN)
	zip $(ZIP_FILE) $(DIST_BIN)

install: dist
	mkdir -p $(PROVIDER_INSTALL_DIR)
	rm -f $(PROVIDER_INSTALL_DIR)/$(ZIP_FILE)
	cp $(ZIP_FILE) $(PROVIDER_INSTALL_DIR)

.PHONY: clean

clean:
	go clean --testcache
	rm -rf $(BINARY)*
