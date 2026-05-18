package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int     `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	UserRole  string  `json:"role"` // 'User' veya 'Instructor'
	Biography string  `json:"biography"`
	Balance   float64 `json:"balance"`
}

// 2. Login Yanıtı (Frontend'e rolü de söylemeliyiz)
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Role    string `json:"role"` // Vue tarafı buna bakıp butonları gösterecek
	UserID  int    `json:"userId"`
	Email   string `json:"email"`
}

// 3. Ekleme isteği için
type CreateArtworkRequest struct {
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"imageUrl"`
	Description string  `json:"description"`
	ArtistID    int     `json:"artistId"`
}
type Artwork struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	ArtistID     int     `json:"artistId"`
	Artist       string  `json:"artist"`
	Price        float64 `json:"price"`
	ImageUrl     string  `json:"imageUrl"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	IsCampaign   bool    `json:"IsCampaign"`
	DiscountRate int     `json:"DiscountRate"`
	IsSold       bool    `json:"issold"`
}
type Workshop struct {
	Id             int     `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	InstructorName string  `json:"instructorName"`
	Location       string  `json:"location"`
	Capacity       int     `json:"capacity"`
	Price          float64 `json:"price"`
	ImageUrl       string  `json:"image"`
	AvailableDates string  `json:"availableDates"`
}

// Eğitmenin kendi atölyelerini döndürmek için struct
type MyWorkshop struct {
	Id             int     `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Location       string  `json:"location"`
	Capacity       int     `json:"capacity"`
	Price          float64 `json:"price"`
	ImageUrl       string  `json:"imageUrl"`
	AvailableDates string  `json:"availableDates"`
}
type EnrollmentRequest struct {
	Email            string `json:"email"`
	WorkshopId       int    `json:"workshopId"`
	ParticipantCount int    `json:"participantCount"`
	ReservedDate     string `json:"reservedDate"`
}
type BuyArtwork struct {
	Email     string  `json:"email"`
	ArtworkId int     `json:"artworkId"`
	Price     float64 `json:"price"`
}
type ArtworkRequest struct {
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageUrl    string  `json:"imageUrl"`
	Description string  `json:"description"`
	Email       string  `json:"email"`
}
type WorkshopRequest struct {
	Email          string  `json:"email"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Location       string  `json:"location"`
	Capacity       int     `json:"capacity"`
	Price          float64 `json:"price"`
	ImageUrl       string  `json:"imageUrl"`
	AvailableDates string  `json:"availableDates"`
}

// ! KAYIT
func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
		return
	}

	if u.UserRole == "" {
		u.UserRole = "User"
	}

	// 1. Şifreyi Hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Şifre işlenemedi", http.StatusInternalServerError)
		return
	}

	// 2. Veritabanına Bağlan
	connString := "sqlserver://localhost:1433?database=SanatProjesi&trusted_connection=yes&encrypt=disable"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Println("Bağlantı hatası:", err)
		http.Error(w, "Veritabanı bağlantı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 3. Veriyi Kaydet
	query := "INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) VALUES (@p1, @p2, @p3, @p4, @p5)"
	_, err = db.Exec(query, u.FirstName, u.LastName, u.Email, string(hashedPassword), u.UserRole)

	if err != nil {
		log.Println("DB Yazma Hatası:", err)
		// Eğer e-posta zaten varsa UNIQUE kısıtlamasından dolayı hata verebilir
		http.Error(w, "Kayıt yapılamadı (E-posta kullanımda olabilir)", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kayıt başarıyla tamamlandı!",
		"role":    u.UserRole,
	})
}

// ! GİRİŞ
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	// 1. Veritabanı Bağlantısı
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var dbPassword string
	var firstName string
	var userRole string

	// Sorgu Güncellendi: UserRole bilgisini de çekiyoruz
	query := "SELECT Password, FirstName, UserRole FROM Users WHERE Email = @p1"
	err = db.QueryRow(query, creds.Email).Scan(&dbPassword, &firstName, &userRole)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	// 2. Şifre Karşılaştırma
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Hatalı şifre", http.StatusUnauthorized)
		return
	}

	// 3. Token Oluşturma (Rol bilgisini içine ekledik)
	var jwtKey = []byte("cok_gizli_anahtar_123")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     creds.Email,
		"firstName": firstName,
		"role":      userRole, // Vue tarafı için token içinde rolü saklıyoruz
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// 4. Yanıt Gönderme (JSON içine rolü ekledik)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Giriş başarılı!",
		"token":     tokenString,
		"role":      userRole, // 'User' veya 'Instructor' döner
		"firstName": firstName,
	})
}

