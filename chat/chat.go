package chat

import (
	"github.com/bobilev/golang-chat-bot-vk"
	"fmt"
	"github.com/bobilev/chatBOOK/dbwork"
	"strconv"
	"strings"
)
type StatusUser struct {
	LastStore int
	LastStep string
	Answer map[string]string
}
func (su *StatusUser) Clear() {
	for k := range su.Answer {
		delete(su.Answer,k)
	}
}
func InitStatusUsers() map[int]StatusUser{
	mapStatusUsers := make(map[int]StatusUser)
	allUser := dbwork.SelectAllUsers()
	for id,user := range allUser {
		var st StatusUser
		st.LastStore = user.LastStore
		st.LastStep = user.LastStep
		st.Answer = make(map[string]string)
		mapStatusUsers[id] = st
	}
	return mapStatusUsers
}
func InitChatBot() {
	accessToken := "b25e0478970ebcde8977b7c7b9b8562e28cce81c9f80518b0fa72196fdc0588d833ff6f298a821d12ba18"

	bot := vkchatbot.InitBot(accessToken)
	bot.Log = 1 // 0,1,2 - уровни отображения логов
	updates := bot.StartLongPollServer()

	mapStatusUsers := InitStatusUsers()
	for update := range updates {
		//fmt.Println("UserID:",update.UserId,"Text Message:",update.Body)
		if update.Body == "hi" {
			dbwork.SelectStep(3,"1")


			//res , _ := bot.SendMessage(update.UserId,"Hello")
			//fmt.Println("[res]",res.MessageID)
		}
		//if update.Body == "sex" {
		//	res , _ := bot.SendDoc(update.UserId,"photo",456239017,"секси эльфийка")
		//	fmt.Println("[res]",res.MessageID)
		//}
		if update.Body != "" {
			if _,ok := mapStatusUsers[update.UserId]; ok {
				fmt.Println("Есть в базе")
				if _,ok := mapStatusUsers[update.UserId].Answer[update.Body]; ok {
					//lastStore := mapStatusUsers[update.UserId].LastStore
					nextStep := mapStatusUsers[update.UserId].Answer[update.Body]

					//Определение Step от store или catalog
					if strings.HasPrefix(nextStep,"ct") {
						//catalog ct00

					} else {
						//store 00

						//var answer string
						//
						//for k,v := range step.Answer {
						//	answer += k+" - "+v
						//}
					}
				}else {
					//Answer нет такого
				}

			} else {
				dbwork.InsertNewUser(update.UserId,0,"")
				mapStatusUsers[update.UserId] = StatusUser{0,"",nil}
				fmt.Println("Нету в базе")
				res , _ := bot.SendMessage(update.UserId,"Добрый день, новичок")
				fmt.Println("[res]",res.MessageID)
			}

		}
		if update.Body == "00" {
			mapStores := dbwork.SelectStores()
			var catalog string
			i := 1
			for _, store := range mapStores {
				catalog += strconv.Itoa(i) +" - "+ store.Text+"\n"
				i++
			}
			res , _ := bot.SendMessage(update.UserId,catalog)
			fmt.Println("[res]",res.MessageID)
		}

	}
}
