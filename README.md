# Websocket Proxy Example

## Requirements

go get golang.org/x/net/websocket

## Endpoints
    
* **:9090**: php webserver, public
* **:9191**: websocket, public
* **:9292**: udp server, private only reachable from localhost

## Dockerized Example

### Setup
This builds two docker images.

* websocket: Contains go and the `server.go` script which is executed at container start. 
* websocket-php: A PHP 7 docker image, the container will run the built-in webserver


    ./bin/build
    
### Usage

Start environment:

    ./bin/start
    
Call the demo page:

    http://<your development hostname>:9090/index.php
    
This opens a websocket and listens to it.
    
Send messages through the web socket:

    ./bin/send-messages
    
The latest script starts an endless loop! You have to stop the command by hand though!

For the lazy dev there is a shortcut for all those commands:

    ./bin/start-dev
    
This builds the Docker Images, starts the environment and sends messages in endless loop!

