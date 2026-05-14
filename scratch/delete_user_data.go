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
		log.Fatal("Veritabanı bağlantı hatası: ", err)
	}
	defer db.Close()

	userEmail := "eneskrgln6@gmail.com"

	// 1. Eser Satın Alımlarını Sil
	res1, err := db.Exec("DELETE FROM ArtworkPurchases WHERE UserEmail = @p1", userEmail)
	if err != nil {
		log.Printf("ArtworkPurchases silme hatası: %v", err)
	} else {
		count, _ := res1.RowsAffected()
		fmt.Printf("ArtworkPurchases: %d kayıt silindi.\n", count)
	}

	// 2. Atölye Kayıtlarını Sil
	res2, err := db.Exec("DELETE FROM WorkshopEnrollments WHERE UserEmail = @p1", userEmail)
	if err != nil {
		log.Printf("WorkshopEnrollments silme hatası: %v", err)
	} else {
		count, _ := res2.RowsAffected()
		fmt.Printf("WorkshopEnrollments: %d kayıt silindi.\n", count)
	}

	fmt.Println("İşlem tamamlandı.")
}
