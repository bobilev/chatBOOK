package dbwork

import "strconv"

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

		err = res.Scan(&userid,&laststep,&laststore)
		checkErr(err)
		var stateuserid StateUser
		stateuserid.LastStore , _ = strconv.Atoi(laststore)
		stateuserid.LastStep = laststep

		mapList[userid] = stateuserid


		//fmt.Println("LogSelectInfo - date",date)
	}
	return mapList
}