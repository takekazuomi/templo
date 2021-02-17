SRC := templo.go

GOLANG_CROSS := dockercore/golang-cross:1.13.15

all: build/templo-darwin-amd64 build/templo-linux-amd64 build/templo-windows-amd64.exe

build/templo-darwin-amd64: $(SRC)
	docker run -it --rm -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=darwin $(GOLANG_CROSS) go build -o $@

build/templo-linux-amd64: $(SRC)
	docker run -it --rm -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=linux $(GOLANG_CROSS) go build -o $@

build/templo-windows-amd64.exe: $(SRC)
	docker run -it --rm -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=windows $(GOLANG_CROSS) go build -o $@

version:
	docker run -it --rm -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=windows $(GOLANG_CROSS) go version

