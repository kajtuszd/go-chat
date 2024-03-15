package main

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func contains(slice []int, element int) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func generateRandomID() string {
	var generatedIDs = make([]int, 20)
	index := 0
	return func() string {
		ID := rand.Int()
		for contains(generatedIDs, ID) {
			ID = rand.Int()
		}
		generatedIDs[index] = ID
		index++
		return strconv.Itoa(rand.Int())
	}()
}

func handler(s Server, w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	username := queryParams.Get("username")
	roomID := queryParams.Get("room")
	if username == "" {
		http.Error(w, "Required username", http.StatusBadRequest)
		return
	}
	user := &User{
		ID:       generateRandomID(),
		Username: username,
		Conn:     nil,
	}
	if roomID == "" {
		roomID = s.addRoom(NewRoom())
	}
	room, err := s.getRoomByID(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	user.Conn = conn
	go room.handleMessages()
	go room.readMessages(user)
	room.Join <- user
	log.Println("Connection open for: ", user.Username)
}

func main() {
	server := NewServer()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler(*server, w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("./src")))
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
