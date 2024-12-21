package utils

import (
	"log"
	"time"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/google/uuid"
)

func StartRoomTimer(rl *models.RoomList, rId uuid.UUID) {
	timer := time.NewTimer(20 * time.Second)
	log.Printf("TIMER: starting room %v timer", rId)
	go func() {
		<-timer.C
		log.Printf("TIMER: room %v timer ran out", rId)
		room := rl.GetRoomById(rId)
		// room.Hub.Stop()
		close(room.Hub.Broadcast)
		rl.RemoveRoomById(rId)
	}()
}
