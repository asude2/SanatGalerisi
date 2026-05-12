package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	connString := "server=localhost\\SQLEXPRESS;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Open failed: ", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES")
	if err != nil {
		log.Fatal("Query failed: ", err)
	}
	defer rows.Close()

	fmt.Println("Tables found:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println("- ", name)
	}
}
