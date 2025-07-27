package handlers

import (
	"messenger/database"
	"sync"

	"github.com/gorilla/websocket"
)

type DB struct { // database for all logic
	S *database.Storage // storage
	H *Hub              //websocket hub
}

func NewDB(s *database.Storage, h *Hub) *DB {
	return &DB{S: s, H: h}
}

type Client struct { // client
	Conn     *websocket.Conn
	Username string
	Send     chan []byte
}

type Hub struct { // хаб, всем управляет
	Mu         sync.Mutex            // что бы мапу с горутинами подружить
	Clients    map[string]*Client    // мапа клиентов по юзернейму
	Register   chan *Client          // поток для регистрации клиентов на вебсокет
	Unregister chan *Client          // поток для того что бы удалить
	Broadcast  chan database.Message // канал для сообщений
}

func NewClient(c *websocket.Conn, us string, s chan []byte) *Client {
	return &Client{c, us, s}
}
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan database.Message),
	}
}
