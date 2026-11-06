GO ?= go

.PHONY: test vet check gallery build snapshots snapshots-all smoke-native

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

check: test vet

gallery:
	$(GO) run ./cmd/gallery

build:
	$(GO) build -o bin/gallery ./cmd/gallery

snapshots:
	$(GO) run ./cmd/gallery -snapshot artifacts/gallery-light.png -theme light
	$(GO) run ./cmd/gallery -snapshot artifacts/gallery-dark.png -theme dark
	$(GO) run ./cmd/gallery -snapshot artifacts/gallery-forms.png -theme light -page forms
	$(GO) run ./cmd/gallery -snapshot artifacts/gallery-narrow.png -theme light -width 420 -height 780

smoke-native:
	$(GO) build -tags e2e -o bin/gallery-e2e ./cmd/gallery
	python3 scripts/smoke_native.py

snapshots-all:
	$(GO) run ./cmd/gallery -snapshot-suite artifacts/gallery-suite
