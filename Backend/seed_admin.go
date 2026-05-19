package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

// Admin Seeder: Bu script terminalden çalıştırılarak manuel admin ekler.
// Kullanım: go run seed_admin.go <firstName> <lastName> <email> <password>
func main() {
	if len(os.Args) < 5 {
		fmt.Println("Kullanım: go run seed_admin.go <Ad> <Soyad> <Email> <Şifre>")
		return
	}

	firstName := os.Args[1]
	lastName := os.Args[2]
	email := os.Args[3]
	password := os.Args[4]

	// Şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Şifre hashleme hatası:", err)
	}

	// Veritabanı bağlantısı
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}
	defer db.Close()

	// Admin kullanıcısını ekle
	query := "INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) VALUES (@p1, @p2, @p3, @p4, 'Admin')"
	_, err = db.Exec(query, firstName, lastName, email, string(hashedPassword))

	if err != nil {
		log.Fatal("Admin ekleme hatası (Email zaten var olabilir):", err)
	}

	fmt.Printf("Başarılı: '%s' email adresiyle Admin kullanıcısı oluşturuldu.\n", email)
}
