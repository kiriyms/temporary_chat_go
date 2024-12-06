package models

type IndexData struct {
	Rooms     bool
	UserData  *User
	UserRooms []*Room
}

func NewIndexData() *IndexData {
	return &IndexData{
		Rooms:     false,
		UserData:  NewUser("anonymous", "static/images/avatar_placeholder.png"),
		UserRooms: make([]*Room, 0),
	}
}
