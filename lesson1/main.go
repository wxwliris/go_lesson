package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var Db *sql.DB

func init() {
	db, err := sql.Open("mysql",
		"test:test@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("ping ", err)
	}
	Db = db
}

func main() {
	name, err := getAccount(100)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows){
			fmt.Printf("ID = %v has no data, err is %+v\n", 100,err)
			return
		}
	}
	fmt.Println(name)
}

func getAccount(id int) (string, error) {
	var (
		name string
		err  error
	)
	err = Db.QueryRow("select account from user where id = ?", id).Scan(&name)
	err = errors.Wrapf(err, "getAccount")
	return name, err
}
