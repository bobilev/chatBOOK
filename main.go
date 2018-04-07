package main

import (
	"log"
	"github.com/Dimonchik0036/vk-api"
)

func main() {
	client, err := vkapi.NewClientFromToken("abc574aa7a283eeb99278f25c67d5d2e7b777c536dfa45778e6199916050cceaaad29de2991021da27800")
	if err != nil {
		log.Panic(err)
	}

	client.Log(true)

	if err := client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, _, err := client.GetLPUpdatesChan(100, vkapi.LPConfig{25, vkapi.LPModeAttachments})
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil || !update.IsNewMessage() || update.Message.Outbox(){
			continue
		}

		log.Printf("%s", update.Message.String())
		if update.Message.Text == "/start" {
			client.SendMessage(vkapi.NewMessage(vkapi.NewDstFromUserID(update.Message.FromID), "Hello!"))
		}

	}
}
