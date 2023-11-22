ARG GOLANG_VERSION=1.21.4-alpine

# Dev Stage
FROM golang:${GOLANG_VERSION} AS dev
WORKDIR /app
CMD ["sh"]

# Build Stage
FROM golang:${GOLANG_VERSION} AS build
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN go mod download
RUN mkdir /app && CGO_ENABLED=0 GOOS=linux go build -o /app/app cmd/main.go

# Release Stage
FROM alpine:latest AS release
WORKDIR /app
COPY --from=build /app ./
CMD ["./app"]
