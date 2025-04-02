.SILENT:
default: run
bin_name = quiz

.PHONY: clear
clear:
	clear

.PHONY: build
build:
	go build -o ./bin/${bin_name} ./main.go

.PHONY: run
run: clear build
	./bin/${bin_name}