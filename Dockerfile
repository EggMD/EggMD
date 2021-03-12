FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

RUN mkdir /etc/EggMD
WORKDIR /etc/EggMD

ADD EggMD /etc/EggMD

RUN chmod 655 /etc/EggMD/EggMD
ENV MACARON_ENV production

ENTRYPOINT ["/etc/EggMD/EggMD", "web"]
EXPOSE 1999