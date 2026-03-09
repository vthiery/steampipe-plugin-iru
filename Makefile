STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/local/iru/iru.plugin -tags "${BUILD_TAGS}" .

test: install
	@bash scripts/test_tables.sh

.PHONY: install test
