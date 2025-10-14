# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
ENV HTTP_ADDR=:8080
ENV UPLOAD_DIR=/data/uploads
COPY --from=build /out/api /app/api
VOLUME ["/data"]
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/api"]
