package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Open failed: ", err)
	}
	defer db.Close()

	fmt.Println("\nUsers Sample:")
	rows, _ := db.Query("SELECT TOP 3 UserID, FirstName, LastName, Email FROM Users")
	for rows.Next() {
		var id int
		var f, l, e string
		rows.Scan(&id, &f, &l, &e)
		fmt.Printf("ID: %d, Name: %s %s, Email: %s\n", id, f, l, e)
	}
	rows.Close()

	fmt.Println("\nArtists Sample:")
	rows, _ = db.Query("SELECT TOP 3 ArtistID, FullName FROM Artists")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	rows.Close()
}
