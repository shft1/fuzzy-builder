FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/migrator ./cmd/migrator

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/migrator /app/migrator
ENTRYPOINT ["/app/migrator"]

