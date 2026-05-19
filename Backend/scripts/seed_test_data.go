package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;TrustServerCertificate=true;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}
	defer db.Close()

	fmt.Println("Test verileri ekleniyor...")

	// 1. Test Kullanıcıları Ekle
	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	users := []struct {
		FirstName, LastName, Email, Role string
	}{
		{"Müşteri", "Test", "musteri@test.com", "User"},
		{"Sanatçı", "Test", "sanatci@test.com", "Instructor"},
		{"Yönetici", "Test", "admin@test.com", "Admin"},
	}

	for _, u := range users {
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM Users WHERE Email = @p1", u.Email).Scan(&exists)
		if exists == 0 {
			db.Exec("INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) VALUES (@p1, @p2, @p3, @p4, @p5)",
				u.FirstName, u.LastName, u.Email, string(password), u.Role)
			fmt.Printf("Kullanıcı eklendi: %s\n", u.Email)
		}
	}

	// Kullanıcı ID'lerini alalım
	var userID, artistID, adminID int
	db.QueryRow("SELECT UserID FROM Users WHERE Email = 'musteri@test.com'").Scan(&userID)
	db.QueryRow("SELECT UserID FROM Users WHERE Email = 'sanatci@test.com'").Scan(&artistID)
	db.QueryRow("SELECT UserID FROM Users WHERE Email = 'admin@test.com'").Scan(&adminID)

	// 2. Artist Kaydı Ekle (Eğer yoksa)
	var artistTableID int
	err = db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", artistID).Scan(&artistTableID)
	if err != nil {
		db.Exec("INSERT INTO Artists (UserID, ArtistName, Biography) VALUES (@p1, 'Sanatçı Test', 'Usta bir test sanatçısı.')", artistID)
		db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", artistID).Scan(&artistTableID)
	}

	// 3. Test Eserleri Ekle
	artworks := []struct {
		Title, Category, ImageUrl string
		Price                     float64
	}{
		{"Yıldızlı Gece Çakması", "Yağlı Boya", "https://images.unsplash.com/photo-1579783902614-a3fb3927b6a5", 1500.0},
		{"Modern Karmaşa", "Soyut", "https://images.unsplash.com/photo-1541963463532-d68292c34b19", 2800.0},
	}

	for _, art := range artworks {
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM Artworks WHERE Title = @p1", art.Title).Scan(&exists)
		if exists == 0 {
			db.Exec("INSERT INTO Artworks (Title, ArtistID, Price, ImageURL, Description, Category) VALUES (@p1, @p2, @p3, @p4, 'Test açıklaması', @p5)",
				art.Title, artistTableID, art.Price, art.ImageUrl, art.Category)
			fmt.Printf("Eser eklendi: %s\n", art.Title)
		}
	}

	// 4. Test Atölyeleri Ekle
	workshops := []struct {
		Title, Location string
		Capacity        int
		Price           float64
	}{
		{"Resim Atölyesi", "İstanbul Galeri", 10, 200.0},
		{"Heykel Atölyesi", "Ankara Stüdyo", 5, 450.0},
	}

	for _, ws := range workshops {
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM Workshops WHERE Title = @p1", ws.Title).Scan(&exists)
		if exists == 0 {
			db.Exec("INSERT INTO Workshops (Title, Description, InstructorID, Location, Capacity, Price, ImageUrl, AvailableDates) VALUES (@p1, 'Eğlenceli bir test atölyesi.', @p2, @p3, @p4, @p5, 'https://images.unsplash.com/photo-1460662131955-99f44f6b44f8', '2026-06-01, 2026-06-15')",
				ws.Title, artistID, ws.Location, ws.Capacity, ws.Price)
			fmt.Printf("Atölye eklendi: %s\n", ws.Title)
		}
	}

	// ID'leri çekelim
	var artID, wsID int
	db.QueryRow("SELECT TOP 1 Id FROM Artworks ORDER BY Id DESC").Scan(&artID)
	db.QueryRow("SELECT TOP 1 Id FROM Workshops ORDER BY Id DESC").Scan(&wsID)

	// 5. Test Destek Talepleri
	tickets := []struct {
		Subject, Status string
	}{
		{"Siparişim Nerede?", "Açık"},
		{"Ödeme Hatası", "Çözüldü"},
	}

	for _, t := range tickets {
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM SupportTickets WHERE Subject = @p1 AND UserID = @p2", t.Subject, userID).Scan(&exists)
		if exists == 0 {
			db.Exec("INSERT INTO SupportTickets (UserID, Subject, Message, SupportType, Status) VALUES (@p1, @p2, 'Test mesajı', 'Teknik', @p3)",
				userID, t.Subject, t.Status)
		}
	}

	// 6. Test Yorumları ve Doğrulanmış Alıcı İşlemi
	// Önce bir satın alma ekleyelim ki "doğrulanmış" olsun
	db.Exec("IF NOT EXISTS (SELECT 1 FROM ArtworkPurchases WHERE ArtworkId = @p1 AND UserEmail = 'musteri@test.com') INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, PurchasePrice) VALUES ('musteri@test.com', @p1, 1500.0)", artID)

	comments := []struct {
		Text       string
		Rating     int
		IsVerified bool
	}{
		{"Harika bir eser, bayıldım!", 5, true},
		{"Fiyatı biraz pahalı ama güzel.", 4, false},
		{"Beklediğim gibi çıkmadı.", 2, false},
	}

	for _, c := range comments {
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM Comments WHERE CommentText = @p1 AND TargetID = @p2", c.Text, artID).Scan(&exists)
		if exists == 0 {
			db.Exec("INSERT INTO Comments (UserID, TargetID, TargetType, CommentText, Rating, IsVerified, Upvotes) VALUES (@p1, @p2, 'Artwork', @p3, @p4, @p5, @p6)",
				userID, artID, c.Text, c.Rating, c.IsVerified, c.Rating*2)
		}
	}

	fmt.Println("İşlem başarıyla tamamlandı!")
}
