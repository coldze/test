FROM golang:1.13.4

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

HEALTHCHECK --start-period=5s --interval=10s --timeout=3s CMD curl -f http://localhost/ping || exit 1

ENTRYPOINT ["test"]