// ! ESERLERİ GETİRME (IsSold Kolonu Eklendi ve Düzeltildi)
func getArtworksHandler(w http.ResponseWriter, r *http.Request) {
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı bağlantı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 🚩 ISNULL(aw.IsSold, 0) SELECT sorgusuna eklendi!
	rows, err := db.Query(`
            SELECT aw.Id, aw.Title, aw.ArtistID, aw.Price, aw.ImageUrl,
                ISNULL(aw.Description, ''), ISNULL(aw.Category, ''),
                ISNULL(ar.ArtistName, 'Bilinmeyen Sanatçı'),
                ISNULL(aw.IsCampaign, 0), ISNULL(aw.DiscountRate, 0), ISNULL(aw.IsSold, 0)
            FROM Artworks aw
            LEFT JOIN Artists ar ON aw.ArtistID = ar.ArtistID
    `)

	if err != nil {
		http.Error(w, "Sorgu hatası", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	artworks := []Artwork{}
	for rows.Next() {
		var a Artwork
		err := rows.Scan(
			&a.ID, &a.Title, &a.ArtistID, &a.Price, &a.ImageUrl,
			&a.Description, &a.Category, &a.Artist,
			&a.IsCampaign, &a.DiscountRate, &a.IsSold, // 🚩 Tarayıcıya satıldı bilgisi artık gidiyor
		)
		if err != nil {
			fmt.Println("Scan hatası:", err)
			continue
		}
		artworks = append(artworks, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artworks)
}

// ! SATIN ALMA HANDLER (Geliştirilmiş ve Hataları Ayıklanmış Versiyon)
func buyArtworkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	pathID := strings.TrimPrefix(r.URL.Path, "/buy-artwork/")
	artworkId, err := strconv.Atoi(pathID)
	if err != nil || artworkId == 0 {
		http.Error(w, "Geçersiz Eser ID'si", http.StatusBadRequest)
		return
	}

	var data struct {
		Email         string  `json:"email"`
		Price         float64 `json:"price"`
		PaymentMethod string  `json:"paymentMethod"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Geçersiz veri formatı", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "İşlem başlatılamadı", http.StatusInternalServerError)
		return
	}

	// --- ADIM 1: Eser Durumu Kontrolü (OwnerID hatası düzeltildi, ArtistID kullanılıyor) ---
	var isSold bool
	var artistId int

	// 🕵️‍♂️ Temizlenen DB şemasına uygun olarak ArtistID çekiliyor
	checkQuery := "SELECT IsSold, ArtistID FROM Artworks WHERE Id = @p1"
	err = tx.QueryRow(checkQuery, artworkId).Scan(&isSold, &artistId)
	if err != nil {
		tx.Rollback()
		log.Println("Eser bulunamadı hatası:", err)
		http.Error(w, "Eser bulunamadı", http.StatusNotFound)
		return
	}

	if isSold {
		tx.Rollback()
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Bu eser çoktan satılmış! 😔"})
		return
	}

	// --- ADIM 2: Bakiye Kontrol ve Düşme İşlemi ---
	var userBalance float64
	var userId int
	err = tx.QueryRow("SELECT UserID, ISNULL(Balance, 0) FROM Users WHERE Email = @p1", data.Email).Scan(&userId, &userBalance)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Kullanıcı hesabı doğrulanamadı", http.StatusUnauthorized)
		return
	}

	if userBalance < data.Price {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Yetersiz bakiye! 💸 Lütfen bakiye yükleyin."})
		return
	}

	// Bakiyeyi düşüyoruz
	_, err = tx.Exec("UPDATE Users SET Balance = Balance - @p1 WHERE UserID = @p2", data.Price, userId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Bakiye düşülürken hata oluştu", http.StatusInternalServerError)
		return
	}

	// --- ADIM 3: Eseri Satıldı Olarak İşaretle (Onay Bekliyor Aşaması İçin) ---
	// Not: İrem satın alınca anasayfada gizlenmesi istendiği için geçici olarak IsSold=1 yapıyoruz, Asude onaylayınca kesinleşecek.
	_, err = tx.Exec("UPDATE Artworks SET IsSold = 1 WHERE Id = @p1", artworkId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Eser durumu güncellenemedi", http.StatusInternalServerError)
		return
	}

	// --- ADIM 4: Satın Alma Kaydı Oluşturma (SellerID artık ArtistID'nin bağlı olduğu UserID'dir) ---
	var sellerUserId int
	err = tx.QueryRow("SELECT UserID FROM Artists WHERE ArtistID = @p1", artistId).Scan(&sellerUserId)
	if err != nil {
		sellerUserId = userId // Fallback durumunda işlemi engellemesin diye mevcut kullanıcıyı ata
	}

	purchaseQuery := `
        INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, PurchasePrice, PaymentMethod, SellerID, Status) 
        VALUES (@p1, @p2, @p3, @p4, @p5, 'Hazırlanıyor')`

	_, err = tx.Exec(purchaseQuery, data.Email, artworkId, data.Price, data.PaymentMethod, sellerUserId)
	if err != nil {
		tx.Rollback()
		log.Println("Purchase tablosuna yazma hatası:", err)
		http.Error(w, "Kayıt oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// --- ADIM 5: İşlemleri Onayla ---
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "İşlem tamamlanamadı", http.StatusInternalServerError)
		return
	}

	// --- ADIM 6: Öneri Sistemini Hafızaya Al ---
	updateCatQuery := `
            UPDATE Users 
            SET LastPurchasedCategory = (SELECT Category FROM Artworks WHERE Id = @p1) 
            WHERE Email = @p2`
	_, _ = db.Exec(updateCatQuery, artworkId, data.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Başarılı! 🎉"})
}

// ! PROFİL BİLGİLERİNİ GETİRME
func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var u User
	var biography sql.NullString
	var balance sql.NullFloat64

	// Users ve Artists tablosundan biography'yi ve balance'ı al
	query := `
		SELECT u.FirstName, u.LastName, u.Email, ISNULL(a.Biography, ''), ISNULL(u.Balance, 0)
		FROM Users u
		LEFT JOIN Artists a ON u.UserID = a.UserID
		WHERE u.Email = @p1
	`
	err = db.QueryRow(query, email).Scan(&u.FirstName, &u.LastName, &u.Email, &biography, &balance)

	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	if biography.Valid {
		u.Biography = biography.String
	}

	if balance.Valid {
		u.Balance = balance.Float64
	} else {
		u.Balance = 0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// ! PROFİL BİLGİLERİNİ GÜNCELLEME
func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// 1. Frontend'den gelen güncel bilgileri oku
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// E-posta adresini anahtar olarak kullanıp diğer alanları güncelliyoruz
	query := "UPDATE Users SET FirstName = @p1, LastName = @p2 WHERE Email = @p3"
	_, err = db.Exec(query, u.FirstName, u.LastName, u.Email)

	if err != nil {
		http.Error(w, "Güncelleme başarısız", http.StatusInternalServerError)
		return
	}

	// Eğer biography varsa Artists tablosuna kaydet
	if u.Biography != "" {
		var userID int
		err = db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", u.Email).Scan(&userID)
		if err == nil {
			// Önce kontrolü yap - artist kaydı var mı?
			var artistID int
			err = db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", userID).Scan(&artistID)

			if err != nil { // Yoksa ekle
				db.Exec("INSERT INTO Artists (UserID, Biography) VALUES (@p1, @p2)", userID, u.Biography)
			} else { // Varsa güncelle
				db.Exec("UPDATE Artists SET Biography = @p1 WHERE UserID = @p2", u.Biography, userID)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Profil güncellendi!"})
}

// ! ŞİFRE DEĞİŞTİRME
func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Frontend'den gelecek veri yapısı
	type PasswordRequest struct {
		Email       string `json:"email"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	var req PasswordRequest
	json.NewDecoder(r.Body).Decode(&req)

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	// 1. Mevcut hash'li şifreyi veritabanından getir
	var storedPassword string
	err := db.QueryRow("SELECT Password FROM Users WHERE Email = @p1", req.Email).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	// 2. Eski şifre doğru mu? (bcrypt ile kontrol)
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.OldPassword))
	if err != nil {
		http.Error(w, "Eski şifre hatalı", http.StatusUnauthorized)
		return
	}

	// 3. Yeni şifreyi hash'le
	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	// 4. Veritabanında güncelle
	_, err = db.Exec("UPDATE Users SET Password = @p1 WHERE Email = @p2", string(newHash), req.Email)
	if err != nil {
		http.Error(w, "Şifre güncellenemedi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Şifre başarıyla değiştirildi!"})
}

