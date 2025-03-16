FROM alpine:3.14
ARG ARCH
ENTRYPOINT ["/iracelog-graphql"]
HEALTHCHECK --interval=2s --timeout=2s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/healthz || exit 1

COPY samples /
COPY iracelog-graphql /


EXPOSE 8080