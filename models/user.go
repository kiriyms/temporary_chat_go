package models

import (
	"sync"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID
	Name       string
	AvatarPath string
}

func NewUser(name, avatarPath string) *User {
	return &User{
		Id:         uuid.New(),
		Name:       name,
		AvatarPath: avatarPath,
	}
}

type UserList struct {
	mu    sync.Mutex
	users map[uuid.UUID]*User
}

func NewUserList() *UserList {
	return &UserList{
		users: make(map[uuid.UUID]*User),
	}
}

func (ul *UserList) AddUser(u *User) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	ul.users[u.Id] = u
}

func (ul *UserList) RemoveUserById(id uuid.UUID) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	delete(ul.users, id)
}

func (ul *UserList) GetUserById(id uuid.UUID) *User {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	return ul.users[id]
}
