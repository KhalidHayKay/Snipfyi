kFROM golang:1.25-alpine
WORKDIR /app

ARG APP_ENV

RUN if [ "$APP_ENV" = "development" ]; then \
      go install github.com/air-verse/air@latest; \
    fi

COPY . .

RUN if [ "$APP_ENV" != "development" ]; then \
      go build -o main .; \
    fi

CMD if [ "$APP_ENV" = "development" ]; then \
      air; \
    else \
      ./main; \
    fi
