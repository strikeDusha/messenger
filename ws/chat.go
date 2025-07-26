package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"messenger/database"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan database.Message),
	}
}

func (h *Hub) Run() { //функция обслуживающая наши потоки в main go

	log.Println("run hub")
	for {
		select { // читаем потоки
		case client := <-h.Register: // если рега  - добавляем в клиенты
			h.Mu.Lock()
			h.Clients[client.Username] = client // мапу только с мутехом
			fmt.Println("added user to our hub")
			h.Mu.Unlock()
		case client := <-h.Unregister: // анрега  - удаляем
			if _, ok := h.Clients[client.Username]; ok {
				h.Mu.Lock()
				delete(h.Clients, client.Username) // удаление из мапы тоже только с мутексом
				close(client.Send)
				h.Mu.Unlock()
			}

		case message := <-h.Broadcast: // если пришло сообщение
			log.Printf("Hub got message %+v\n", message)
			if reciever, ok := h.Clients[message.Receiver]; ok { // достаем клиента из мапы
				b, err := json.Marshal(message)
				log.Println("пришло сообщение !" + string(b)) // маршалим месаж, пишем слайс байтов в б
				if err == nil {
					reciever.Send <- b // передаем его в канал приемки ресивера
				} else {
					log.Printf("user offline%+v", message)
				}
			}
		}
	}
}

func (h *Hub) ReadPump(c *Client) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("readPump ReadMessage error:", err)
			break
		}
		log.Printf("readPump raw message: %s\n", string(msgBytes))

		var m database.Message
		if err := json.Unmarshal(msgBytes, &m); err != nil {
			log.Println("readPump Unmarshal error:", err)
			continue
		}
		log.Printf("readPump unmarshaled: %+v\n", m)

		m.Sender = c.Username
		m.Sent = time.Now()
		log.Printf("readPump sending to broadcast: %+v\n", m)
		h.Broadcast <- m
	}
}
func (h *Hub) WritePump(c *Client) {
	log.Println("w pump")
	for msg := range c.Send { // а тут не ебу вообще че происходит и движа с жсоном нету
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
