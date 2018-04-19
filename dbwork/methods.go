package dbwork

import (
	"log"
	"fmt"
	"github.com/lib/pq"
	"strings"
)

func SelectAllUsers() map[int]StateUser{
	db := dbConnect()
	defer db.Close()

	mapList := make(map[int]StateUser)

	res, err := db.Query("SELECT * FROM users ")
	checkErr(err)
	for res.Next() {
		var userid,laststore int
		var laststep string

		var stateuserid StateUser

		err = res.Scan(&userid,&laststore,&laststep)
		checkErr(err)

		stateuserid.LastStore = laststore
		stateuserid.LastStep = laststep

		mapList[userid] = stateuserid


		fmt.Println("userid -",userid)
		fmt.Println("laststep -",laststep)
		fmt.Println("laststore -",laststore)
	}
	return mapList
}
func SelectStores() map[int]Store{
	db := dbConnect()
	defer db.Close()

	mapList := make(map[int]Store)

	res, err := db.Query("SELECT * FROM stores ")
	checkErr(err)
	for res.Next() {
		var storeid int
		var text string
		var media int

		var Store Store
		err = res.Scan(&storeid,&text,&media)
		checkErr(err)

		Store.Media = media
		Store.Text = text

		mapList[storeid] = Store


		fmt.Println("storeid -",storeid)
		fmt.Println("text -",text)
		fmt.Println("media -",media)
	}
	return mapList
}
func SelectStep(storeid int,stepid string) {
	db := dbConnect()
	defer db.Close()

	var answer []string
	var step Step
	err := db.QueryRow("SELECT storeid,stepid,text,media,answer FROM steps WHERE storeid=$1 AND stepid=$2 ",storeid,stepid).Scan(&step.StoreId,&step.StepID,&step.Text,&step.Media,pq.Array(&answer))
	checkErr(err)

	fmt.Println(answer)

	for _,val := range answer {
		arrey := strings.Split(val,"|")
		ans := Answer{arrey[0],arrey[1]}
		step.Answers = append(step.Answers,ans)
	}
	fmt.Println(step)
	fmt.Println(step.Answers[0].Text)
	fmt.Println(step.Answers[0].NextStep)

}
func InsertNewUser(userid int, laststore int, laststep string) int{
	db := dbConnect()
	defer db.Close()

	var LastInsertId int
	// insert
	err := db.QueryRow("INSERT INTO users(userid,laststore,laststep) values($1,$2,$3) returning userid",userid, laststore, laststep).Scan(&LastInsertId)
	checkErr(err)

	log.Println("{INSERT | NewUser}",LastInsertId)
	return LastInsertId
}