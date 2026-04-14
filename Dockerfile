FROM golang:1.25-alpine

WORKDIR /app

ARG APP_ENV

RUN if [ "$APP_ENV" = "development" ]; then \
    go install github.com/air-verse/air@latest; \
    fi

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN if [ "$APP_ENV" != "development" ]; then \
    go build -o cmd/web/main ./cmd/web && \
    go build -o cmd/worker/main ./cmd/worker; \
    fi

ENTRYPOINT ["/bin/sh", "-c", "exec sh start.sh"]