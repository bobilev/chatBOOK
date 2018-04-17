package dbwork

import (
	_ "github.com/lib/pq"
	"database/sql"
	"os"
	"log"
)

func dbConnect() *sql.DB{
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL_CHATBOOK"))
	checkErr(err)
	return db
}






func checkErr(err error) {
	if err != nil {
		log.Println("===========[Err]========\n",err)
	}
}
// createTable creates the table, and if necessary, the database.
//func CreateAllTable() error {
//	db := dbConnect()
//	defer db.Close()
//
//	for _, stmt := range createTableStatements {
//		_, err := db.Exec(stmt)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}