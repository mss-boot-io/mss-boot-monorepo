FROM alpine:latest
WORKDIR /app
ADD application application
COPY ../../bin/docker-start.sh .

ENV GIN_MODE release

CMD ["sh", "./docker-start.sh"]
