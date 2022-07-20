FROM alpine:latest
WORKDIR /app
ADD generator application
COPY bin/docker-start.sh .

ENV GIN_MODE release

CMD ["sh", "./docker-start.sh"]
