FROM golang:1.22 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN useradd -u 10001 snex

WORKDIR /go/src/snex

# Update dependencies: On unchanged dependencies, cached layer will be reused
COPY go.* ./
RUN go mod tidy

# Build
COPY . ./

RUN go build -o snex

# Pack
FROM gcr.io/distroless/static-debian11

COPY --from=build /go/src/snex/snex .
USER snex
LABEL maintainer="Daan Gerits <daan@shono.io>"

ENTRYPOINT ["/snex"]