// ! ATÖLYELERİ GETİRME
func getWorkshopsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
        SELECT 
            w.Id, w.Title, ISNULL(w.Description, ''), 
            (u.FirstName + ' ' + u.LastName) as InstructorName, 
            ISNULL(w.AvailableDates, ''), 
            w.Location, w.Capacity, w.Price, ISNULL(w.ImageUrl, '') 
        FROM Workshops w
        JOIN Users u ON w.InstructorID = u.UserID
    `
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Sorgu hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workshops := []Workshop{}
	for rows.Next() {
		var ws Workshop
		err := rows.Scan(
			&ws.Id,
			&ws.Title,
			&ws.Description,
			&ws.InstructorName,
			&ws.AvailableDates,
			&ws.Location,
			&ws.Capacity,
			&ws.Price,
			&ws.ImageUrl,
		)

		if err != nil {
			fmt.Println("Scan hatası:", err)
			continue
		}

		workshops = append(workshops, ws)
	}

	if workshops == nil {
		workshops = []Workshop{}
	}
	json.NewEncoder(w).Encode(workshops)
}

// ! FAVORİLERE EKLEME
func addFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		ArtworkId int    `json:"artworkId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//favorilere ekle(eğer zaten varsa UNIQUE constraint hata verecektir)
	query := "INSERT INTO Favorites (UserEmail, ArtworkId) VALUES (@p1, @p2)"
	_, err = db.Exec(query, req.Email, req.ArtworkId)

	if err != nil {
		http.Error(w, "Bu eser zaten favorilerinizde veya bir hata oluştu", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Favorilere eklendi! ❤️"})
}

// ! FAVORİLERDEN ÇIKARMA
func deleteFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		ArtworkId int    `json:"artworkId"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	// Favoriyi sil
	query := "DELETE FROM Favorites WHERE UserEmail = @p1 AND ArtworkId = @p2"
	_, err := db.Exec(query, req.Email, req.ArtworkId)

	if err != nil {
		http.Error(w, "Silme işlemi başarısız", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Favorilerden çıkarıldı! 💔"})
}

// ! FAVORİ Mİ KONTROLÜ
func checkFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	email := r.URL.Query().Get("email")
	artworkId := r.URL.Query().Get("artworkId")

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var count int
	// Bu ikilinin tabloda olup olmadığını sayıyoruz
	query := "SELECT COUNT(*) FROM Favorites WHERE UserEmail = @p1 AND ArtworkId = @p2"
	err = db.QueryRow(query, email, artworkId).Scan(&count)

	exists := false
	if err == nil && count > 0 {
		exists = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"isFavorite": exists})
}

// ! FAVORİ LİSTESİNİ GETİRME (Arkadaşının Kodunun En Kararlı ve Uyumlu Hali)
func getUserFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parametresi gerekli", http.StatusBadRequest)
		return
	}

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 🚩 SQL'den tüm alanları eksiksiz ve struct sırasına göre çekiyoruz
	query := `
        SELECT aw.Id, aw.Title, aw.ArtistID, aw.Price, aw.ImageUrl, 
               ISNULL(aw.Description, ''), ISNULL(aw.Category, ''),
               ISNULL(ar.ArtistName, 'Bilinmeyen Sanatçı'),
               ISNULL(aw.IsCampaign, 0), ISNULL(aw.DiscountRate, 0), ISNULL(aw.IsSold, 0)
        FROM Favorites f
        JOIN Artworks aw ON f.ArtworkId = aw.Id
        LEFT JOIN Artists ar ON aw.ArtistID = ar.ArtistID
        WHERE f.UserEmail = @p1
    `
	rows, err := db.Query(query, email)
	if err != nil {
		http.Error(w, "Sorgu hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var favs []Artwork
	for rows.Next() {
		var a Artwork
		// 🚀 BURASI ALTIN DEĞERİNDE: BIT (bool) alanları Go doğrudan yakalıyor, hata vermiyor!
		err := rows.Scan(
			&a.ID, &a.Title, &a.ArtistID, &a.Price, &a.ImageUrl,
			&a.Description, &a.Category, &a.Artist,
			&a.IsCampaign, &a.DiscountRate, &a.IsSold,
		)
		if err != nil {
			fmt.Println("Favori Scan hatası:", err)
			continue
		}
		favs = append(favs, a)
	}

	if favs == nil {
		favs = []Artwork{}
	}

	json.NewEncoder(w).Encode(favs)
}

// ! ATÖLYEYE KAYIT OLMA
func enrollWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	var req EnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı bağlantı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Veritabanına kayıt ekleme (Katılımcı sayısı ve Tarih dahil)
	query := `INSERT INTO WorkshopEnrollments (UserEmail, WorkshopId, ParticipantCount, ReservedDate) 
              VALUES (@p1, @p2, @p3, @p4)`

	_, err = db.Exec(query, req.Email, req.WorkshopId, req.ParticipantCount, req.ReservedDate)
	if err != nil {
		http.Error(w, "Rezervasyon oluşturulamadı. Lütfen tekrar deneyin.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyonunuz başarıyla oluşturuldu! 🎉"})
}

// ! KULLANICININ ATÖLYE KAYITLARINI GETİRME
func getUserEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parametresi gerekli", http.StatusBadRequest)
		return
	}

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
        SELECT e.Id, w.Title, e.ParticipantCount, e.ReservedDate, e.CreatedAt, w.Location, w.AvailableDates
        FROM WorkshopEnrollments e
        JOIN Workshops w ON e.WorkshopId = w.Id
        WHERE e.UserEmail = @p1
        ORDER BY e.CreatedAt DESC`

	rows, err := db.Query(query, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var enrollments []map[string]interface{}
	for rows.Next() {
		var id, pCount int
		var title, rDate, cAt, location, aDates string
		rows.Scan(&id, &title, &pCount, &rDate, &cAt, &location, &aDates)

		enrollments = append(enrollments, map[string]interface{}{
			"id":               id,
			"workshopTitle":    title,
			"participantCount": pCount,
			"reservedDate":     rDate,
			"createdAt":        cAt,
			"location":         location,
			"availableDates":   aDates,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// ! ATÖLYE KAYIT BİLGİLERİNİ GÜNCELLEME
func updateEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Sadece PUT metodu destekleniyor", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID               int    `json:"id"`
		ParticipantCount int    `json:"participantCount"`
		ReservedDate     string `json:"reservedDate"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "UPDATE WorkshopEnrollments SET ParticipantCount = @p1, ReservedDate = @p2 WHERE Id = @p3"
	_, err = db.Exec(query, data.ParticipantCount, data.ReservedDate, data.ID)

	if err != nil {
		http.Error(w, "Güncelleme hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon başarıyla güncellendi"})
}

// ! ATÖLYE KAYIT SİLME
func deleteEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Sadece DELETE metodu destekleniyor", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "DELETE FROM WorkshopEnrollments WHERE Id = @p1"
	_, err = db.Exec(query, id)

	if err != nil {
		http.Error(w, "Silme hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon başarıyla iptal edildi"})
}

// ! KULLANICININ SATIN ALDIĞI ESERLERİ GETİRME
func getUserPurchasesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	email := r.URL.Query().Get("email")

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// Artworks tablosuyla birleştirip başlığı ve resmi de alıyoruz
	query := `
        SELECT p.Id, a.Title, p.PurchasePrice, p.Status, p.CreatedAt, a.ImageUrl 
        FROM ArtworkPurchases p
        JOIN Artworks a ON p.ArtworkId = a.Id
        WHERE p.UserEmail = @p1`

	rows, _ := db.Query(query, email)
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var title, status, cAt, img string
		var price float64
		rows.Scan(&id, &title, &price, &status, &cAt, &img)
		results = append(results, map[string]interface{}{
			"id": id, "title": title, "price": price, "status": status, "date": cAt, "image": img,
		})
	}
	json.NewEncoder(w).Encode(results)
}

// ! ESER EKLEME (INSTRUCTOR İÇİN)
func addArtworkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req ArtworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Kullanıcı e-posta bilgisi eksik", http.StatusBadRequest)
		return
	}

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// 1. ADIM: Kullanıcının UserID'sini bulalım (Email üzerinden)
	var userID int
	var firstName string
	var lastName string
	err := db.QueryRow("SELECT UserID, FirstName, LastName FROM Users WHERE Email = @p1", req.Email).Scan(&userID, &firstName, &lastName)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	// 2. ADIM: Bu kullanıcı Artists tablosunda var mı? Yoksa ekle.
	var artistID int
	artistName := firstName + " " + lastName

	err = db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", userID).Scan(&artistID)

	if err != nil { // Kayıtlı değilse sanatçı olarak ekle
		err = db.QueryRow("INSERT INTO Artists (UserID, ArtistName) OUTPUT INSERTED.ArtistID VALUES (@p1, @p2)",
			userID, artistName).Scan(&artistID)
	} else { // Kayıtlı ise ismi güncelleyelim
		_, err = db.Exec("UPDATE Artists SET ArtistName = @p1 WHERE ArtistID = @p2", artistName, artistID)
	}

	// 3. ADIM: Eseri artık gerçek ArtistID ile ekle
	query := `INSERT INTO Artworks (Title, ArtistID, Price, ImageURL, Description, Category) 
              VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`

	_, err = db.Exec(query, req.Title, artistID, req.Price, req.ImageUrl, req.Description, req.Category)

	if err != nil {
		http.Error(w, "Eser kaydedilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Eseriniz başarıyla yayınlandı!"})
}

// ! ESER SİLME
func deleteArtworkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ArtworkID int    `json:"artworkId"`
		Email     string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// 1. ADIM: Kullanıcının UserID'sini al
	var userID int
	err := db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", req.Email).Scan(&userID)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	// 2. ADIM: Eseri kontrol et - bu eser bu kullanıcının mı?
	var artworkArtistID int
	err = db.QueryRow("SELECT ArtistID FROM Artworks WHERE Id = @p1", req.ArtworkID).Scan(&artworkArtistID)
	if err != nil {
		http.Error(w, "Eser bulunamadı", http.StatusNotFound)
		return
	}

	// 3. ADIM: Kullanıcının artist kaydını kontrol et
	var artistID int
	err = db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", userID).Scan(&artistID)
	if err != nil || artistID != artworkArtistID {
		http.Error(w, "Bu eseri silme yetkiniz yok", http.StatusForbidden)
		return
	}

	// 4. ADIM: Önce Favorites tablosundan sil
	db.Exec("DELETE FROM Favorites WHERE ArtworkId = @p1", req.ArtworkID)

	// 5. ADIM: Sonra ArtworkPurchases tablosundan sil
	db.Exec("DELETE FROM ArtworkPurchases WHERE ArtworkId = @p1", req.ArtworkID)

	// 6. ADIM: Son olarak Artworks tablosundan sil
	_, err = db.Exec("DELETE FROM Artworks WHERE Id = @p1", req.ArtworkID)
	if err != nil {
		http.Error(w, "Eser silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Eser başarıyla silindi!"})
}

func confirmArtworksSaleHandler(w http.ResponseWriter, r *http.Request) {
	// 1. CORS Ayarları (Diğer handler'lardaki gibi sabit tutuyoruz)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Sadece POST metodu desteklenir", http.StatusMethodNotAllowed)
		return
	}

	// 2. Gelen Request Body'yi Oku
	var req struct {
		PurchaseId int `json:"purchaseId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz istek gövdesi! ❌"})
		return
	}

	// 3. Veritabanı Bağlantısı
	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "DB Bağlantı Hatası", 500)
		return
	}
	defer db.Close()

	// 4. İŞLEM: Satış Durumunu Güncelle
	// ArtworkPurchases tablosundaki Status sütununu 'Onaylandı' yapıyoruz
	query := "UPDATE ArtworkPurchases SET Status = 'Onaylandı' WHERE Id = @p1"
	_, err = db.Exec(query, req.PurchaseId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Satış onaylanırken hata oluştu: " + err.Error()})
		return
	}

	// 5. Başarılı Yanıt Dön
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Satış başarıyla onaylandı! ✅🎨"})
}

func getSellerOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// URL'den satıcı ID'sini alıyoruz (Örn: /seller-orders?sellerId=2)
	sellerId := r.URL.Query().Get("sellerId")

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// ÖNEMLİ: ArtworkPurchases ile Artworks tablolarını birleştirip eserin adını alıyoruz
	query := `
        SELECT p.Id, a.Title, p.UserEmail, p.PurchasePrice, p.Status 
        FROM ArtworkPurchases p
        JOIN Artworks a ON p.ArtworkId = a.Id
        WHERE p.SellerID = @p1`

	rows, err := db.Query(query, sellerId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var orders []map[string]interface{}
	for rows.Next() {
		var id int
		var title, email, status string
		var price float64
		rows.Scan(&id, &title, &email, &price, &status)
		orders = append(orders, map[string]interface{}{
			"id": id, "title": title, "email": email, "price": price, "status": status,
		})
	}
	json.NewEncoder(w).Encode(orders)
}

// ! SANATÇI DETAYLARI GETIRME (Kusursuz Harf ve Tip Uyumu)
func getArtistHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	artistName := r.URL.Query().Get("name")

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var artistID int
	var name, biography string
	var artworksCount, workshopsCount int

	// 1. ADIM: Sanatçı kimlik ve istatistik bilgilerini çekiyoruz (Hata payını sıfırlamak için ArtistID eklendi)
	query := `
		SELECT 
			a.ArtistID,
			ISNULL(a.ArtistName, ''), 
			ISNULL(a.Biography, ''), 
			COUNT(DISTINCT aw.Id) as ArtworkCount,
			COUNT(DISTINCT w.Id) as WorkshopCount
		FROM Artists a
		LEFT JOIN Artworks aw ON a.ArtistID = aw.ArtistID
		LEFT JOIN Users u ON a.UserID = u.UserID
		LEFT JOIN Workshops w ON u.UserID = w.InstructorID
		WHERE a.ArtistName = @p1
		GROUP BY a.ArtistID, a.ArtistName, a.Biography
	`
	err = db.QueryRow(query, artistName).Scan(&artistID, &name, &biography, &artworksCount, &workshopsCount)

	if err != nil {
		http.Error(w, "Sanatçı bulunamadı", http.StatusNotFound)
		return
	}

	// 2. ADIM: Sanatçıya ait eserleri, harf/boşluk çakışması yaşamadan direkt netleşen artistID ile çekiyoruz
	rows, err := db.Query(`
		SELECT aw.Id, aw.Title, aw.Price, aw.ImageUrl, ISNULL(aw.Category, ''),
		       ISNULL(aw.IsCampaign, 0), ISNULL(aw.DiscountRate, 0), ISNULL(aw.IsSold, 0)
		FROM Artworks aw
		WHERE aw.ArtistID = @p1
	`, artistID)

	var artworksList []map[string]interface{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, discountRate int
			var price float64
			var title, imageUrl, category string
			// SQL Server BIT tipi (true/false) için Go tarafında bool değişkenler kullanıyoruz
			var isCampaign, isSold bool

			err := rows.Scan(&id, &title, &price, &imageUrl, &category, &isCampaign, &discountRate, &isSold)
			if err != nil {
				fmt.Println("Sanatçı Sayfası Eser Scan Hatası:", err)
				continue
			}

			// 🔥 Çift yönlü harf korumalı map yapısı (Vue ne ararsa bulacak)
			artMap := map[string]interface{}{
				"id":           id,
				"Id":           id,
				"title":        title,
				"Title":        title,
				"price":        price,
				"Price":        price,
				"imageUrl":     imageUrl,
				"ImageUrl":     imageUrl,
				"image":        imageUrl,
				"category":     category,
				"Category":     category,
				"artist":       name,
				"Artist":       name,
				"IsCampaign":   isCampaign,
				"iscampaign":   isCampaign,
				"DiscountRate": discountRate,
				"discountrate": discountRate,
				"IsSold":       isSold,
				"issold":       isSold,
			}
			artworksList = append(artworksList, artMap)
		}
	}

	if artworksList == nil {
		artworksList = []map[string]interface{}{}
	}

	// 3. ADIM: Hem istatistikleri hem de harf uyumlu eser listesini tek pakette birleştiriyoruz
	response := map[string]interface{}{
		"name":           name,
		"biography":      biography,
		"artworksCount":  artworksCount,
		"workshopsCount": workshopsCount,
		"artworks":       artworksList, // Frontend'deki v-for'u besleyecek gıcır gıcır liste!
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ! ATÖLYE EKLEME (INSTRUCTOR İÇİN)
func addWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req WorkshopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Kullanıcı e-posta bilgisi eksik", http.StatusBadRequest)
		return
	}

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// 1. ADIM: Eğitmenin UserID'sini bulalım (Email üzerinden)
	var userID int
	err := db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", req.Email).Scan(&userID)
	if err != nil {
		http.Error(w, "Eğitmen hesabı bulunamadı", http.StatusUnauthorized)
		return
	}

	// 2. ADIM: Atölyeyi Workshops tablosuna kaydet
	// Not: Veritabanında InstructorID, Users tablosundaki UserID'ye bağlıdır.
	query := `INSERT INTO Workshops (Title, Description, InstructorID, Location, Capacity, Price, ImageUrl, AvailableDates) 
              VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`

	_, err = db.Exec(query,
		req.Title,
		req.Description,
		userID,
		req.Location,
		req.Capacity,
		req.Price,
		req.ImageUrl,
		req.AvailableDates)

	if err != nil {
		http.Error(w, "Atölye oluşturulamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Atölyeniz başarıyla oluşturuldu!"})
}

