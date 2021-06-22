ARG GO_VERSION=1.15

FROM golang:${GO_VERSION}-alpine AS build_base
RUN apk add --no-cache git
WORKDIR /build
COPY . /build
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -x -installsuffix cgo -o ws-examples .

# Build the Go app
RUN go build -o ./out/ws-examples .

FROM golang:${GO_VERSION}-alpine
RUN apk add ca-certificates
WORKDIR /app
COPY --from=build_base /build/ws-examples .
EXPOSE 7000
ENTRYPOINT ["./ws-examples"]

