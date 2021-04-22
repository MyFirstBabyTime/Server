FROM alpine
MAINTAINER Park, Jinhong <jinhong0719@naver.com>

RUN apk add curl

COPY ./first-baby-time ./first-baby-time
ENTRYPOINT [ "/first-baby-time" ]
