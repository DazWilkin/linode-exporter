ARG GOLANG_VERSION=1.13

FROM golang:${GOLANG_VERSION} as build

ARG VERSION=""
ARG COMMIT=""

WORKDIR /linode-exporter

# TODO(dazwilkin) Local go.mod includes: replace github.com/linode/linodego => .../linodego
RUN echo "module github.com/DazWilkin/linode-exporter\ngo 1.13\nrequire (\n)\n" > ./go.mod

COPY main.go .
COPY collector ./collector
# TODO(dazwilkin) remove this
COPY mock ./mock

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/linode-exporter \
    ./main.go

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/bin/linode-exporter /

ENTRYPOINT ["/linode-exporter"]
