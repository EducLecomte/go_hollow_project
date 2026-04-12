.PHONY: build clean install

BINARY_NAME=hollow

build:
	@echo "Compilation de Hollow..."
	go build -o $(BINARY_NAME) ./cmd/hollow

clean:
	@echo "Nettoyage..."
	rm -f $(BINARY_NAME)

install: build
	@echo "Installation dans /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/