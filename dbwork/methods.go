package dbwork

import (
	"strconv"
	"log"
	"fmt"
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
		var laststore string
		var stateuserid StateUser

		err = res.Scan(&userid,&laststep,&laststore)
		checkErr(err)
		stateuserid.LastStore , _ = strconv.Atoi(laststore)
		stateuserid.LastStep = laststep

		mapList[userid] = stateuserid


		fmt.Println("userid -",userid)
		fmt.Println("laststep -",laststep)
		fmt.Println("laststore -",laststore)
	}
	return mapList
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