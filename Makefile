GOARCH?=amd64
GOOS?=linux

all: test kongfig

kongfig: main.go
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o kongfig

clean:
	@rm -f kongfig *.out

ci:
	@ go test -covermode=atomic -coverprofile=coverage.out -race

test: clean
	@ go test -covermode=count -coverprofile=coverage.out

cover: test
	@ go tool cover -html=coverage.out

count: test
	@ go tool cover -func=coverage.out
