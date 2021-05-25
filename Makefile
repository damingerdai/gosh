# build gosh
clean:
	rm -r -f gosh

build: clean
	go build -o gosh *.go

run: build
	./gosh
