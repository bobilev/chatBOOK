package dbwork

import (
	"log"
	"fmt"
	"github.com/lib/pq"
)

func SelectAllUsers() map[int]StateUser{
	db := dbConnect()
	defer db.Close()

	mapList := make(map[int]StateUser)

	res, err := db.Query("SELECT * FROM users ")
	checkErr(err)
	for res.Next() {
		var userid int
		var laststep string
		var laststore int
		var stateuserid StateUser

		err = res.Scan(&userid,&laststep,&laststore)
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

	var answers []string
	err := db.QueryRow("SELECT answers FROM steps2 WHERE storeid=$1 AND stepid=$2 ",storeid,stepid).Scan(pq.Array(&answers))
	checkErr(err)
	fmt.Println(answers)

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