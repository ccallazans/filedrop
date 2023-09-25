FROM golang:1.21-alpine

WORKDIR /app

ADD ./bin/app .
ADD ./.env .
ADD ./migrations ./

EXPOSE 3000
ENTRYPOINT ./app