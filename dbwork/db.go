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
//createTable creates the table, and if necessary, the database.
func CreateAllTable() error {
	db := dbConnect()
	defer db.Close()

	for _, stmt := range createTableStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
var createTableStatements = []string{
	`CREATE SEQUENCE IF NOT EXISTS public.stores_id_seq
    	INCREMENT 1
    	START 1
    	MINVALUE 1
    	MAXVALUE 9223372036854775807
    	CACHE 1;`,
	`CREATE TABLE IF NOT EXISTS public.users(
    	userid integer,
    	laststore integer,
    	laststep text COLLATE pg_catalog."default"
	)`,
	`CREATE TABLE IF NOT EXISTS public.stores(
    	id integer NOT NULL DEFAULT nextval('stores_id_seq'::regclass),
		name text COLLATE pg_catalog."default",
    	media integer,
    	CONSTRAINT stores_pkey PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS public.steps(
    	id integer NOT NULL DEFAULT nextval('steps2_id_seq'::regclass),
    	storeid integer NOT NULL,
    	stepid text COLLATE pg_catalog."default",
    	text text COLLATE pg_catalog."default",
    	media integer,
    	answer text[] COLLATE pg_catalog."default" NOT NULL,
    	CONSTRAINT steps2_pkey PRIMARY KEY (id)
	)`,
}