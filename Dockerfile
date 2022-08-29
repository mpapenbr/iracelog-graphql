FROM alpine:3.15
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/iracelog-graphql"]
COPY iracelog-graphql /
# COPY config.yml /
EXPOSE 8080