# Step 1: Modules caching
FROM golang:1.25.5-alpine3.23 AS modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Step 2: Tests
FROM golang:1.25.5-alpine3.23

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN go env -w CGO_ENABLED=0 \
    && go env -w GOOS=linux

CMD ["go", "test", "-v", "./integration-test/..."]
