FROM google/cloud-sdk:alpine

RUN apk update && apk add jq && rm -rf /var/cache/apk/*
RUN addgroup -S mantle && adduser -S mantle -G mantle
RUN curl -s -H "Accept: application/octet-stream" -L \
      $(curl -s  https://api.github.com/repos/ovotech/mantle/releases/latest | \
      jq '.assets[]  | select(.browser_download_url | contains("linux")) | .url' | tr -d '"')  \
    -o /usr/local/bin/mantle && \
    chown mantle /usr/local/bin/mantle

USER mantle
RUN chmod u+x /usr/local/bin/mantle
