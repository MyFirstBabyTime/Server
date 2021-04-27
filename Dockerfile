FROM alpine
MAINTAINER Park, Jinhong <jinhong0719@naver.com>

RUN apk add curl
RUN apk add docker

COPY ./first-baby-time ./first-baby-time
ENTRYPOINT [ "/first-baby-time" ]
