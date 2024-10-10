.DEFAULT_GOAL := goapp

.PHONY: all
all: clean goapp

.PHONY: goapp
goapp:
	mkdir -p bin
	go build -o bin/goapp/cmd/server ./main.go

.PHONY: clean
clean:
	go clean
	rm -f bin/*