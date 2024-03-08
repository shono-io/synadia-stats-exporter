FROM golang:1.21 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN useradd -u 10001 snex

WORKDIR /go/src/snex
# Update dependencies: On unchanged dependencies, cached layer will be reused
COPY go.* /go/src/snex
RUN go mod download

# Build
COPY main.go /go/src/snex/
COPY cmd /go/src/snex/
COPY pkg /go/src/snex/

RUN go build -o snex

# Pack
FROM gcr.io/distroless/static-debian11

COPY --from=build /go/src/snex/snex .
USER snex:snex
LABEL maintainer="Daan Gerits <daan@shono.io>"

ENTRYPOINT ["/snex"]