FROM google/cloud-sdk:alpine

RUN apk update && apk add jq && rm -rf /var/cache/apk/*
ADD mantle /go/bin/mantle
RUN addgroup -S mantle && adduser -S mantle -G mantle

USER mantle