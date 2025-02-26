package models

type IndexData struct {
	Rooms        bool
	EditModal    bool
	UserData     *User
	UserRooms    []*RoomWithTimer
	CurrentRooms int
	MaxRooms     int
}

func NewIndexData() *IndexData {
	return &IndexData{
		Rooms:        false,
		EditModal:    false,
		UserData:     NewUser("anonymous", "static/images/avatar_placeholder.png"),
		UserRooms:    make([]*RoomWithTimer, 0),
		CurrentRooms: 0,
		MaxRooms:     5,
	}
}
