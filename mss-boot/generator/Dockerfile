FROM alpine:latest
WORKDIR /app
ADD admin application
COPY ../../bin/docker-start.sh .
COPY cfg/*.yaml ./cfg/

ENV GIN_MODE release

CMD ["sh", "./docker-start.sh"]
