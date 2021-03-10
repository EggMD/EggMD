FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

RUN mkdir /etc/EggMD
WORKDIR /etc/EggMD

ADD Elaina /etc/EggMD

RUN chmod 655 /etc/EggMD/EggMD

ENTRYPOINT ["/etc/EggMD/EggMD"]
EXPOSE 1999