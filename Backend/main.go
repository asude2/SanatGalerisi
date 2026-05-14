package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserRole  string `json:"role"` // 'User' veya 'Instructor'
	Biography string `json:"biography"`
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
	Rating       float64 `json:"rating"`
	CommentCount int     `json:"commentCount"`
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

// ! ESERLERİ GETİRME
func getArtworksHandler(w http.ResponseWriter, r *http.Request) {
	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı bağlantı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Hem eser bilgilerini hem de dinamik puan/yorum sayılarını getiriyoruz
	rows, err := db.Query(`
        SELECT 
            aw.Id, aw.Title, aw.ArtistID, aw.Price, aw.ImageUrl,
            ISNULL(aw.Description, ''),
            ISNULL(ar.FullName, 'Bilinmeyen Sanatçı'),
            ISNULL((SELECT AVG(CAST(Rating AS FLOAT)) FROM Comments WHERE TargetType = 'Artwork' AND TargetId = aw.Id), 0),
            ISNULL((SELECT COUNT(*) FROM Comments WHERE TargetType = 'Artwork' AND TargetId = aw.Id), 0)
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
		err := rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.Price, &a.ImageUrl, &a.Description, &a.Artist, &a.Rating, &a.CommentCount)
		if err != nil {
			fmt.Println("Scan hatası:", err)
			continue
		}
		artworks = append(artworks, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artworks)
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

	// Users tablosundan profil bilgilerini al
	query := `
		SELECT u.FirstName, u.LastName, u.Email, ISNULL(u.UserRole, 'User'),
			   ISNULL((SELECT TOP 1 a.Biography FROM Artists a WHERE LTRIM(RTRIM(a.FullName)) = LTRIM(RTRIM(u.FirstName + ' ' + u.LastName))), '')
		FROM Users u
		WHERE u.Email = @p1
	`
	var role string
	err = db.QueryRow(query, email).Scan(&u.FirstName, &u.LastName, &u.Email, &role, &biography)
	u.UserRole = role

	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}

	if biography.Valid {
		u.Biography = biography.String
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

// ! FAVORİ LİSTESİNİ GETİRME
func getUserFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
        SELECT aw.Id, aw.Title, aw.ArtistID, aw.Price, aw.ImageUrl, 
               ISNULL(ar.FullName, 'Bilinmeyen Sanatçı')
        FROM Favorites f
        JOIN Artworks aw ON f.ArtworkId = aw.Id
        LEFT JOIN Artists ar ON aw.ArtistID = ar.ArtistID
        WHERE f.UserEmail = @p1
    `
	rows, err := db.Query(query, email)
	if err != nil {
		http.Error(w, "Sorgu hatası", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var favs []Artwork
	for rows.Next() {
		var a Artwork
		err := rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.Price, &a.ImageUrl, &a.Artist)
		if err != nil {
			continue
		}
		favs = append(favs, a)
	}

	w.Header().Set("Content-Type", "application/json")
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

	// Önce bu kayıtla ilgili yorumları temizle (Kullanıcının bu workshop için yaptığı yorumlar)
	// Enrollment'tan workshopId ve kullanıcı email'ini almamız lazım.
	var workshopId int
	var userEmail string
	err = db.QueryRow("SELECT WorkshopId, UserEmail FROM WorkshopEnrollments WHERE Id = @p1", id).Scan(&workshopId, &userEmail)
	if err == nil {
		// Puanları/Yorumları temizle (Cascade manuel)
		db.Exec(`DELETE FROM CommentHelpfulVotes WHERE CommentId IN (SELECT Id FROM Comments WHERE TargetType = 'Workshop' AND TargetId = @p1 AND UserEmail = @p2)`, workshopId, userEmail)
		db.Exec(`DELETE FROM CommentReplies WHERE CommentId IN (SELECT Id FROM Comments WHERE TargetType = 'Workshop' AND TargetId = @p1 AND UserEmail = @p2)`, workshopId, userEmail)
		db.Exec("DELETE FROM Comments WHERE TargetType = 'Workshop' AND TargetId = @p1 AND UserEmail = @p2", workshopId, userEmail)
	}

	query := "DELETE FROM WorkshopEnrollments WHERE Id = @p1"
	_, err = db.Exec(query, id)

	if err != nil {
		http.Error(w, "Silme hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon başarıyla iptal edildi"})
}

// ! ESER SATIN ALMA
func buyArtworkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	var req BuyArtwork
	json.NewDecoder(r.Body).Decode(&req)

	db, _ := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	defer db.Close()

	// --- Eser zaten satılmış mı? ---
	var exists int
	checkQuery := "SELECT COUNT(*) FROM ArtworkPurchases WHERE ArtworkId = @p1"
	db.QueryRow(checkQuery, req.ArtworkId).Scan(&exists)

	if exists > 0 {
		// 409 Conflict: "Bu eser zaten satılmış/alınmış"
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Bu eser daha önce satın alınmış! ❌"})
		return
	}

	// satılmadıysa devam et
	query := `INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, PurchasePrice) VALUES (@p1, @p2, @p3)`
	_, err := db.Exec(query, req.Email, req.ArtworkId, req.Price)

	if err != nil {
		http.Error(w, "DB Hatası: "+err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Eser başarıyla satın alındı! 🎉"})
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
	query := `INSERT INTO Artworks (Title, ArtistID, Price, ImageURL, Description) 
              VALUES (@p1, @p2, @p3, @p4, @p5)`

	_, err = db.Exec(query, req.Title, artistID, req.Price, req.ImageUrl, req.Description)

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

	// 5.5 ADIM: Yorumları temizle (Artwork'e ait tüm yorumlar silinir)
	db.Exec(`DELETE FROM CommentHelpfulVotes WHERE CommentId IN (SELECT Id FROM Comments WHERE TargetType = 'Artwork' AND TargetId = @p1)`, req.ArtworkID)
	db.Exec(`DELETE FROM CommentReplies WHERE CommentId IN (SELECT Id FROM Comments WHERE TargetType = 'Artwork' AND TargetId = @p1)`, req.ArtworkID)
	db.Exec("DELETE FROM Comments WHERE TargetType = 'Artwork' AND TargetId = @p1", req.ArtworkID)

	// 6. ADIM: Son olarak Artworks tablosundan sil
	_, err = db.Exec("DELETE FROM Artworks WHERE Id = @p1", req.ArtworkID)
	if err != nil {
		http.Error(w, "Eser silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Eser başarıyla silindi!"})
}

// ! SANATÇI DETAYLARI GETIRME
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

	var artistInfo struct {
		Name           string `json:"name"`
		Biography      string `json:"biography"`
		ArtworksCount  int    `json:"artworksCount"`
		WorkshopsCount int    `json:"workshopsCount"`
	}

	// Sanatçı bilgisini ve eserlerin sayısını ve atölyelerin sayısını getir
	query := `
		SELECT 
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
	err = db.QueryRow(query, artistName).Scan(&artistInfo.Name, &artistInfo.Biography, &artistInfo.ArtworksCount, &artistInfo.WorkshopsCount)

	if err != nil {
		http.Error(w, "Sanatçı bulunamadı", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artistInfo)
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

	// 3.5 ADIM: Yorumları temizle (Atölyeye ait tüm yorumlar silinir)
	db.Exec(`DELETE FROM CommentHelpfulVotes WHERE CommentId IN (SELECT Id FROM Comments WHERE TargetType = 'Workshop' AND TargetId = @p1)`, req.WorkshopID)
	db.Exec(`DELETE FROM CommentReplies WHERE CommentId IN (SELECT Id FROM Comments WHERE TargetType = 'Workshop' AND TargetId = @p1)`, req.WorkshopID)
	db.Exec("DELETE FROM Comments WHERE TargetType = 'Workshop' AND TargetId = @p1", req.WorkshopID)

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

// ==========================================
// ! DESTEK (SUPPORT) MODÜLÜ
// ==========================================
type SupportTicket struct {
	Id        int    `json:"id"`
	UserEmail string `json:"userEmail"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

func createSupportTicketHandler(w http.ResponseWriter, r *http.Request) {
	var req SupportTicket
	json.NewDecoder(r.Body).Decode(&req)
	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "Veritabanı bağlantı hatası", 500)
		return
	}
	defer db.Close()
	query := "INSERT INTO SupportTickets (UserEmail, Subject, Message) VALUES (@p1, @p2, @p3)"
	_, err = db.Exec(query, req.UserEmail, req.Subject, req.Message)
	if err != nil {
		http.Error(w, "Kayıt başarısız: "+err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Destek talebi oluşturuldu!"})
}

func getSupportTicketsHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "Veritabanı hatası", 500)
		return
	}
	defer db.Close()
	rows, err := db.Query(`
		SELECT Id, UserEmail, Subject, Message, ISNULL(Status, 'Beklemede'), CONVERT(VARCHAR, CreatedAt, 120)
		FROM SupportTickets WHERE UserEmail = @p1 ORDER BY CreatedAt DESC
	`, email)
	if err != nil {
		http.Error(w, "Sorgu hatası", 500)
		return
	}
	defer rows.Close()
	var tickets []SupportTicket
	for rows.Next() {
		var t SupportTicket
		rows.Scan(&t.Id, &t.UserEmail, &t.Subject, &t.Message, &t.Status, &t.CreatedAt)
		tickets = append(tickets, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

// ==========================================
// ! YORUMLAR (COMMENTS) MODÜLÜ
// ==========================================
type Comment struct {
	Id           int    `json:"id"`
	TargetType   string `json:"targetType"`
	TargetId     int    `json:"targetId"`
	UserEmail    string `json:"userEmail"`
	UserName     string `json:"userName"`
	Content      string `json:"content"`
	HelpfulCount int    `json:"helpfulCount"`
	CreatedAt    string `json:"createdAt"`
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req Comment
	json.NewDecoder(r.Body).Decode(&req)
	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "Veritabanı hatası", 500)
		return
	}
	defer db.Close()
	query := "INSERT INTO Comments (TargetType, TargetId, UserEmail, UserName, Content) VALUES (@p1, @p2, @p3, @p4, @p5)"
	_, err = db.Exec(query, req.TargetType, req.TargetId, req.UserEmail, req.UserName, req.Content)
	if err != nil {
		http.Error(w, "Yorum eklenemedi: "+err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Yorum eklendi!"})
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("targetType")
	targetId := r.URL.Query().Get("targetId")
	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "Veritabanı hatası", 500)
		return
	}
	defer db.Close()
	rows, err := db.Query(`
		SELECT Id, TargetType, TargetId, UserEmail, UserName, Content, 
		       ISNULL(HelpfulCount, 0), CONVERT(VARCHAR, CreatedAt, 120)
		FROM Comments WHERE TargetType = @p1 AND TargetId = @p2
		ORDER BY CreatedAt DESC
	`, targetType, targetId)
	if err != nil {
		http.Error(w, "Sorgu hatası", 500)
		return
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var c Comment
		rows.Scan(&c.Id, &c.TargetType, &c.TargetId, &c.UserEmail, &c.UserName, &c.Content, &c.HelpfulCount, &c.CreatedAt)
		comments = append(comments, c)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func getAverageRatingHandler(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("targetType")
	targetId := r.URL.Query().Get("targetId")

	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "Veritabanı hatası", 500)
		return
	}
	defer db.Close()

	var average float64
	var count int

	query := `SELECT ISNULL(AVG(CAST(Rating AS FLOAT)), 0), COUNT(*) 
              FROM Comments WHERE TargetType = @p1 AND TargetId = @p2`
	err = db.QueryRow(query, targetType, targetId).Scan(&average, &count)
	if err != nil {
		average = 0
		count = 0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"average": average,
		"count":   count,
	})
}

func rateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CommentId int    `json:"commentId"`
		UserEmail string `json:"userEmail"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, err := sql.Open("sqlserver", "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
	if err != nil {
		http.Error(w, "Veritabanı hatası", 500)
		return
	}
	defer db.Close()

	// Daha önce bu yorumu faydalı işaretlemiş mi?
	var count int
	db.QueryRow("SELECT COUNT(*) FROM CommentHelpfulVotes WHERE CommentId = @p1 AND UserEmail = @p2",
		req.CommentId, req.UserEmail).Scan(&count)

	if count > 0 {
		// Zaten işaretlemiş → geri al (toggle)
		db.Exec("DELETE FROM CommentHelpfulVotes WHERE CommentId = @p1 AND UserEmail = @p2", req.CommentId, req.UserEmail)
		db.Exec("UPDATE Comments SET HelpfulCount = HelpfulCount - 1 WHERE Id = @p1", req.CommentId)
		json.NewEncoder(w).Encode(map[string]string{"message": "faydalı işareti kaldırıldı"})
	} else {
		// Henüz işaretlememiş → ekle
		db.Exec("INSERT INTO CommentHelpfulVotes (CommentId, UserEmail) VALUES (@p1, @p2)", req.CommentId, req.UserEmail)
		db.Exec("UPDATE Comments SET HelpfulCount = HelpfulCount + 1 WHERE Id = @p1", req.CommentId)
		json.NewEncoder(w).Encode(map[string]string{"message": "faydalı işaretlendi"})
	}
}

func main() {
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
	mux.HandleFunc("/artworks/buy", buyArtworkHandler)
	mux.HandleFunc("/user-purchases", getUserPurchasesHandler)
	mux.HandleFunc("/add-artwork", addArtworkHandler)
	mux.HandleFunc("/delete-artwork", deleteArtworkHandler)
	mux.HandleFunc("/artist", getArtistHandler)
	mux.HandleFunc("/add-workshop", addWorkshopHandler)
	mux.HandleFunc("/delete-workshop", deleteWorkshopHandler)
	mux.HandleFunc("/my-workshops", getMyWorkshopsHandler)
	mux.HandleFunc("/artists", getArtistsHandler)
	mux.HandleFunc("/support/create", createSupportTicketHandler)
	mux.HandleFunc("/support/tickets", getSupportTicketsHandler)
	mux.HandleFunc("/comments/add", addCommentHandler)
	mux.HandleFunc("/comments/list", getCommentsHandler)
	mux.HandleFunc("/comments/helpful", rateCommentHandler)
	mux.HandleFunc("/comments/average", getAverageRatingHandler)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Güvenlik için daha sonra frontend adresini yazabilirsin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Eğer tarayıcı "izin var mı?" (OPTIONS) diye soruyorsa, direkt OK de ve bitir.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Değilse normal akışa devam et
		mux.ServeHTTP(w, r)
	})

	fmt.Println("Server 8080 portunda çalışıyor...")
	// ListenAndServe içine 'mux' yerine 'finalHandler' yazıyoruz!
	err := http.ListenAndServe(":8080", finalHandler)
	if err != nil {
		log.Fatal(err)
	}
}
