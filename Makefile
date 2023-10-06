
install:
	go mod tidy

build: install
	go build -o output/dashboard

run: build
	./output/dashboard

