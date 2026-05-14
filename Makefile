BINARY   := anuncios
TMPDIR   := /tmp/go_build_tmp
GOCACHE  := /tmp/go_cache
GO_ENV   := TMPDIR=$(TMPDIR) GOCACHE=$(GOCACHE)

.PHONY: all build run run-unlimited run-tracks clean help

all: build

## build: Compila el binario
build:
	@mkdir -p $(TMPDIR) $(GOCACHE)
	$(GO_ENV) go build -o $(BINARY) .
	@echo "✅ Build exitoso → ./$(BINARY)"

## run: Ejecuta el binario (modo normal)
run: build
	./$(BINARY)

## run-unlimited: Ejecuta sin límite de jobs (-i)
run-unlimited: build
	./$(BINARY) -i

## run-tracks: Ejecuta con matching de tracks (-t)
run-tracks: build
	./$(BINARY) -t

## clean: Elimina el binario compilado
clean:
	@rm -f $(BINARY)
	@echo "🗑️  Binario eliminado"

## help: Muestra esta ayuda
help:
	@echo ""
	@echo "Uso: make [target]"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## /  /'
	@echo ""
