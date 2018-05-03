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
func (su *StatusUser) Defalt() {
	su.SetStore(0)
	su.SetStep("0")
	su.Answer = make(map[string]string)
	su.Answer["0"]  = "ct0"
	su.Answer["00"] = "ct00"
}
func (su *StatusUser) Clear() {
	for k := range su.Answer {
		delete(su.Answer,k)
	}
	su.Answer["0"]  = "ct0"
	su.Answer["00"] = "ct00"
}
func (su *StatusUser) SetStore(newStore int) {
	su.LastStore = newStore
}
func (su *StatusUser) SetStep(newStep string) {
	su.LastStep = newStep
}
func (su *StatusUser) NewAnswerStep(answers []dbwork.Answer) {
	su.Clear()
	for key,val := range answers {
		su.Answer[strconv.Itoa(key+1)] = val.NextStep
	}
}
func (su *StatusUser) NewAnswerStore(answers []dbwork.Store) {
	su.Clear()
	for key,val := range answers {
		su.Answer[strconv.Itoa(key+1)] = "ct"+strconv.Itoa(val.Storeid)
	}
}
var mapStatusUsers map[int]*StatusUser
func InitStatusUsers() map[int]*StatusUser{
	mapStatusUsers := make(map[int]*StatusUser)
	allUser := dbwork.SelectAllUsers()
	for id,user := range allUser {
		st := new(StatusUser)
		st.LastStore = user.LastStore
		st.LastStep = user.LastStep

		st.Defalt()
		mapStatusUsers[id] = st
	}
	return mapStatusUsers
}
func InitChatBot() {
	accessToken := "b25e0478970ebcde8977b7c7b9b8562e28cce81c9f80518b0fa72196fdc0588d833ff6f298a821d12ba18"

	bot := vkchatbot.InitBot(accessToken)
	bot.Log = 3 // 0,1,2 - уровни отображения логов
	updates := bot.StartLongPollServer()

	mapStatusUsers = InitStatusUsers()
	HelloMain := "Добрый день, новичок\n0 - меню\n00 - каталог"
	for update := range updates {

		if update.Body != "" {
			fmt.Println("1[StatusUser]",mapStatusUsers[update.UserId])
			//Проверка на нахождения user в локальной базе
			if _,ok := mapStatusUsers[update.UserId]; ok {
				fmt.Println("Есть в базе")
				if nextStep,ok := mapStatusUsers[update.UserId].Answer[update.Body]; ok {
					//Определение Step от store или catalog
					if strings.HasPrefix(nextStep,"ct") {//===========================================catalog ctN
						fmt.Println("nextStep[2:]",nextStep[2:])
						if nextStep[2:] == "00" {// Листать каталог дальше
							arrStores := dbwork.SelectStores()
							mapStatusUsers[update.UserId].NewAnswerStore(arrStores)
							fmt.Println("2[StatusUser]",mapStatusUsers[update.UserId])

							//media
							res0 , _ := bot.SendDocs(update.UserId,SendCategory(arrStores),"")
							fmt.Println("[res]",res0.MessageID)
							//answer
							res1 , _ := bot.SendMessage(update.UserId,ConstructAnswerStore(arrStores))
							fmt.Println("[res]",res1.MessageID)
						} else if nextStep[2:] == "0" {//Главное меню
							res , _ := bot.SendMessage(update.UserId,HelloMain)
							fmt.Println("[res]",res.MessageID)
						} else{//Загрузить выбраный Store
							Store ,_ := strconv.Atoi(nextStep[2:])
							SendStep(bot,update,Store,"1")
						}

					} else {//===================================================================================store N
						SendStep(bot,update,mapStatusUsers[update.UserId].LastStore,nextStep)
					}
				}else {
					//Answer нет такого
					res , _ := bot.SendMessage(update.UserId,HelloMain)
					fmt.Println("[res]",res.MessageID)
				}
			} else {
				fmt.Println("Нету в базе")
				dbwork.InsertNewUser(update.UserId,0,"0")
				st := new(StatusUser)
				st.Defalt()
				mapStatusUsers[update.UserId] = st
				res , _ := bot.SendMessage(update.UserId,HelloMain)
				fmt.Println("[res]",res.MessageID)
			}
		}
	}
}
func SendStep(bot *vkchatbot.BotVkApiGroup,update vkchatbot.ObjectUpdate,LastStore int,NextStep string) {
	Step := dbwork.SelectStep(LastStore,NextStep)
	mapStatusUsers[update.UserId].SetStep(Step.StepID)
	mapStatusUsers[update.UserId].SetStore(LastStore)
	mapStatusUsers[update.UserId].NewAnswerStep(Step.Answers)
	fmt.Println("2[StatusUser]",mapStatusUsers[update.UserId])
	var Attach vkchatbot.Attachment
	Attach.TypeDoc = Step.TypeDoc
	Attach.MediaId = Step.Media
	Attach.OwnerId = 165847301
	Attach.AccessKey = Step.AccessKey

	if Step.Media != 0 {

		res , _ := bot.SendDoc(update.UserId,Attach,Step.Text)
		fmt.Println("[res]",res.MessageID)
	} else {
		res , _ := bot.SendMessage(update.UserId,Step.Text)
		fmt.Println("[res]",res.MessageID)
	}

	res , _ := bot.SendMessage(update.UserId,ConstructAnswer(Step))
	fmt.Println("[res]",res.MessageID)
}
func SendCategory(arrStores []dbwork.Store) []vkchatbot.Attachment {

	var arrAttach []vkchatbot.Attachment

	for _,store := range arrStores {
		var Attach vkchatbot.Attachment
		Attach.MediaId = store.Media
		Attach.TypeDoc = "photo"
		Attach.OwnerId = 164670950
		arrAttach = append(arrAttach,Attach)
	}
	return arrAttach

}
func ConstructAnswer(Step dbwork.Step) string{
	var answer string
	for k,v := range Step.Answers {
		answer += strconv.Itoa(k+1)+" - "+v.Text+"\n"
	}
	//answer += "____________.__________\n0 - меню | 00 - каталог"
	return answer
}
func ConstructAnswerStore(Store []dbwork.Store) string{
	var answer string
	for k,v := range Store {
		answer += strconv.Itoa(k+1)+" - "+v.Text+"\n"
	}
	//answer += "____.____\n0 - меню"
	return answer
}
