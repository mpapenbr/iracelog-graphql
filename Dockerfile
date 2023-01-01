FROM alpine:3.17
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/iracelog-graphql"]
COPY iracelog-graphql /
COPY scripts/wait-for-it.sh /wait-for-it.sh
# COPY config.yml /
EXPOSE 8080