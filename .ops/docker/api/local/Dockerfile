FROM golang:1.21-alpine
ENV CGO_ENABLED 1

RUN set -ex \
    && apk update && apk add \
      build-base \
      bash \
      inotify-tools \
      git \
      sqlite \
    && git clone https://github.com/go-delve/delve \
    && cd delve \
    && go install github.com/go-delve/delve/cmd/dlv \
    && git clone https://github.com/go-swagger/go-swagger go-swagger \
    && cd go-swagger \
    && go install ./cmd/swagger

WORKDIR /go/src/app
ADD ./goapi .

COPY ./.ops/docker/api/local/entrypoint.sh /
RUN chmod 777 /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]
