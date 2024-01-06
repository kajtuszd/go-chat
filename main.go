package main

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
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

func handleUserCommunication(user *User, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		messageType, p, err := user.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}
		for _, receiver := range activeUsers {
			if receiver != user {
				log.Println(user.Username, ": ", string(p))
				if err := receiver.Conn.WriteMessage(messageType, p); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := &User{
		ID:       generateRandomID(),
		Username: "Anonymous",
		Conn:     nil,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	user.Conn = conn
	log.Println("Connection open for: ", user.ID)
	if err != nil {
		log.Println(err)
		return
	}

	activeUsers[user.ID] = user
	var wg sync.WaitGroup
	wg.Add(1)

	defer func() {
		wg.Wait()
		if user.Conn != nil {
			user.Conn.Close()
			delete(activeUsers, user.ID)
			log.Println("Connection closed for user:", user.ID)
		}
	}()
	go handleUserCommunication(user, &wg)
}

func main() {
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
