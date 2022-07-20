FROM alpine:latest
WORKDIR /app
ADD store application
COPY bin/docker-start.sh .

CMD ["sh", "./docker-start.sh"]
