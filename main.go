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

func handleUserCommunication(user *User) {
	defer func() {
		if user.Conn != nil {
			user.Conn.Close()
			delete(activeUsers, user.ID)
			log.Println("Connection closed for user:", user.Username)
		}
	}()
	for {
		messageType, p, err := user.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}
		for _, receiver := range activeUsers {
			log.Println(user.Username, ": ", string(p))
			if err := receiver.Conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	username := queryParams.Get("username")
	if username == "" {
		http.Error(w, "Required username", http.StatusBadRequest)
		return
	}
	user := &User{
		ID:       generateRandomID(),
		Username: username,
		Conn:     nil,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	user.Conn = conn
	log.Println("Connection open for: ", user.Username)
	activeUsers[user.ID] = user
	go handleUserCommunication(user)
}

func main() {
	http.HandleFunc("/ws", handler)
	http.Handle("/", http.FileServer(http.Dir("./src")))
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
