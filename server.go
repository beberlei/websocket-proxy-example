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
	server    *Server
	sendEvent chan *Event
}

func NewClient(ws *websocket.Conn, s *Server) *Client {
	maxId++ // this is not thread-safe
	return &Client{
		ws:        ws,
		server:    s,
		id:        maxId,
		sendEvent: make(chan *Event, 10),
	}
}

func (c *Client) ListenEvents() {
	log.Println("Listening to events")
	for {
		select {
		case event := <-c.sendEvent:
			log.Printf("Send to %d: %s\n", c.id, event)
			err := websocket.JSON.Send(c.ws, event)

			if err != nil {
				c.server.clientDisconnected <- c
			}
		}
	}
}

// Chat server.
type Server struct {
	clients            map[int]*Client
	clientAdded        chan *Client
	clientDisconnected chan *Client
	eventReceived      chan *Event
}

func NewServer() *Server {
	return &Server{
		clients:            make(map[int]*Client),
		clientAdded:        make(chan *Client),
		clientDisconnected: make(chan *Client),
		eventReceived:      make(chan *Event, 100),
	}
}

func (s *Server) ListenHttp() {

	log.Println("Listening server...")

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		client := NewClient(ws, s)
		s.clientAdded <- client
		client.ListenEvents() // blocks
	}))
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}

func (s *Server) EventLoop() {
	timer := time.NewTicker(time.Duration(60) * time.Second).C

	for {
		select {

		case <-timer:
			log.Printf("Max clients %d\n", maxId)
			log.Printf("Sending pings to %d to check for timeouts.\n", len(s.clients))
			if len(s.clients) > 0 {
				for id, client := range s.clients {
					select {
					case client.sendEvent <- &Event{UserId: id, Name: "Ping"}:
					default:
					}
				}
			}

		case event := <-s.eventReceived:
			if client, ok := s.clients[event.UserId]; ok {
				select {
				case client.sendEvent <- event:
				default:
				}
			}

		case c := <-s.clientAdded:
			s.clients[c.id] = c
			log.Printf("Added new client %d, now %d clients connected.\n", c.id, len(s.clients))

		case c := <-s.clientDisconnected:
			log.Println("Delete client")
			delete(s.clients, c.id)
		}
	}
}

func (s *Server) ListenUdp() {
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8081")
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

		s.eventReceived <- event
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	server := NewServer()
	go server.ListenHttp()
	go server.ListenUdp()

	server.EventLoop()
}