// ! ATÖLYE SİLME (EĞİTMEN İÇİN)
func deleteWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Sadece DELETE metodu destekleniyor", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		WorkshopID int    `json:"workshopId"`
		Email      string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Eğer body'den okumada sorun varsa, query parametresi olarak deneyelim
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Atölye ID'si gerekli", http.StatusBadRequest)
			return
		}
		workshopID, convErr := strconv.Atoi(id)
		if convErr != nil {
			http.Error(w, "Geçersiz ID", http.StatusBadRequest)
			return
		}
		req.WorkshopID = workshopID
	}

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// 1. ADIM: Workshopı kontrol et
	var instructorID int
	err := db.QueryRow("SELECT InstructorID FROM Workshops WHERE Id = @p1", req.WorkshopID).Scan(&instructorID)
	if err != nil {
		http.Error(w, "Atölye bulunamadı", http.StatusNotFound)
		return
	}

	// 2. ADIM: Eğer email verilmişse, sahibi olup olmadığını kontrol et
	if req.Email != "" {
		var userID int
		err = db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", req.Email).Scan(&userID)
		if err != nil || userID != instructorID {
			http.Error(w, "Bu atölyeyi silme yetkiniz yok", http.StatusForbidden)
			return
		}
	}

	// 3. ADIM: Önce WorkshopEnrollments tablosundan sil
	_, err = db.Exec("DELETE FROM WorkshopEnrollments WHERE WorkshopId = @p1", req.WorkshopID)
	if err != nil {
		http.Error(w, "Enrollment silme hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. ADIM: Sonra Workshops tablosundan sil
	_, err = db.Exec("DELETE FROM Workshops WHERE Id = @p1", req.WorkshopID)
	if err != nil {
		http.Error(w, "Atölye silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Atölye başarıyla silindi!"})
}

// ! EĞİTMENİN KENDİ ATÖLYELERİNİ GETİRME
func getMyWorkshopsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parametresi gerekli", http.StatusBadRequest)
		return
	}

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// Önce kullanıcının ID'sini buluyoruz
	var userID int
	err := db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", email).Scan(&userID)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	// Bu kullanıcıya ait workshopları çekiyoruz
	rows, err := db.Query("SELECT Id, Title, Description, Location, Capacity, Price, ImageUrl, AvailableDates FROM Workshops WHERE InstructorID = @p1", userID)
	if err != nil {
		http.Error(w, "Workshoplar getirilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workshops []MyWorkshop
	for rows.Next() {
		var ws MyWorkshop
		err := rows.Scan(&ws.Id, &ws.Title, &ws.Description, &ws.Location, &ws.Capacity, &ws.Price, &ws.ImageUrl, &ws.AvailableDates)
		if err != nil {
			fmt.Println("Scan hatası:", err)
			continue
		}
		workshops = append(workshops, ws)
	}

	if workshops == nil {
		workshops = []MyWorkshop{}
	}
	json.NewEncoder(w).Encode(workshops)
}

// ! SANATÇILARI GETİRME (ADMIN İÇİN)
func getArtistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// Artists ve Users tablolarını birleştirerek isim ve biyografi bilgilerini çekiyoruz
	query := `
        SELECT a.ArtistID, a.ArtistName, a.Biography, a.Nationality, u.Email 
        FROM Artists a
        JOIN Users u ON a.UserID = u.UserID`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var artists []map[string]interface{}
	for rows.Next() {
		var id int
		var name, bio, nationality, email string
		rows.Scan(&id, &name, &bio, &nationality, &email)

		artists = append(artists, map[string]interface{}{
			"id":          id,
			"name":        name,
			"biography":   bio,
			"nationality": nationality,
			"email":       email,
		})
	}

	json.NewEncoder(w).Encode(artists)
}

