FROM golang:alpine

RUN apk update && apk add curl git && rm -rf /var/cache/apk/*

ADD . /code

RUN go get golang.org/x/net/websocket

EXPOSE 9191

CMD ["go", "run", "/code/server.go"]
