package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/google/uuid"
)

func NewUserTimer(ul *models.UserList, uId uuid.UUID) chan (bool) {
	resetChan := make(chan (bool))
	timer := time.NewTimer(15 * time.Second)

	log.Printf("TIMER: starting user %v timer", uId)
	go func() {
		for {
			select {
			case <-timer.C:
				log.Printf("TIMER: user %v timer ran out", uId)
				avatarPath := ul.GetUserById(uId).AvatarPath
				if avatarPath == "static/images/avatar_placeholder.png" {
					ul.RemoveUserById(uId)
					return
				}
				err := os.Remove(fmt.Sprintf("%v", avatarPath))
				if err != nil {
					log.Printf("[ERROR]: could not remove image %v. ErrMsg: %v", avatarPath, err)
				}
				ul.RemoveUserById(uId)
				return
			case <-resetChan:
				log.Printf("TIMER: user %v timer resetting", uId)
				timer.Reset(15 * time.Second)
			}
		}
	}()
	return resetChan
}
