FROM alpine:latest

RUN apk update && \
    apk add ca-certificates tzdata && \
    update-ca-certificates && \
    apk add shadow && \
    groupadd -r zbi && \
    useradd -r -g zbi -s /sbin/nologin -c "zbi user" zbi

USER zbi
WORKDIR /zbi

COPY ./zbi-controller /usr/bin/zbi-controller

EXPOSE 8080

CMD ["zbi-controller", "--port", "8080"]
