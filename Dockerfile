# This file is used by goreleaser

ARG BUILDPLATFORM
FROM --platform=$BUILDPLATFORM alpine:3.22
# TARGETPLATFORM needs to be set after FROM
ARG TARGETPLATFORM

ENTRYPOINT ["/iracelog-graphql"]
HEALTHCHECK --interval=2s --timeout=2s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/healthz || exit 1

COPY samples /
COPY $TARGETPLATFORM/iracelog-graphql /


EXPOSE 8080