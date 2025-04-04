package controller

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Game

	egress chan []byte
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Game) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (self *Client) readMessage() {
	defer func() {
		self.manager.removeClient(self)
	}()

	for {
		messageType, payload, err := self.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		for wsclient := range self.manager.clients {
			wsclient.egress <- payload
		}

		log.Println(messageType)
		log.Println(string(payload))
	}
}

func (self *Client) writeMessage() {
	defer func() {
		self.manager.removeClient(self)
	}()

	for {
		select {
		case message, ok := <-self.egress:
			if !ok {
				if err := self.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("Connection closed: %v", err)
				}

				return
			}

			if err := self.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failled to send message: %v", err)
			}
			log.Println("Message sent")
		}
	}
}
