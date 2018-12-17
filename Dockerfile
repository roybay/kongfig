FROM golang:1.11-stretch

WORKDIR /go/src/github.com/pagerinc/kongfig/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o kongfig

FROM scratch

ENV PORT=3000
COPY --from=0 /go/src/github.com/pagerinc/kongfig/kongfig /go/kongfig
EXPOSE 3000

ENTRYPOINT ["/go/kongfig"]
