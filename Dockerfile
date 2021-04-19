FROM alpine
MAINTAINER Park, Jinhong <jinhong0719@naver.com>

COPY ./first-baby-time ./first-baby-time
ENTRYPOINT [ "/first-baby-time" ]
