package controller

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Game
	gameID     uint

	egress chan []byte
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Game, gameID uint) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		gameID:     gameID,
		egress:     make(chan []byte),
	}
}

func (self *Client) readMessage() {
	defer func() {
		self.manager.removeClient(self)
	}()

	for {
		_, payload, err := self.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		self.manager.playAMove(self.gameID, string(payload))

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
		}
	}
}
