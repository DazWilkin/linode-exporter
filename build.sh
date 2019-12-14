#!/bin/bash

docker build --rm --file=Dockerfile --tag=dazwilkin/linode-exporter:$(git rev-parse HEAD) . && \
docker push dazwilkin/linode-exporter:$(git rev-parse HEAD) && \
sed --in-place "s|dazwilkin/linode-exporter:[0-9a-f]\{40\}|dazwilkin/linode-exporter:$(git rev-parse HEAD)|g" ./docker-compose.yaml &&
sed --in-place "s|dazwilkin/linode-exporter:[0-9a-f]\{40\}|dazwilkin/linode-exporter:$(git rev-parse HEAD)|g" ./README.md
