APP_NAME := sylos
UI_DIR := ui
UI_DIST := internal/ui/dist
BIN := bin/$(APP_NAME)

.PHONY: build run ui ui/dist tidy fmt test clean

build: ui/dist
	go build -o $(BIN) ./cmd/sylos

run: build
	./$(BIN)

ui: ui/dist

ui/dist:
	@test -d $(UI_DIR) || (echo "missing $(UI_DIR) submodule; run: git submodule update --init --recursive" && exit 1)
	cd $(UI_DIR) && npm ci --legacy-peer-deps && OAUTH_CREDS_DIR=$(CURDIR)/creds VITE_API_BASE= npm run build
	rm -rf $(UI_DIST)
	mkdir -p $(UI_DIST)
	cp -r $(UI_DIR)/dist/. $(UI_DIST)/

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf $(BIN) $(UI_DIST)
