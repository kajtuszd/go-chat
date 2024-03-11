package main

import "github.com/gorilla/websocket"

type User struct {
	ID       string
	Username string
	Conn     *websocket.Conn
}
