FROM golang:1.24.5-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# ensure migrations are present in the build context
RUN test -d migrations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/migrator ./cmd/migrator

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/migrator /app/migrator
COPY --from=build /app/migrations /app/migrations
ENTRYPOINT ["/app/migrator"]

