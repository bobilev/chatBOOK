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
		var keyboardlayout string

		err = res.Scan(&userid,&laststore,&laststep,&keyboardlayout)
		checkErr(err)

		stateuserid.LastStore = laststore
		stateuserid.LastStep = laststep
		stateuserid.KeyboardLayout = keyboardlayout

		mapList[userid] = stateuserid

		fmt.Println("userid         -",userid)
		fmt.Println("laststep       -",laststep)
		fmt.Println("laststore      -",laststore)
		fmt.Println("keyboardlayout -",keyboardlayout)
	}
	return mapList
}
func SelectStores() []Store{
	db := dbConnect()
	defer db.Close()

	res, err := db.Query("SELECT * FROM stores ORDER BY id ASC")
	checkErr(err)
	var Lists []Store
	for res.Next() {
		var storeid int
		var text string
		var media int

		var Store Store
		err = res.Scan(&storeid,&text,&media)
		checkErr(err)

		Store.Storeid = storeid
		Store.Media = media
		Store.Text = text

		Lists = append(Lists,Store)
	}
	return Lists
}
func SelectStoreName(id int) string {
	db := dbConnect()
	defer db.Close()

	var name string
	fmt.Println("DB ->",id)
	err := db.QueryRow("SELECT name FROM public.stores WHERE id = $1", id).Scan(&name)
	checkErr(err)

	return name
}
func SelectStep(storeid int,stepid string) Step{
	db := dbConnect()
	defer db.Close()

	var answer []string
	var step Step
	fmt.Println("DB ->",storeid,stepid)
	err := db.QueryRow("SELECT storeid,stepid,text,media,answer,typedoc,accesskey FROM steps WHERE storeid=$1 AND stepid=$2 ",
		storeid,stepid).Scan(&step.StoreId,&step.StepID,&step.Text,&step.Media,pq.Array(&answer),&step.TypeDoc,&step.AccessKey)
	checkErr(err)

	for _,val := range answer {
		arrey := strings.Split(val,"|")
		ans := Answer{arrey[0],arrey[1]}
		step.Answers = append(step.Answers,ans)
	}
	fmt.Println("-[step]",step)
	return step
}
func SelectGetKeyboardLayout(userid int) string{
	db := dbConnect()
	defer db.Close()

	var KeyboardLayout string
	// SELECT
	err := db.QueryRow("SELECT keyboardlayout FROM users WHERE userid=$1",userid,).Scan(&KeyboardLayout)
	checkErr(err)

	log.Println("{SELECT | KeyboardLayout}",KeyboardLayout)
	return KeyboardLayout
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
func UpdateUserStep(userid int,laststore int,laststep string) {
	db := dbConnect()
	defer db.Close()

	// update
	stmt, err := db.Prepare("UPDATE users SET laststore=$2,laststep=$3 WHERE userid=$1")
	checkErr(err)

	res, err := stmt.Exec(userid, laststore, laststep)
	checkErr(err)

	affect, err := res.RowsAffected()//Сколько записей удалось обновить
	checkErr(err)
	if affect == 0 {
		fmt.Println("[ErrDB: UPDATE] UpdateUserStep")
	}
	fmt.Println(affect)
}
func UpdateSetKeyboardLayout(userid int,NewLayout string) {
	db := dbConnect()
	defer db.Close()

	// update
	stmt, err := db.Prepare("UPDATE users SET keyboardlayout=$2 WHERE userid=$1")
	checkErr(err)

	res, err := stmt.Exec(userid, NewLayout)
	checkErr(err)

	affect, err := res.RowsAffected()//Сколько записей удалось обновить
	checkErr(err)
	if affect == 0 {
		fmt.Println("[ErrDB: UPDATE] UpdateSetKeyboardLayout")
	}
	fmt.Println(affect)
}