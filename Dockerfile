FROM golang:1.16-alpine AS build

WORKDIR /app
ADD . /app
RUN cd /app && go build -o challenge
ENTRYPOINT ./challenge