// ! Bakiye Yükleme Handler'ı
func updateBalanceHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email  string  `json:"email"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Veri formatı hatalı", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "DB Hatası", 500)
		return
	}
	defer db.Close()

	// Bakiyeyi güncelle (NULL ise 0 kabul et)
	_, err = db.Exec("UPDATE Users SET Balance = ISNULL(Balance, 0) + @p1 WHERE Email = @p2", data.Amount, data.Email)
	if err != nil {
		http.Error(w, "SQL Hatası", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bakiye başarıyla güncellendi!"})
}

// ! Kupon Kontrol Handler'ı
func checkCouponHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	var discount float64
	var isActive bool
	err := db.QueryRow("SELECT DiscountAmount, IsActive FROM Coupons WHERE Code = @p1", code).Scan(&discount, &isActive)

	if err != nil || !isActive {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz veya pasif kupon!"})
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"discount": discount})
}

func ApplyRandomCampaigns(db *sql.DB) error {
	// 1. Önce tüm indirimleri tertemiz sıfırla
	_, err := db.Exec("UPDATE Artworks SET IsCampaign = 0, DiscountRate = 0")
	if err != nil {
		log.Println("Sıfırlama hatası:", err)
		return err
	}

	// 2. Sadece geçerli ID'leri çek (ID'nin boş olmadığını garanti ediyoruz)
	rows, err := db.Query("SELECT Id FROM Artworks WHERE Id IS NOT NULL")
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}

	total := len(ids)
	if total == 0 {
		log.Println("Veritabanında eser bulunamadı.")
		return nil
	}

	// 3. %40 hesapla ve listeyi karıştır
	countToDiscount := int(math.Round(float64(total) * 0.4))
	// Eğer 1-2 eser varsa en az 1 tanesine indirim yapsın diye:
	if countToDiscount == 0 && total > 0 {
		countToDiscount = 1
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })

	// 4. İndirimleri Uygula (SQL Server tip hatasını CAST ile çözüyoruz)
	for i := 0; i < countToDiscount; i++ {
		randomRate := r.Intn(41) + 10 // %10-50 arası rastgele oran

		// CAST(Id AS INT) diyerek SQL'in kafasındaki soru işaretlerini siliyoruz
		query := "UPDATE Artworks SET [DiscountRate] = @p1, [IsCampaign] = 1 WHERE CAST(Id AS INT) = @p2"

		_, err := db.Exec(query,
			sql.Named("p1", randomRate),
			sql.Named("p2", ids[i]))

		if err != nil {
			log.Printf("❌ ID %d güncellenemedi: %v", ids[i], err)
		} else {
			log.Printf("✅ BAŞARILI: ID %d için %% %d indirim tanımlandı.", ids[i], randomRate)
		}
	}
	return nil
}

// SeedCoupons veritabanı sıfırlandığında sistem kuponlarının otomatik eklenmesini sağlar
func SeedCoupons(db *sql.DB) error {
	coupons := []struct {
		Code   string
		Amount float64
	}{
		{"SANAT100", 100.00},
		{"ILK500", 500.00},
	}

	for _, c := range coupons {
		// Kupon daha önce eklenmiş mi kontrol et (Hata vermemesi için)
		var exists bool
		checkQuery := "SELECT CASE WHEN EXISTS (SELECT 1 FROM Coupons WHERE Code = @p1) THEN 1 ELSE 0 END"
		err := db.QueryRow(checkQuery, c.Code).Scan(&exists)
		if err != nil {
			return err
		}

		// Eğer kupon yoksa içeriye yapıştır
		if !exists {
			insertQuery := "INSERT INTO Coupons (Code, DiscountAmount, IsActive) VALUES (@p1, @p2, 1)"
			_, err = db.Exec(insertQuery, c.Code, c.Amount)
			if err != nil {
				log.Printf("⚠️ %s kuponu otomatik eklenirken hata: %v", c.Code, err)
				return err
			}
			log.Printf("🎫 Sistem Kuponu Tanımlandı: %s (%v TL)", c.Code, c.Amount)
		}
	}
	return nil
}

func main() {
	// 1. Veritabanı bağlantısını kur (Burada kendi connString'in olduğunu varsayıyorum)
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/artworks", getArtworksHandler)
	mux.HandleFunc("/profile", getUserProfileHandler)
	mux.HandleFunc("/profile/update", updateProfileHandler)
	mux.HandleFunc("/profile/change-password", changePasswordHandler)
	mux.HandleFunc("/workshops", getWorkshopsHandler)
	mux.HandleFunc("/favorites/add", addFavoriteHandler)
	mux.HandleFunc("/favorites/remove", deleteFavoriteHandler)
	mux.HandleFunc("/favorites/check", checkFavoriteHandler)
	mux.HandleFunc("/favorites/list", getUserFavoritesHandler)
	mux.HandleFunc("/workshops/enroll", enrollWorkshopHandler)
	mux.HandleFunc("/user-enrollments", getUserEnrollmentsHandler)
	mux.HandleFunc("/update-enrollment", updateEnrollmentHandler)
	mux.HandleFunc("/delete-enrollment", deleteEnrollmentHandler)

	mux.HandleFunc("/buy-artwork/", buyArtworkHandler)
	mux.HandleFunc("/user-purchases", getUserPurchasesHandler)
	mux.HandleFunc("/add-artwork", addArtworkHandler)
	mux.HandleFunc("/delete-artwork", deleteArtworkHandler)
	mux.HandleFunc("/artist", getArtistHandler)
	mux.HandleFunc("/add-workshop", addWorkshopHandler)
	mux.HandleFunc("/delete-workshop", deleteWorkshopHandler)
	mux.HandleFunc("/my-workshops", getMyWorkshopsHandler)
	mux.HandleFunc("/artists", getArtistsHandler)
	mux.HandleFunc("/profile/update-balance", updateBalanceHandler)
	mux.HandleFunc("/check-coupon", checkCouponHandler)
	mux.HandleFunc("/confirm-sale", confirmArtworksSaleHandler)
	mux.HandleFunc("/seller-orders", getSellerOrdersHandler)
	// --- KRİTİK DOKUNUŞ: KAMPANYALARI VE KUPONLARI BURADA TETİKLİYORUZ ---
	fmt.Println("Günlük kampanyalar veritabanına işleniyor...")
	err = ApplyRandomCampaigns(db)
	if err != nil {
		log.Printf("⚠️ Kampanya tanımlama hatası: %v", err)
	} else {
		fmt.Println("✅ Kampanyalar başarıyla sabitlendi!")
	}

	// 🎫 Go backend kuponları server açılmadan hemen önce veritabanına ekliyor
	err = SeedCoupons(db)
	if err != nil {
		log.Printf("⚠️ Kupon ekleme hatası: %v", err)
	} else {
		fmt.Println("✅ Sistem kuponları başarıyla doğrulandı/eklendi!")
	}
	// -------------------------------------------------------------------

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	})

	// En sonda server ayağa kalkıyor (Bunun altında başka kod olmamalı)
	fmt.Println("Server 8080 portunda çalışıyor...")
	err = http.ListenAndServe(":8080", finalHandler)
	if err != nil {
		log.Fatal(err)
	}
}
