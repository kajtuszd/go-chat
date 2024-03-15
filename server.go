package main

import (
	"errors"
)

type Server struct {
	Rooms map[string]*Room
}

var roomNotFoundError = errors.New("room not found")

type server interface {
	getRoomByID(ID string) (*Room, error)
	addRoom(r *Room) string
	deleteRoom(ID string)
}

func NewServer() *Server {
	return &Server{
		Rooms: make(map[string]*Room),
	}
}

func (s *Server) getRoomByID(ID string) (*Room, error) {
	for _, r := range s.Rooms {
		if r.ID == ID {
			return r, nil
		}
	}
	return nil, roomNotFoundError
}

func (s *Server) addRoom(r *Room) string {
	s.Rooms[r.ID] = r
	return r.ID
}

func (s *Server) deleteRoom(ID string) {
	delete(s.Rooms, ID)
}
