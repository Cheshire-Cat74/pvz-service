FROM golang:1.23.0

WORKDIR ${GOPATH}/pvz-service/
COPY . ${GOPATH}/pvz-service/

RUN go mod download

RUN go test -cover ./internal/handler/... ./internal/service/... ./internal/repository/...


RUN go build -o /build ./cmd \
    && go clean -cache -modcache


EXPOSE 8080 9000 3000

CMD ["/build"]
