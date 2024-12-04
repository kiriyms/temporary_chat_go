package models

type IndexData struct {
	Rooms    bool
	UserData *User
}

func NewIndexData() *IndexData {
	return &IndexData{
		Rooms:    false,
		UserData: NewUser("anonymous", "static/images/avatar_placeholder.png"),
	}
}
