ARG GOLANG_VERSION=1.23

ARG COMMIT
ARG VERSION

ARG TARGETOS
ARG TARGETARCH

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /linode-exporter

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go .
COPY collector ./collector

ARG TARGETOS
ARG TARGETARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /bin/exporter \
    ./main.go


FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.source=https://github.com/DazWilkin/linode-exporter

COPY --from=build /bin/exporter /

ENTRYPOINT ["/exporter"]
