package models

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id           uuid.UUID
	Users        []uuid.UUID
	Hub          *Hub
	Name         string
	Code         string
	ExpireTime   time.Time
	TimerSeconds int
}

func NewRoom(creatorId uuid.UUID, h *Hub, n string, c string) *Room {
	room := &Room{
		Id:         uuid.New(),
		Users:      make([]uuid.UUID, 0),
		Hub:        h,
		Name:       n,
		Code:       c,
		ExpireTime: time.Now().Add(180 * time.Second),
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

func (rl *RoomList) GetRoomByCode(c string) *Room {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for _, v := range rl.rooms {
		if v.Code == c {
			return v
		}
	}

	return nil
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

func (rl *RoomList) GetRooms() map[uuid.UUID]*Room {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	return rl.rooms
}

func (rl *RoomList) IsRoomCodeUsed(c string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for _, v := range rl.rooms {
		if v.Code == c {
			return true
		}
	}

	return false
}

func (rl *RoomList) AddUserToRoom(rId, uId uuid.UUID) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.rooms[rId].Users = append(rl.rooms[rId].Users, uId)
	return nil
}

func (rl *RoomList) RemoveUserFromRoom(rId, uId uuid.UUID) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	room := rl.rooms[rId]
	if room == nil {
		return fmt.Errorf("room with id %v not found", rId)
	}

	index := -1
	for k, v := range room.Users {
		if uId == v {
			index = k
		}
	}

	if index == -1 {
		return fmt.Errorf("room with id %v does not contain user %v", rId, uId)
	}

	room.Users[index] = room.Users[len(room.Users)-1]
	room.Users = room.Users[:len(room.Users)-1]

	return nil
}

type RoomWithTimer struct {
	Room         Room
	TimerSeconds int
}
