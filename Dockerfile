ARG GO_VERSION=1.16.5

FROM golang:${GO_VERSION}-buster AS build_base
WORKDIR /build
COPY . /build
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -x -installsuffix cgo -o ws-examples .

# Build the Go app
RUN go build -o ./out/ws-examples .

FROM debian:buster-slim
RUN set -x && apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build_base /build/ws-examples .
EXPOSE 8081
ENTRYPOINT ["./ws-examples"]

