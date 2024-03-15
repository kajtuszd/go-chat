package main

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"strconv"
)

type Room struct {
	ID        string
	Join      chan *User
	Exit      chan *User
	Broadcast chan string
	Clients   map[string]*User
}

type room interface {
	addUser(user *User)
	removeUser(user *User)
	broadcastMessage(message string)
	handleMessages()
	readMessages(user *User)
}

func NewRoom() *Room {
	return &Room{
		ID:        strconv.Itoa(rand.Int() % 1000),
		Join:      make(chan *User),
		Exit:      make(chan *User),
		Broadcast: make(chan string),
		Clients:   make(map[string]*User),
	}
}

func (r *Room) addUser(user *User) {
	r.Clients[user.ID] = user
}

func (r *Room) removeUser(user *User) {
	delete(r.Clients, user.ID)
	user.Conn.Close()
}

func (r *Room) broadcastMessage(message string) {
	for _, user := range r.Clients {
		if err := user.Conn.WriteMessage(websocket.TextMessage, []byte("["+r.ID+"]~ "+message)); err != nil {
			return
		}
	}
}

func (r *Room) handleMessages() {
	for {
		select {
		case u := <-r.Join:
			r.addUser(u)
		case u := <-r.Exit:
			r.removeUser(u)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		}
	}
}

func (r *Room) readMessages(user *User) {
	defer func() {
		if user.Conn != nil {
			r.Exit <- user
			log.Println("Connection closed for user:", user.Username)
		}
		return
	}()
	for {
		_, p, err := user.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}
		r.Broadcast <- string(p)
	}
}
