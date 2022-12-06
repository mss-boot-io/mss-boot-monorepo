FROM alpine:latest
WORKDIR /app
ADD application application
COPY ../../../bin/docker-start.sh .

CMD ["sh", "./docker-start.sh"]
