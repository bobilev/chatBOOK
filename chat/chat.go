package chat

import (
	"github.com/bobilev/golang-chat-bot-vk"
	"fmt"
	"github.com/bobilev/chatBOOK/dbwork"
	"strconv"
	"strings"
	"os"
	"log"
	"github.com/bobilev/chatBOOK/util"
)
type StatusUser struct {
	Id int
	LastStore int
	LastStoreName string
	LastStep string
	KeyboardLayout string
	Answer map[string]string
}
func (su *StatusUser) Defalt() {
	su.Answer = make(map[string]string)
	su.Answer["0"]  = "ct0"
}
func (su *StatusUser) Start() {
	su.KeyboardLayout = "num"
	su.Answer = make(map[string]string)
	su.Answer["1"]  = "faq"
	su.Answer["2"]  = "ct0"
}
func (su *StatusUser) Clear() {
	for k := range su.Answer {
		delete(su.Answer,k)
	}
	if su.LastStore != 999 {
		su.Answer["0"]  = "ct0"
	}
}
func (su *StatusUser) Continue() {
	if su.LastStore != 0 && su.LastStore != 999 {
		su.Answer["9"]  = "continue"
	}
}
func (su *StatusUser) SetStore(newStore int) {
	su.LastStore = newStore
}
func (su *StatusUser) SetLastStoreName(newLastStoreName string) {
	su.LastStoreName = newLastStoreName
}
func (su *StatusUser) SetStep(newStep string) {
	su.LastStep = newStep
}
func (su *StatusUser) DoneStore() {
	su.SetStore(0)
	su.SetLastStoreName("")
	su.SetStep("0")
	dbwork.UpdateUserStep(su.Id,0,"0")
}
func (su *StatusUser) GetKeyboardLayout() {
	res := dbwork.SelectGetKeyboardLayout(su.Id)
	su.KeyboardLayout = res
}
func (su *StatusUser) SetKeyboardLayout(NewLayout string) {
	su.KeyboardLayout = NewLayout
	dbwork.UpdateSetKeyboardLayout(su.Id,NewLayout)
}
func (su *StatusUser) NewAnswerStep(answers []dbwork.Answer) {
	su.Clear()
	for key,val := range answers {
		if su.LastStore == 999 {
			if su.LastStep == "2" || strings.HasPrefix(su.LastStep,"9") {
				su.Answer[strconv.Itoa(key)] = val.NextStep
			} else {
				su.Answer[strconv.Itoa(key+1)] = val.NextStep
			}
		}
	}
}
func (su *StatusUser) NewAnswerStore(answers []dbwork.Store) {
	su.Clear()
	for key,val := range answers {
		su.Answer[strconv.Itoa(key+1)] = "ct"+strconv.Itoa(val.Storeid)
	}
}
func (su *StatusUser) RecoveryAnswer() {
	if su.LastStore != 0 {
		Step := dbwork.SelectStep(su.LastStore,su.LastStep)
		su.NewAnswerStep(Step.Answers)
	}
}
var mapStatusUsers map[int]*StatusUser
func InitStatusUsers() map[int]*StatusUser{
	mapStatusUsers := make(map[int]*StatusUser)
	allUser := dbwork.SelectAllUsers()
	for id,user := range allUser {
		st := new(StatusUser)
		st.Id = id
		st.LastStore = user.LastStore
		st.LastStep = user.LastStep
		st.KeyboardLayout = user.KeyboardLayout

		st.Defalt()
		st.RecoveryAnswer()
		mapStatusUsers[id] = st
	}
	return mapStatusUsers
}
func InitChatBot() {
	accessToken := os.Getenv("accesstokenvk")
	if accessToken == "" {
		log.Fatal("$accessToken must be set")
	}

	bot := vkchatbot.InitBot(accessToken)
	bot.Log = 2 // 0,1,2 - уровни отображения логов
	updates := bot.StartLongPollServer()

	mapStatusUsers = InitStatusUsers()
	for update := range updates {
		if update.Body != "" {
			//Проверка на нахождения user в локальной базе
			if _,ok := mapStatusUsers[update.UserId]; ok {
				fmt.Println("old [Answer]:",mapStatusUsers[update.UserId].Answer)

				sendText := strings.ToLower(update.Body)
				var okDeRKL bool
				if mapStatusUsers[update.UserId].KeyboardLayout == "рус" {
					sendText , okDeRKL = util.DeRKL(sendText)
					if !okDeRKL {
						sendText = ""
					}
				}
				if nextStep,ok := mapStatusUsers[update.UserId].Answer[sendText]; ok {
					//Определение Step от store или catalog
					if strings.HasPrefix(nextStep,"ct") {//===========================================catalog ctN
						fmt.Println("nextStep[2:]",nextStep[2:])
						if nextStep[2:] == "0" {// Листать каталог дальше
							arrStores := dbwork.SelectStores()
							mapStatusUsers[update.UserId].NewAnswerStore(arrStores)
							//mapStatusUsers[update.UserId].Continue()
							fmt.Println("2[StatusUser]",mapStatusUsers[update.UserId])

							//media
							res0 , _ := bot.SendDocs(update.UserId,SendCategory(arrStores),"")
							fmt.Println("[res]",res0.MessageID)
							//answer
							res1 , _ := bot.SendMessage(update.UserId,ConstructAnswerStore(mapStatusUsers[update.UserId].KeyboardLayout,arrStores))
							fmt.Println("[res]",res1.MessageID)
							//continue
							continueStore := ""

							us := mapStatusUsers[update.UserId].LastStore
							if us != 0 && us != 999 {
								mapStatusUsers[update.UserId].Continue()
								kbl := ""
								switch mapStatusUsers[update.UserId].KeyboardLayout {
								case "num" :
									kbl = "9"
								case "рус":
									kbl = "п"
								}
								continueStore = kbl+" - Продолжить ("+mapStatusUsers[update.UserId].LastStoreName+")\n"
								res2 , _ := bot.SendMessage(update.UserId,continueStore)
								fmt.Println("[res]",res2.MessageID)
							}
						} else{//Загрузить выбраный Store
							Store ,_ := strconv.Atoi(nextStep[2:])
							mapStatusUsers[update.UserId].LastStoreName = dbwork.SelectStoreName(Store)
							SendStep(bot,update,Store,"1")
						}

					} else if nextStep == "continue"{//=========================================================continue
						SendStep(bot,update,mapStatusUsers[update.UserId].LastStore,mapStatusUsers[update.UserId].LastStep)

					} else if nextStep == "faq"{//=========================================================continue
						SendStep(bot,update,999,"1")

					} else if nextStep == "end"{//================================================================== end
						mapStatusUsers[update.UserId].DoneStore()
						arrStores := dbwork.SelectStores()
						mapStatusUsers[update.UserId].NewAnswerStore(arrStores)
						//mapStatusUsers[update.UserId].Continue()
						fmt.Println("2[StatusUser]",mapStatusUsers[update.UserId])

						//media
						res0 , _ := bot.SendDocs(update.UserId,SendCategory(arrStores),"")
						fmt.Println("[res]",res0.MessageID)
						//answer
						res1 , _ := bot.SendMessage(update.UserId,ConstructAnswerStore(mapStatusUsers[update.UserId].KeyboardLayout,arrStores))
						fmt.Println("[res]",res1.MessageID)

					} else {//===================================================================================store N
						SendStep(bot,update,mapStatusUsers[update.UserId].LastStore,nextStep)
					}
				} else if strings.HasPrefix(sendText,"!"){
					if sendText == "!num" {
						mapStatusUsers[update.UserId].SetKeyboardLayout("num")

						res1 , _ := bot.SendMessage(update.UserId,"Раскладка клавиатуры изменена на: num")
						fmt.Println("[res]",res1.MessageID)
						SendStep(bot,update,mapStatusUsers[update.UserId].LastStore,mapStatusUsers[update.UserId].LastStep)
					}
					if sendText == "!рус" {
						mapStatusUsers[update.UserId].SetKeyboardLayout("рус")

						res1 , _ := bot.SendMessage(update.UserId,"Раскладка клавиатуры изменена на: рус")
						fmt.Println("[res]",res1.MessageID)
						SendStep(bot,update,mapStatusUsers[update.UserId].LastStore,mapStatusUsers[update.UserId].LastStep)
					}
				} else {//Answer нет такого
					//answer
					var help string
					fmt.Println("mapStatusUsers[update.Id].KeyboardLayout",mapStatusUsers[update.UserId].KeyboardLayout)
					switch mapStatusUsers[update.UserId].KeyboardLayout {
					case "num":
						help = "0"
					case "рус":
						help = "к"
					}
					res1 , _ := bot.SendMessage(update.UserId,"Нет такой команды, чтобы вернутся в каталог отправьте: "+help)
					fmt.Println("[res]",res1.MessageID)
				}
			} else {
				fmt.Println("Нету в базе")
				dbwork.InsertNewUser(update.UserId,0,"0")
				st := new(StatusUser)
				//st.Defalt()
				st.Start()
				mapStatusUsers[update.UserId] = st

				res , _ := bot.SendMessage(update.UserId,"Добрый день, показать обучение? Это займет всего 1 минуту.\n1 - Давай\n2 - Неа")
				fmt.Println("[res]",res.MessageID)
			}
		}
	}
}
func SendStep(bot *vkchatbot.BotVkApiGroup,update vkchatbot.ObjectUpdate,LastStore int,NextStep string) {
	Step := dbwork.SelectStep(LastStore,NextStep)
	mapStatusUsers[update.UserId].SetStore(LastStore)
	mapStatusUsers[update.UserId].SetStep(Step.StepID)
	mapStatusUsers[update.UserId].NewAnswerStep(Step.Answers)

	dbwork.UpdateUserStep(update.UserId,LastStore,NextStep)

	var Attach vkchatbot.Attachment
	Attach.TypeDoc = Step.TypeDoc
	Attach.MediaId = Step.Media
	Attach.OwnerId = 165847301
	//if LastStore == 999 {//faq
	//	Attach.OwnerId = 164670950
	//}
	Attach.AccessKey = Step.AccessKey

	if Step.Media != 0 {

		res , _ := bot.SendDoc(update.UserId,Attach,Step.Text)
		fmt.Println("[res]",res.MessageID)
	} else {
		res , _ := bot.SendMessage(update.UserId,Step.Text)
		fmt.Println("[res]",res.MessageID)
	}

	res , _ := bot.SendMessage(update.UserId,ConstructAnswer(mapStatusUsers[update.UserId].KeyboardLayout,Step))
	fmt.Println("[res]",res.MessageID)
	fmt.Println("new [Answer]:",mapStatusUsers[update.UserId].Answer)
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
func ConstructAnswer(KeyboardLayout string,Step dbwork.Step) string{
	var answer string
	for k,v := range Step.Answers {
		switch KeyboardLayout {
		case "num" :
			//для FAQ
			if Step.StoreId == 999 && Step.StepID == "2" || Step.StoreId == 999 && strings.HasPrefix(Step.StepID, "9") {
				answer += strconv.Itoa(k)+" - "+v.Text+"\n"
			} else {
				answer += strconv.Itoa(k+1)+" - "+v.Text+"\n"
			}
		case "рус":
			answer += util.RKL(k+1)+" - "+v.Text+"\n"
		}
	}
	return answer
}
func ConstructAnswerStore(KeyboardLayout string,Store []dbwork.Store) string{
	var answer string
	for k,v := range Store {
		switch KeyboardLayout {
		case "num" :
			answer += strconv.Itoa(k+1)+" - "+v.Text+"\n"
		case "рус":
			answer += util.RKL(k+1)+" - "+v.Text+"\n"
		}
	}
	return answer
}
