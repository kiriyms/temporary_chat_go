package models

import (
	"slices"
	"sync"

	"github.com/google/uuid"
)

type Room struct {
	Id    uuid.UUID
	Users []uuid.UUID
	Hub   *Hub
}

func NewRoom(creatorId uuid.UUID, h *Hub) *Room {
	room := &Room{
		Id:    uuid.New(),
		Users: make([]uuid.UUID, 0),
		Hub:   h,
	}
	room.Users = append(room.Users, creatorId)

	return room
}

type RoomList struct {
	mu    sync.Mutex
	rooms map[uuid.UUID]*Room
}

func NewRoomList() *RoomList {
	return &RoomList{
		rooms: make(map[uuid.UUID]*Room),
	}
}

func (rl *RoomList) AddRoom(r *Room) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.rooms[r.Id] = r
}

func (rl *RoomList) RemoveRoomById(id uuid.UUID) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.rooms, id)
}

func (rl *RoomList) GetRoomById(id uuid.UUID) *Room {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	return rl.rooms[id]
}

func (rl *RoomList) GetUserRooms(uId uuid.UUID) []*Room {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	list := make([]*Room, 0)
	for _, v := range rl.rooms {
		if slices.Contains(v.Users, uId) {
			list = append(list, v)
		}
	}

	return list
}

func (rl *RoomList) GetRoomUsers(rId uuid.UUID) []uuid.UUID {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	return rl.rooms[rId].Users
}
