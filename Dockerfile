FROM scratch
ARG ARCH
ENTRYPOINT ["/iracelog-graphql"]
HEALTHCHECK --interval=2s --timeout=2s --start-period=5s --retries=3 CMD [ "/grpc_health_probe", "-addr", "localhost:8080" ]
COPY iracelog-graphql /
COPY ext/healthcheck/grpc_health_probe.$ARCH /grpc_health_probe

EXPOSE 8080