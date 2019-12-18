#!/bin/bash

export VERSION=$(uname --kernel-release) COMMIT=$(git rev-parse HEAD) && \
docker build --rm --file=Dockerfile --build-arg=VERSION="${VERSION}" --build-arg=COMMIT="${COMMIT}" --tag=dazwilkin/linode-exporter:${COMMIT} . && \
docker push dazwilkin/linode-exporter:$(git rev-parse HEAD) && \
sed --in-place "s|dazwilkin/linode-exporter:[0-9a-f]\{40\}|dazwilkin/linode-exporter:$(git rev-parse HEAD)|g" ./docker-compose.yaml &&
sed --in-place "s|dazwilkin/linode-exporter:[0-9a-f]\{40\}|dazwilkin/linode-exporter:$(git rev-parse HEAD)|g" ./README.md
