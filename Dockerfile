ARG GOLANG_VERSION=1.21.4-alpine

# Dev Stage
FROM golang:${GOLANG_VERSION} AS dev
WORKDIR /app
CMD ["sh"]

# Build Stage
FROM golang:${GOLANG_VERSION} AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app

# Release Stage
FROM alpine:latest AS release
WORKDIR /app
COPY --from=build /app ./
CMD ["./app"]
