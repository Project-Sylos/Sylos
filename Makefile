APP_NAME := sylos
UI_DIR := ui
UI_DIST := internal/ui/dist
BIN := bin/$(APP_NAME)
DIST_DIR := dist
OAUTH_ENV := oauth.env
OAUTH_KEY := oauth.key
PARENT := $(abspath ..)

# Reveal enc: values from oauth.env using oauth.key for -ldflags bake.
OAUTH_CMD := go run ../Sylos-API/cmd/oauthenv
OAUTH_LDFLAGS := $(shell $(OAUTH_CMD) ldflags $(OAUTH_ENV) $(OAUTH_KEY) 2>/dev/null)

# Repos to test before build/release. Paths are relative to this Makefile's parent.
TEST_REPOS := Sylos Sylos-API Migration-Engine Sylos-FS Spectra go-path-linter

.PHONY: build run ui ui/dist release tidy fmt test test-all clean oauth-genkey oauth-encrypt

build: test-all ui/dist
	go build -ldflags "$(OAUTH_LDFLAGS)" -o $(BIN) ./cmd/sylos

run: build
	./$(BIN)

ui: ui/dist

ui/dist:
	@if [ -d "$(UI_DIR)" ]; then \
		cd $(UI_DIR) && npm ci --legacy-peer-deps && VITE_API_BASE= npm run build; \
		rm -rf $(UI_DIST); \
		mkdir -p $(UI_DIST); \
		cp -r $(UI_DIR)/dist/. $(UI_DIST)/; \
	elif [ -f "$(UI_DIST)/index.html" ]; then \
		echo "using existing $(UI_DIST)"; \
	else \
		echo "missing $(UI_DIR) submodule and $(UI_DIST); run: git submodule update --init --recursive"; \
		exit 1; \
	fi

release: build
	mkdir -p $(DIST_DIR)
	cp $(BIN) $(DIST_DIR)/
	cd $(DIST_DIR) && zip -q -r $(APP_NAME)-alpha.zip $(APP_NAME)
	@echo "wrote $(DIST_DIR)/$(APP_NAME)-alpha.zip"

# Create oauth.key (32-byte hex). Intended to be committed with oauth.env.
oauth-genkey:
	$(OAUTH_CMD) genkey $(OAUTH_KEY)

# Rewrite plaintext values in oauth.env as enc:… (skips already-encrypted lines).
oauth-encrypt:
	$(OAUTH_CMD) encrypt $(OAUTH_ENV) $(OAUTH_KEY)

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

# Run go test in this repo and each replace sibling that exists. Continues after
# failures so the log shows every broken suite, then exits non-zero with a summary.
test-all:
	@failed=""; \
	for name in $(TEST_REPOS); do \
		if [ "$$name" = "Sylos" ]; then \
			dir="$(CURDIR)"; \
		else \
			dir="$(PARENT)/$$name"; \
		fi; \
		if [ ! -d "$$dir" ]; then \
			echo "=== Skipping $$name (missing $$dir) ==="; \
			continue; \
		fi; \
		if [ ! -f "$$dir/go.mod" ]; then \
			echo "=== Skipping $$name (no go.mod in $$dir) ==="; \
			continue; \
		fi; \
		echo "=== Testing $$name ==="; \
		if ! (cd "$$dir" && go test ./... -count=1); then \
			echo "=== FAILED $$name ==="; \
			failed="$$failed $$name"; \
		else \
			echo "=== OK $$name ==="; \
		fi; \
	done; \
	if [ -n "$$failed" ]; then \
		echo ""; \
		echo "test-all failed for:$$failed"; \
		echo "See FAIL output and stack traces above."; \
		exit 1; \
	fi; \
	echo ""; \
	echo "test-all: all repos passed."

clean:
	rm -rf $(BIN) $(UI_DIST) $(DIST_DIR)
