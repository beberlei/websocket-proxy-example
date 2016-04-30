package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"time"
)

var maxId = 0
var udpAddr string = ":9292"
var httpAddr string = ":9191"
var clients            map[int]*Client
var clientAdded        chan *Client
var clientDisconnected chan *Client
var eventReceived      chan *Event

type Event struct {
	UserId  int
	Name    string
	Payload interface{}
}

func (self *Event) String() string {
	return self.Name
}

type Client struct {
	id        int
	ws        *websocket.Conn
	sendEvent chan *Event
}

func (c *Client) Start() {
	log.Println("Listening to events")
	clientAdded <- c
	for {
		select {
		case event := <-c.sendEvent:
			log.Printf("Send to %d: %s\n", c.id, event)
			err := websocket.JSON.Send(c.ws, event)

			if err != nil {
				clientDisconnected <- c
				return
			}
		}
	}
}

func HandleWsRequest(ws *websocket.Conn) {
	defer ws.Close()
	maxId++ // this is not thread-safe

	client := &Client{
		id: maxId,
		ws: ws,
		sendEvent: make(chan *Event, 10),
	}

	client.Start()
}

func ListenHttp() {
	log.Println("Listening server...")

	http.Handle("/ws", websocket.Handler(HandleWsRequest))
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(httpAddr, nil)
}

func EventLoop() {
	timer := time.NewTicker(time.Duration(60) * time.Second).C

	for {
		select {

		case <-timer:
			log.Printf("Max clients %d\n", maxId)
			log.Printf("Sending pings to %d to check for timeouts.\n", len(clients))
			if len(clients) > 0 {
				for id, client := range clients {
					select {
					case client.sendEvent <- &Event{UserId: id, Name: "Ping"}:
					default:
					}
				}
			}

		case event := <-eventReceived:
			if client, ok := clients[event.UserId]; ok {
				select {
				case client.sendEvent <- event:
				default:
				}
			}

		case c := <-clientAdded:
			clients[c.id] = c
			log.Printf("Added new client %d, now %d clients connected.\n", c.id, len(clients))

		case c := <-clientDisconnected:
			log.Println("Delete client")
			delete(clients, c.id)
		}
	}
}

func ListenUdp() {
	serverAddr, _ := net.ResolveUDPAddr("udp", udpAddr)
	conn, _ := net.ListenUDP("udp", serverAddr)

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		log.Printf("Received UDP msg %s\n", buf[0:n])

		if err != nil {
			log.Printf("err")
			continue
		}

		var event *Event
		err = json.Unmarshal(buf[0:n], &event)

		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}

		eventReceived <- event
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	clients =            make(map[int]*Client)
	clientAdded =        make(chan *Client)
	clientDisconnected = make(chan *Client)
	eventReceived =     make(chan *Event, 100)

	go ListenHttp()
	go ListenUdp()

	EventLoop()
}
