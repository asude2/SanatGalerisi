package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type Artwork struct {
    Id       int     `json:"id"`
    Title    string  `json:"title"`
    Artist   string  `json:"artist"`
    Price    float64 `json:"price"`
    ImageUrl string  `json:"image"`
	Description string  `json:"description"`
    ArtistInfo  string  `json:"artistInfo"`
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
type EnrollmentRequest struct {
    Email            string `json:"email"`
    WorkshopId       int    `json:"workshopId"`
    ParticipantCount int    `json:"participantCount"`
    ReservedDate     string `json:"reservedDate"`
}
type BuyArtwork struct {
    Email            string `json:"email"`
    ArtworkId        int    `json:"artworkId"`
    Price            float64 `json:"price"`
}
type SupportTicket struct {
    Id        int    `json:"id"`
    UserEmail string `json:"userEmail"`
    Subject   string `json:"subject"`
    Message   string `json:"message"`
    Status    string `json:"status"`
    CreatedAt string `json:"createdAt"`
}
type Comment struct {
    Id           int     `json:"id"`
    UserEmail    string  `json:"userEmail"`
    TargetType   string  `json:"targetType"` // "Artwork" veya "Workshop"
    TargetId     int     `json:"targetId"`
    Content      string  `json:"content"`
    Rating       int     `json:"rating"` // 1-5 arası (nullable)
    CreatedAt    string  `json:"createdAt"`
    HelpfulCount int     `json:"helpfulCount"`
}

// ! KAYIT
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Şifreyi Hashle (Güvenlik için)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	// 2. Veritabanına Bağlan (Trusted Connection kullanıyoruz)
	connString := "sqlserver://localhost:1433?database=SanatProjesi&trusted_connection=yes&encrypt=disable"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Bağlantı hatası:", err)
	}
	defer db.Close()

	// 3. Veriyi Kaydet
	// SQL Server soru işareti yerine @p1, @p2 formatını kullanır
	query := "INSERT INTO Users (FirstName, LastName, Email, Password) VALUES (@p1, @p2, @p3, @p4)"
	_, err = db.Exec(query, u.FirstName, u.LastName, u.Email, string(hashedPassword))

	if err != nil {
		fmt.Println("DB Yazma Hatası:", err)
		http.Error(w, "Veritabanına kaydedilemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kayıt başarıyla tamamlandı!"})
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

    // 1. Veritabanı Bağlantısı ve Veri Çekme
    connString := "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;"
    db, err := sql.Open("sqlserver", connString)
    if err != nil {
        http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var dbPassword string
    var firstName string // İsmi tutmak için değişken oluşturduk

    // Sorguyu güncelledik: Hem Password hem FirstName alıyoruz
    err = db.QueryRow("SELECT Password, FirstName FROM Users WHERE Email = @p1", creds.Email).Scan(&dbPassword, &firstName)
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

    // 3. Token Oluşturma
    // jwtKey'i fonksiyon içinde değil, genellikle fonksiyon dışında tanımlamak daha iyidir
    var jwtKey = []byte("cok_gizli_anahtar_123") 

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email":     creds.Email,
        "firstName": firstName, // Artık 'firstName' değişkenini kullanıyoruz
        "exp":       time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
        return
    }

    // 4. Tek bir yanıt gönderiyoruz
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Giriş başarılı!",
        "token":   tokenString,
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

    rows, err := db.Query(`
        SELECT 
            aw.Id, aw.Title, aw.Artist, aw.Price, aw.ImageUrl, 
            ISNULL(aw.Description, ''), 
            ISNULL(ar.Biography, '') 
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
        err := rows.Scan(&a.Id, &a.Title, &a.Artist, &a.Price, &a.ImageUrl, &a.Description, &a.ArtistInfo)
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
    err = db.QueryRow("SELECT FirstName, LastName, Email FROM Users WHERE Email = @p1", email).Scan(&u.FirstName, &u.LastName, &u.Email)
    
    if err != nil {
        http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
        return
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
            w.Location, w.Capacity, w.Price, w.ImageUrl 
        FROM Workshops w
        JOIN Users u ON w.InstructorID = u.UserID
    `
    rows, err := db.Query(query)
    if err != nil {
        http.Error(w, "Sorgu hatası", http.StatusInternalServerError)
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
            &ws.AvailableDates, // SQL'deki w.AvailableDates buraya gelir
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

    w.Header().Set("Content-Type", "application/json")
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
        SELECT aw.Id, aw.Title, aw.Artist, aw.Price, aw.ImageUrl 
        FROM Favorites f
        JOIN Artworks aw ON f.ArtworkId = aw.Id
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
        err := rows.Scan(&a.Id, &a.Title, &a.Artist, &a.Price, &a.ImageUrl)
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
    if r.Method == "OPTIONS" { return }

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

// ==========================================
// ! DESTEK (SUPPORT) MODÜLÜ
// ==========================================
func createSupportTicketHandler(w http.ResponseWriter, r *http.Request) {
    var req SupportTicket
    json.NewDecoder(r.Body).Decode(&req)
    db, _ := sql.Open("sqlserver", "server=localhost\\SQLEXPRESS;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
    defer db.Close()
    query := "INSERT INTO SupportTickets (UserEmail, Subject, Message) VALUES (@p1, @p2, @p3)"
    _, err := db.Exec(query, req.UserEmail, req.Subject, req.Message)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Destek talebi oluşturuldu."})
}

func getSupportTicketsHandler(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    db, _ := sql.Open("sqlserver", "server=localhost\\SQLEXPRESS;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
    defer db.Close()
    rows, _ := db.Query("SELECT Id, UserEmail, Subject, Message, Status, CreatedAt FROM SupportTickets WHERE UserEmail = @p1 ORDER BY CreatedAt DESC", email)
    defer rows.Close()
    var tickets []SupportTicket
    for rows.Next() {
        var t SupportTicket
        rows.Scan(&t.Id, &t.UserEmail, &t.Subject, &t.Message, &t.Status, &t.CreatedAt)
        tickets = append(tickets, t)
    }
    json.NewEncoder(w).Encode(tickets)
}

// ==========================================
// ! YORUMLAR (COMMENTS) MODÜLÜ
// ==========================================
func addCommentHandler(w http.ResponseWriter, r *http.Request) {
    var req Comment
    json.NewDecoder(r.Body).Decode(&req)
    db, _ := sql.Open("sqlserver", "server=localhost\\SQLEXPRESS;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
    defer db.Close()

    var count int
    if req.TargetType == "Artwork" {
        db.QueryRow("SELECT COUNT(*) FROM ArtworkPurchases WHERE UserEmail = @p1 AND ArtworkId = @p2", req.UserEmail, req.TargetId).Scan(&count)
    } else if req.TargetType == "Workshop" {
        db.QueryRow("SELECT COUNT(*) FROM WorkshopEnrollments WHERE UserEmail = @p1 AND WorkshopId = @p2", req.UserEmail, req.TargetId).Scan(&count)
    }

    if count == 0 {
        http.Error(w, "Yorum yapmak için eseri satın almış veya etkinliğe katılmış olmalısınız.", http.StatusForbidden)
        return
    }

    query := "INSERT INTO Comments (UserEmail, TargetType, TargetId, Content, Rating) VALUES (@p1, @p2, @p3, @p4, @p5)"
    _, err := db.Exec(query, req.UserEmail, req.TargetType, req.TargetId, req.Content, req.Rating)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Yorumunuz eklendi."})
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
    targetType := r.URL.Query().Get("type")
    targetId := r.URL.Query().Get("id")
    db, _ := sql.Open("sqlserver", "server=localhost\\SQLEXPRESS;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
    defer db.Close()

    query := `
        SELECT c.Id, c.UserEmail, c.TargetType, c.TargetId, c.Content, ISNULL(c.Rating, 0), c.CreatedAt, 
               (SELECT COUNT(*) FROM CommentHelpfulVotes v WHERE v.CommentId = c.Id) as HelpfulCount
        FROM Comments c
        WHERE c.TargetType = @p1 AND c.TargetId = @p2
        ORDER BY c.CreatedAt DESC`
    rows, _ := db.Query(query, targetType, targetId)
    defer rows.Close()
    
    var comments []Comment
    for rows.Next() {
        var c Comment
        rows.Scan(&c.Id, &c.UserEmail, &c.TargetType, &c.TargetId, &c.Content, &c.Rating, &c.CreatedAt, &c.HelpfulCount)
        comments = append(comments, c)
    }
    json.NewEncoder(w).Encode(comments)
}

func rateCommentHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        CommentId int `json:"commentId"`
        UserEmail string `json:"userEmail"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    db, _ := sql.Open("sqlserver", "server=localhost\\SQLEXPRESS;database=SanatProjesi;trusted_connection=yes;encrypt=disable;")
    defer db.Close()

    query := "INSERT INTO CommentHelpfulVotes (CommentId, UserEmail) VALUES (@p1, @p2)"
    _, err := db.Exec(query, req.CommentId, req.UserEmail)
    if err != nil {
        http.Error(w, "Zaten oy verdiniz", 409)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Faydalı olarak işaretlendi."})
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

    // Yeni endpoint kayıtları
    mux.HandleFunc("/support/create", createSupportTicketHandler)
    mux.HandleFunc("/support/tickets", getSupportTicketsHandler)
    mux.HandleFunc("/comments/add", addCommentHandler)
    mux.HandleFunc("/comments/list", getCommentsHandler)
    mux.HandleFunc("/comments/helpful", rateCommentHandler)

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
