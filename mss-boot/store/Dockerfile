FROM alpine:latest
WORKDIR /app
ADD store application
COPY ../../bin/docker-start.sh .
COPY cfg/*.yaml ./cfg/

CMD ["sh", "./docker-start.sh"]
