package models

type IndexData struct {
	Rooms        bool
	UserData     *User
	UserRooms    []*RoomWithTimer
	CurrentRooms int
	MaxRooms     int
}

func NewIndexData() *IndexData {
	return &IndexData{
		Rooms:        false,
		UserData:     NewUser("anonymous", "static/images/avatar_placeholder.png"),
		UserRooms:    make([]*RoomWithTimer, 0),
		CurrentRooms: 0,
		MaxRooms:     5,
	}
}
