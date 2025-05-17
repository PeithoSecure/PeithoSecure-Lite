# ==== Makefile (PeithoBackend: Stub Mode + Roast Mode + License Force + Screenshot) ====

GO := go

.PHONY: all tidy build run clean lint force-license roast screenshot

all: tidy build

## 🔍 Static Analysis
lint:
	@echo "🔍 Running lint checks (for what it’s worth)..."
	@$(GO) vet ./...
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "⚠️ goimports not found, skipping auto-format"

## 🧼 Clean
clean:
	@echo "🧼 Deleting evidence of ambition..."
	@rm -rf peitho-core peitho-server
	@rm -f ./peitho-core/unlock.lic ./peitho-core/keys/peitho_private.key ./peitho-core/keys/peitho_public.sig
	@echo "✨ Gone. Like your hopes of a prod deployment."

## 🔧 Dependency Management
tidy:
	$(GO) mod tidy

## 🏗️ Build
build:
	@echo "🔨 Compiling with stubbed dreams..."
	$(GO) build -o peitho-server ./cmd/peitho-server

## 🚀 Run (full launch)
run: build
	@echo "🎬 Running PeithoBackend in pure delusion mode..."
	./peitho-server

## 🥩 Roast-Only Mode
roast: build
	@echo "🥩 Simmering roast engine only..."
	@PEITHO_ROAST_ONLY=true ./peitho-server

## 🔐 Force regenerate unlock.lic using dev signer
force-license:
	@echo "🔐 Attempting to reforge unlock.lic — again. Because planning ahead is overrated..."
	go run peitho-core/scripts/signer/dev_signer.go \
		--email=dev@peitho.local \
		--device=web-default \
		--out=peitho-core/unlock.lic || (echo "💥 dev_signer.go exploded — no license for you." && exit 1)
	@echo "✅ unlock.lic forged under duress."
	@echo "🧠 Bound to: dev@peitho.local @ web-default"
	@echo "💀 You are now exactly one unsigned byte away from failure again."

## 📸 Capture roast output to log
screenshot: build
	@echo "📸 Capturing roast for posterity..."
	@PEITHO_ROAST_ONLY=true ./peitho-server > peitho-trace.log 2>&1 || true
	@echo "📝 Roast log saved to peitho-trace.log — frame it, weep later."

## 🐳 Docker Builds
docker-dev:
	@echo "🐳 Building dev image (still fake security)..."
	docker-compose build --no-cache

docker-release:
	@echo "🚀 Building release image. Good luck convincing anyone it’s secure..."
	docker-compose build --no-cache

up:
	docker-compose up

full-rebuild: clean docker-release up
