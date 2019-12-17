FROM golang:1.13-alpine as builder
WORKDIR /go/src/github.com/oliver006/redis_exporter/

ADD *.go /go/src/github.com/oliver006/redis_exporter/
ADD vendor /go/src/github.com/oliver006/redis_exporter/vendor

ARG GOARCH="amd64"
ARG SHA1="[no-sha]"
ARG TAG="[no-tag]"

RUN apk --no-cache add ca-certificates
RUN BUILD_DATE=$(date +%F-%T) && CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o /redis_exporter \
    -ldflags  "-s -w -extldflags \"-static\" -X main.BuildVersion=$TAG -X main.BuildCommitSha=$SHA1 -X main.BuildDate=$BUILD_DATE" .

RUN [ $GOARCH = "amd64" ]  && /redis_exporter -version || ls -la /redis_exporter

FROM alpine:latest

COPY --from=builder /redis_exporter /redis_exporter
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

EXPOSE 9121


ENV TILE38_VERSION 1.19.3
ENV TILE38_DOWNLOAD_URL https://github.com/tidwall/tile38/releases/download/$TILE38_VERSION/tile38-$TILE38_VERSION-linux-amd64.tar.gz

RUN addgroup -S tile38 && adduser -S -G tile38 tile38

RUN apk update \
    && apk add ca-certificates \
    && update-ca-certificates \
    && apk add openssl \
    && wget -O tile38.tar.gz "$TILE38_DOWNLOAD_URL" \
    && tar -xzvf tile38.tar.gz \
    && rm -f tile38.tar.gz \
    && mv tile38-$TILE38_VERSION-linux-amd64/tile38-server /usr/local/bin \
    && mv tile38-$TILE38_VERSION-linux-amd64/tile38-cli /usr/local/bin \
    && mv /redis_exporter /usr/local/bin \
    && rm -fR tile38-$TILE38_VERSION-linux-amd64

RUN mkdir /data && chown tile38:tile38 /data

RUN apk add --no-cache --upgrade bash

COPY exec.sh /exec.sh
RUN chmod 777 -R /exec.sh

VOLUME /data
WORKDIR /data

EXPOSE 9851
CMD ["/exec.sh"]