IMAGE_NAME	?= takekazuomi/templo
TAG		?= 0.1.0
SRC		:= templo.go
GOLANG_CROSS	:= dockercore/golang-cross:1.13.15

.PHONY: build

help:	## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

all: build/templo-darwin-amd64 build/templo-linux-amd64 build/templo-windows-amd64.exe

build/templo-darwin-amd64: $(SRC)
	docker run -it --rm -u "$${UID}:$${GID}" -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=darwin $(GOLANG_CROSS) go build -o $@

build/templo-linux-amd64: $(SRC)
	docker run -it --rm -u "$${UID}:$${GID}" -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=linux $(GOLANG_CROSS) go build -o $@

build/templo-windows-amd64.exe: $(SRC)
	docker run -it --rm -u "$${UID}:$${GID}" -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=windows $(GOLANG_CROSS) go build -o $@

version:
	docker run -it --rm -v $(PWD):/src -w /src -e GOARCH=amd64 -e GOOS=windows $(GOLANG_CROSS) go version

build:	## build
	docker build \
		-t $(IMAGE_NAME):$(TAG) \
		-t $(IMAGE_NAME):latest \
		-f Dockerfile .

push:	## push
push:	build
	docker push $(IMAGE_NAME):$(TAG)
	docker push $(IMAGE_NAME):latest

login:	## login docker shell
	docker run -it --rm -u=$$(id -u):$$(id -g) -v $(PWD):/workspace $(IMAGE_NAME):latest /bin/zsh

lint:	## lint
	staticcheck ./...

test:	## test
	go test  ./...

fmt:	## format
	go fmt ./...

