FROM google/cloud-sdk:alpine

RUN apk update && apk add jq && rm -rf /var/cache/apk/*
ADD mantle /usr/local/bin/mantle
RUN addgroup -S mantle && adduser -S mantle -G mantle

USER mantle
RUN chmod u+x /usr/local/bin/mantle
