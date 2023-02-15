ARG GOLANG_VERSION=1.20
ARG GOLANG_OPTIONS="CGO_ENABLED=0 GOOS=linux GOARCH=amd64"

FROM docker.io/golang:${GOLANG_VERSION} as build

WORKDIR /linode-exporter

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go .
COPY collector ./collector

ARG VERSION=""
ARG COMMIT=""

RUN env ${GOLANG_OPTIONS} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /bin/exporter \
    ./main.go

FROM gcr.io/distroless/base-debian11

LABEL org.opencontainers.image.source https://github.com/DazWilkin/linode-exporter

COPY --from=build /bin/linode-exporter /

ENTRYPOINT ["/exporter"]
