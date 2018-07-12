ARCH?=amd64

all: test kongfig

kongfig: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o kongfig

server: server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o server

clean:
	@rm -f server *.out

ci:
	@ go test -covermode=atomic -coverprofile=coverage.out -race

test: clean
	@ go test -covermode=count -coverprofile=coverage.out

cover: test
	@ go tool cover -html=coverage.out

count: test
	@ go tool cover -func=coverage.out
