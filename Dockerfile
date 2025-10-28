# syntax=docker/dockerfile:1.7
FROM golang:1.24.5-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
# cache модулей Go
RUN --mount=type=cache,target=/go/pkg/mod go mod download
# копируем только backend, без frontend (уменьшаем build context)
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations
# подтянем новые зависимости (например modernc.org/sqlite) и синхронизируем go.mod/go.sum
RUN --mount=type=cache,target=/go/pkg/mod go mod tidy
# кэш компиляции
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
ENV HTTP_ADDR=:8080
ENV UPLOAD_DIR=/data/uploads
COPY --from=build /out/api /app/api
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/app/api"]
