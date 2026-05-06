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
	"github.com/rs/cors"
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
    Date           string  `json:"date"`
    Location       string  `json:"location"`
    Capacity       int     `json:"capacity"`
    Price          float64 `json:"price"`
    ImageUrl       string  `json:"image"`
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
            w.Date, w.Location, w.Capacity, w.Price, w.ImageUrl 
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
        var dateVal time.Time
        err := rows.Scan(&ws.Id, &ws.Title, &ws.Description, &ws.InstructorName, &dateVal, &ws.Location, &ws.Capacity, &ws.Price, &ws.ImageUrl)
        if err != nil {
            continue
        }
        ws.Date = dateVal.Format("02.01.2006 15:04")
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

	handler := cors.Default().Handler(mux)
	fmt.Println("Backend 8080 portunda çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
