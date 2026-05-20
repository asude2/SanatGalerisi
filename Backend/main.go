package main

import (
	"context"
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

var jwtKey = []byte("cok_gizli_anahtar_123")
var connString = "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;TrustServerCertificate=true;"

// --- MODELLER ---

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

type LoginResponse struct {
	Message   string `json:"message"`
	Token     string `json:"token"`
	Role      string `json:"role"`
	UserID    int    `json:"userId"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
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
	Rating       float64 `json:"rating"` // Enes'in eklediği alan korundu
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
	Rating         float64 `json:"rating"` // Enes'in eklediği alan korundu
}

type SupportTicket struct {
	TicketID    int       `json:"ticketId"`
	UserID      int       `json:"userId"`
	Subject     string    `json:"subject"`
	Message     string    `json:"message"`
	SupportType string    `json:"supportType"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type SupportMessage struct {
	MessageID  int       `json:"messageId"`
	TicketID   int       `json:"ticketId"`
	SenderID   int       `json:"senderId"`
	SenderName string    `json:"senderName"`
	SenderRole string    `json:"senderRole"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Comment struct {
	CommentID   int       `json:"commentId"`
	UserID      int       `json:"userId"`
	UserName    string    `json:"userName"`
	TargetID    int       `json:"targetId"`
	TargetType  string    `json:"targetType"`
	CommentText string    `json:"commentText"`
	Rating      int       `json:"rating"`
	Upvotes     int       `json:"upvotes"`
	Downvotes   int       `json:"downvotes"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"createdAt"`
	AdminReply  *string   `json:"adminReply"`
	ReplierName *string   `json:"replierName"`
	ReplierRole *string   `json:"replierRole"`
}

type InteractionLog struct {
	UserID          *int   `json:"userId"`
	TargetID        int    `json:"targetId"`
	TargetType      string `json:"targetType"`
	InteractionType string `json:"interactionType"`
}

type Comparison struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userId"`
	Title      string `json:"title"`
	TargetType string `json:"targetType"`
	TargetIDs  string `json:"targetIds"`
	CreatedAt  string `json:"createdAt"`
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

// --- MIDDLEWARE ---

type contextKey string

const (
	userEmailKey contextKey = "userEmail"
	userRoleKey  contextKey = "userRole"
	userIDKey    contextKey = "userId"
)

func isAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || authHeader == "Bearer" || authHeader == "Bearer null" || authHeader == "Bearer undefined" {
			// 🔥 Senin yazdığın o efsanevi kurtarma zırhı: Ön yüzü çökertmemek için boş dizi fırlatıyoruz kanka!
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
			return
		}

		email := claims["email"].(string)
		role := claims["role"].(string)

		db, _ := sql.Open("sqlserver", connString)
		defer db.Close()
		var userID int
		db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", email).Scan(&userID)

		ctx := context.WithValue(r.Context(), userEmailKey, email)
		ctx = context.WithValue(ctx, userRoleKey, role)
		ctx = context.WithValue(ctx, userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func isAdmin(next http.HandlerFunc) http.HandlerFunc {
	return isAuth(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(userRoleKey)
		if val == nil {
			http.Error(w, "Yetkisiz erişim", http.StatusForbidden)
			return
		}

		role := val.(string)

		if role != "Admin" && role != "Instructor" {
			http.Error(w, "Yetkisiz erişim", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- YARDIMCI FONKSİYONLAR ---

func checkEnrollment(userID int, workshopID int) bool {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var count int
	db.QueryRow("SELECT COUNT(*) FROM WorkshopEnrollments e JOIN Users u ON e.UserEmail = u.Email WHERE u.UserID = @p1 AND e.WorkshopId = @p2", userID, workshopID).Scan(&count)
	return count > 0
}

func checkPurchase(userID int, artworkID int) bool {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var count int
	db.QueryRow("SELECT COUNT(*) FROM ArtworkPurchases p JOIN Users u ON p.UserEmail = u.Email WHERE u.UserID = @p1 AND p.ArtworkId = @p2", userID, artworkID).Scan(&count)
	return count > 0
}

// --- HANDLERS ---

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	query := "INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) VALUES (@p1, @p2, @p3, @p4, @p5)"
	_, err := db.Exec(query, u.FirstName, u.LastName, u.Email, string(hashedPassword), u.UserRole)
	if err != nil {
		http.Error(w, "Kayıt hatası", 500)
		return
	}
	w.WriteHeader(201)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var dbPassword, firstName, role string
	var userID int
	err := db.QueryRow("SELECT UserID, Password, FirstName, UserRole FROM Users WHERE Email = @p1", creds.Email).Scan(&userID, &dbPassword, &firstName, &role)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(creds.Password)) != nil {
		http.Error(w, "Hatalı giriş", 401)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": creds.Email, "role": role, "exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString, "role": role, "firstName": firstName, "userId": userID, "email": creds.Email,
	})
}

func getArtworksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı bağlantı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
        SELECT 
            aw.Id, aw.Title, aw.ArtistID, aw.Price, aw.ImageUrl,
            ISNULL(aw.Description, ''), ISNULL(aw.Category, ''),
            ISNULL(ar.ArtistName, 'Bilinmeyen Sanatçı'),
            ISNULL(aw.IsCampaign, 0), ISNULL(aw.DiscountRate, 0), ISNULL(aw.IsSold, 0)
        FROM Artworks aw
        LEFT JOIN Artists ar ON aw.ArtistID = ar.ArtistID`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Sorgu yürütülemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var artworks []Artwork
	for rows.Next() {
		var a Artwork
		err := rows.Scan(
			&a.ID, &a.Title, &a.ArtistID, &a.Price, &a.ImageUrl,
			&a.Description, &a.Category, &a.Artist,
			&a.IsCampaign, &a.DiscountRate, &a.IsSold,
		)
		if err != nil {
			fmt.Println("Scan hatası:", err)
			continue
		}
		artworks = append(artworks, a)
	}

	json.NewEncoder(w).Encode(artworks)
}

// 🛍️ %100 PARAMETRE UYUMLU VE SELLERID İLİŞKİLİ NİHAİ ESER SATIN ALMA HANDLERI
func buyArtworkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Hem senin dinamik slug yapını hem Enes'in düz rotasını destekleyen hibrit ID çekici kanka
	var idStr string
	if strings.Contains(r.URL.Path, "/buy-artwork/") {
		idStr = strings.TrimPrefix(r.URL.Path, "/buy-artwork/")
	} else {
		idStr = strings.TrimPrefix(r.URL.Path, "/artworks/buy/")
	}
	idStr = strings.Trim(idStr, "/")

	var data struct {
		Email         string  `json:"email"`
		ArtworkId     int     `json:"artworkId"`
		Price         float64 `json:"price"`
		PaymentMethod string  `json:"paymentMethod"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz veri formatı! ❌"})
		return
	}

	artworkID := data.ArtworkId
	if artworkID <= 0 && idStr != "" {
		artworkID, _ = strconv.Atoi(idStr)
	}

	if artworkID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz Eser ID formatı! ❌"})
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "DB Bağlantı Hatası"})
		return
	}
	defer db.Close()

	var artworkPrice float64
	var category string
	var sellerID int
	err = db.QueryRow("SELECT Price, Category, ArtistID FROM Artworks WHERE Id = @p1", artworkID).Scan(&artworkPrice, &category, &sellerID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Eser bulunamadı veya silinmiş! 🔍"})
		return
	}

	pMethod := data.PaymentMethod
	if pMethod == "" {
		pMethod = "Kredi Kartı"
	}
	finalPurchasePrice := data.Price
	if finalPurchasePrice <= 0 {
		finalPurchasePrice = artworkPrice
	}

	// Mükerrer satışı önleyen Enes'in COUNT kontrol zırhı entegre edildi kanka
	var exists int
	db.QueryRow("SELECT COUNT(*) FROM ArtworkPurchases WHERE ArtworkId = @p1", artworkID).Scan(&exists)
	if exists > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Eser zaten satılmış! 💸"})
		return
	}

	if pMethod == "Uygulama Bakiyesi" {
		var userBalance float64
		err = db.QueryRow("SELECT Balance FROM Users WHERE Email = @p1", data.Email).Scan(&userBalance)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Kullanıcı bulunamadı."})
			return
		}

		if userBalance < finalPurchasePrice {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Bakiyeniz yetersiz! 💸"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE Users SET Balance = Balance - @p1 WHERE Email = @p2", finalPurchasePrice, data.Email)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE Artworks SET IsSold = 1 WHERE Id = @p1", artworkID)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, SellerID, PurchasePrice, Status, PaymentMethod, CreatedAt) 
                  VALUES (@p1, @p2, @p3, @p4, @p5, @p6, GETDATE())`

		_, err = tx.Exec(query,
			sql.Named("p1", data.Email),
			sql.Named("p2", artworkID),
			sql.Named("p3", sellerID),
			sql.Named("p4", finalPurchasePrice),
			sql.Named("p5", "Hazırlanıyor"),
			sql.Named("p6", pMethod),
		)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Sipariş kaydedilirken hata oluştu: " + err.Error()})
			return
		}
		tx.Commit()

	} else {
		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE Artworks SET IsSold = 1 WHERE Id = @p1", artworkID)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, SellerID, PurchasePrice, Status, PaymentMethod, CreatedAt) 
                  VALUES (@p1, @p2, @p3, @p4, @p5, @p6, GETDATE())`

		_, err = tx.Exec(query,
			sql.Named("p1", data.Email),
			sql.Named("p2", artworkID),
			sql.Named("p3", sellerID),
			sql.Named("p4", finalPurchasePrice),
			sql.Named("p5", "Hazırlanıyor"),
			sql.Named("p6", "Kredi Kartı"),
		)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Sipariş kaydedilirken hata oluştu: " + err.Error()})
			return
		}
		tx.Commit()
	}

	if category != "" && data.Email != "" {
		db.Exec("UPDATE Users SET LastPurchasedCategory = @p1 WHERE Email = @p2", category, data.Email)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Eser başarıyla satın alındı! 🎉🎨"})
}

type WorkshopEnrollRequest struct {
	Email            string `json:"email"`
	WorkshopID       int    `json:"workshopId"`
	ParticipantCount int    `json:"participantCount"`
	ReservedDate     string `json:"reservedDate"`
	PaymentMethod    string `json:"paymentMethod"`
}

func enrollWorkshopWithPaymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req WorkshopEnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz veri formatı! ❌"})
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if req.PaymentMethod == "Uygulama Bakiyesi" {
		var workshopPrice float64
		err = db.QueryRow("SELECT Price FROM Workshops WHERE Id = @p1", req.WorkshopID).Scan(&workshopPrice)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Atölye bulunamadı! 🔍"})
			return
		}

		totalPrice := workshopPrice * float64(req.ParticipantCount)

		var userBalance float64
		err = db.QueryRow("SELECT Balance FROM Users WHERE Email = @p1", req.Email).Scan(&userBalance)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Kullanıcı hesabı bulunamadı! 👤"})
			return
		}

		if userBalance < totalPrice {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Bakiyeniz yetersiz! 💸 Lütfen cüzdanınıza bakiye yükleyin."})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE Users SET Balance = Balance - @p1 WHERE Email = @p2", totalPrice, req.Email)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO WorkshopEnrollments (UserEmail, WorkshopId, ParticipantCount, ReservedDate, Status) 
                  VALUES (@p1, @p2, @p3, @p4, 'Onay Bekliyor')`
		_, err = tx.Exec(query, req.Email, req.WorkshopID, req.ParticipantCount, req.ReservedDate)
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tx.Commit()

	} else {
		query := `INSERT INTO WorkshopEnrollments (UserEmail, WorkshopId, ParticipantCount, ReservedDate, Status) 
                  VALUES (@p1, @p2, @p3, @p4, 'Onay Bekliyor')`
		_, err = db.Exec(query, req.Email, req.WorkshopID, req.ParticipantCount, req.ReservedDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon oluşturulurken bir hata oluştu."})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Atölye rezervasyonu başarıyla oluşturuldu! 🎉🎨"})
}

func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var u User
	var biography sql.NullString
	var balance sql.NullFloat64

	query := `
		SELECT u.FirstName, u.LastName, u.Email, ISNULL(a.Biography, ''), ISNULL(u.Balance, 0), u.UserRole
		FROM Users u
		LEFT JOIN Artists a ON u.UserID = a.UserID
		WHERE u.Email = @p1
	`
	err := db.QueryRow(query, email).Scan(&u.FirstName, &u.LastName, &u.Email, &biography, &balance, &u.UserRole)

	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", 404)
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

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var userID int
	db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", u.Email).Scan(&userID)
	db.Exec("UPDATE Users SET FirstName = @p1, LastName = @p2 WHERE UserID = @p3", u.FirstName, u.LastName, userID)
	var artistID int
	err := db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", userID).Scan(&artistID)
	if err == nil {
		db.Exec("UPDATE Artists SET Biography = @p1, ArtistName = @p2 WHERE ArtistID = @p3", u.Biography, u.FirstName+" "+u.LastName, artistID)
	} else {
		db.Exec("INSERT INTO Artists (UserID, Biography, ArtistName) VALUES (@p1, @p2, @p3)", userID, u.Biography, u.FirstName+" "+u.LastName)
	}
	w.WriteHeader(200)
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var dbPassword string
	db.QueryRow("SELECT Password FROM Users WHERE Email = @p1", req.Email).Scan(&dbPassword)
	if bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(req.OldPassword)) != nil {
		http.Error(w, "Eski şifre hatalı", 401)
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	db.Exec("UPDATE Users SET Password = @p1 WHERE Email = @p2", string(newHash), req.Email)
	w.WriteHeader(200)
}

func getWorkshopsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	rows, _ := db.Query(`SELECT w.Id, w.Title, w.Description, u.FirstName + ' ' + u.LastName, w.AvailableDates, w.Location, w.Capacity, w.Price, w.ImageUrl FROM Workshops w JOIN Users u ON w.InstructorID = u.UserID`)
	defer rows.Close()
	var workshops []Workshop
	for rows.Next() {
		var ws Workshop
		rows.Scan(&ws.Id, &ws.Title, &ws.Description, &ws.InstructorName, &ws.AvailableDates, &ws.Location, &ws.Capacity, &ws.Price, &ws.ImageUrl)
		workshops = append(workshops, ws)
	}
	json.NewEncoder(w).Encode(workshops)
}

func addFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		ArtworkId int    `json:"artworkId"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("INSERT INTO Favorites (UserEmail, ArtworkId) VALUES (@p1, @p2)", req.Email, req.ArtworkId)
	w.WriteHeader(201)
}

func deleteFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		ArtworkId int    `json:"artworkId"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("DELETE FROM Favorites WHERE UserEmail = @p1 AND ArtworkId = @p2", req.Email, req.ArtworkId)
	w.WriteHeader(200)
}

func checkFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	artworkId := r.URL.Query().Get("artworkId")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var count int
	db.QueryRow("SELECT COUNT(*) FROM Favorites WHERE UserEmail = @p1 AND ArtworkId = @p2", email, artworkId).Scan(&count)
	json.NewEncoder(w).Encode(map[string]bool{"isFavorite": count > 0})
}

func getUserFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parametresi gerekli", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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

func enrollWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email            string `json:"email"`
		WorkshopId       int    `json:"workshopId"`
		ParticipantCount int    `json:"participantCount"`
		ReservedDate     string `json:"reservedDate"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("INSERT INTO WorkshopEnrollments (UserEmail, WorkshopId, ParticipantCount, ReservedDate) VALUES (@p1, @p2, @p3, @p4)", req.Email, req.WorkshopId, req.ParticipantCount, req.ReservedDate)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon başarılı!"})
}

func updateEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sadece PUT metodu destekleniyor"})
		return
	}

	var data struct {
		ID               int    `json:"id"`
		ParticipantCount int    `json:"participantCount"`
		ReservedDate     string `json:"reservedDate"`
		PaymentMethod    string `json:"paymentMethod"`
		Email            string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz veri formatı! ❌"})
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Veritabanı bağlantı hatası!"})
		return
	}
	defer db.Close()

	var oldParticipantCount, workshopID int
	err = db.QueryRow("SELECT WorkshopId, ParticipantCount FROM WorkshopEnrollments WHERE Id = @p1", data.ID).Scan(&workshopID, &oldParticipantCount)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon bulunamadı! 🔍"})
		return
	}

	var workshopPrice float64
	err = db.QueryRow("SELECT Price FROM Workshops WHERE Id = @p1", workshopID).Scan(&workshopPrice)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Atölye fiyatı bulunamadı!"})
		return
	}

	diff := data.ParticipantCount - oldParticipantCount

	if diff > 0 {
		extraPrice := workshopPrice * float64(diff)

		if data.PaymentMethod == "Uygulama Bakiyesi" {
			var userBalance float64
			err = db.QueryRow("SELECT Balance FROM Users WHERE Email = @p1", data.Email).Scan(&userBalance)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"message": "Kullanıcı hesabı bulunamadı!"})
				return
			}

			if userBalance < extraPrice {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"message": "Ekstra katılımcı ödemesi için bakiyeniz yetersiz! 💸"})
				return
			}

			tx, err := db.Begin()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = tx.Exec("UPDATE Users SET Balance = Balance - @p1 WHERE Email = @p2", extraPrice, data.Email)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			query := "UPDATE WorkshopEnrollments SET ParticipantCount = @p1, ReservedDate = @p2, Status = 'Onay Bekliyor' WHERE Id = @p3"
			_, err = tx.Exec(query, data.ParticipantCount, data.ReservedDate, data.ID)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			tx.Commit()
		} else {
			query := "UPDATE WorkshopEnrollments SET ParticipantCount = @p1, ReservedDate = @p2, Status = 'Onay Bekliyor' WHERE Id = @p3"
			_, err = db.Exec(query, data.ParticipantCount, data.ReservedDate, data.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Güncelleme hatası: " + err.Error()})
				return
			}
		}
	} else {
		query := "UPDATE WorkshopEnrollments SET ParticipantCount = @p1, ReservedDate = @p2 WHERE Id = @p3"
		_, err = db.Exec(query, data.ParticipantCount, data.ReservedDate, data.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Güncelleme hatası: " + err.Error()})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon başarıyla güncellendi"})
}

func deleteEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("DELETE FROM WorkshopEnrollments WHERE Id = @p1", id)
	w.WriteHeader(200)
}

func getUserPurchasesHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	query := `SELECT p.Id, a.Title, p.PurchasePrice, p.Status, p.CreatedAt, a.ImageUrl FROM ArtworkPurchases p JOIN Artworks a ON p.ArtworkId = a.Id WHERE p.UserEmail = @p1`
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

func addArtworkHandler(w http.ResponseWriter, r *http.Request) {
	var req ArtworkRequest
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var userID int
	var firstName, lastName string
	db.QueryRow("SELECT UserID, FirstName, LastName FROM Users WHERE Email = @p1", req.Email).Scan(&userID, &firstName, &lastName)
	var artistID int
	artistName := firstName + " " + lastName
	err := db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", userID).Scan(&artistID)
	if err != nil {
		db.QueryRow("INSERT INTO Artists (UserID, ArtistName) OUTPUT INSERTED.ArtistID VALUES (@p1, @p2)", userID, artistName).Scan(&artistID)
	}
	db.Exec("INSERT INTO Artworks (Title, ArtistID, Price, ImageURL, Description, Category) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)", req.Title, artistID, req.Price, req.ImageUrl, req.Description, req.Category)
	w.WriteHeader(201)
}

func deleteArtworkHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ArtworkID int    `json:"artworkId"`
		Email     string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	var userID int
	err := db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", req.Email).Scan(&userID)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	var artworkArtistID int
	err = db.QueryRow("SELECT ArtistID FROM Artworks WHERE Id = @p1", req.ArtworkID).Scan(&artworkArtistID)
	if err != nil {
		http.Error(w, "Eser bulunamadı", http.StatusNotFound)
		return
	}

	log.Printf("🔥 [ESER SİLME RADARI] Ön yüzden Gelen Email: %s | Bulunan UserID: %d", req.Email, userID)

	var artistID int
	err = db.QueryRow("SELECT ArtistID FROM Artists WHERE UserID = @p1", userID).Scan(&artistID)
	if err != nil || artistID != artworkArtistID {
		http.Error(w, "Bu eseri silme yetkiniz yok", http.StatusForbidden)
		return
	}

	db.Exec("DELETE FROM Favorites WHERE ArtworkId = @p1", req.ArtworkID)
	db.Exec("DELETE FROM ArtworkPurchases WHERE ArtworkId = @p1", req.ArtworkID)
	db.Exec("DELETE FROM Artworks WHERE Id = @p1", req.ArtworkID)
	w.WriteHeader(200)
}

func confirmArtworksSaleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		PurchaseId int `json:"purchaseId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz istek gövdesi! ❌"})
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "DB Bağlantı Hatası", 500)
		return
	}
	defer db.Close()

	query := "UPDATE ArtworkPurchases SET Status = 'Onaylandı' WHERE Id = @p1"
	_, err = db.Exec(query, req.PurchaseId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Satış onaylanırken hata oluştu: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Satış başarıyla onaylandı! ✅🎨"})
}

func getUserEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parametresi gerekli", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
        SELECT e.Id, w.Title, e.ParticipantCount, e.ReservedDate, 
               CONVERT(NVARCHAR, e.CreatedAt, 120) as CreatedAt, 
               w.Location, w.AvailableDates, ISNULL(e.Status, 'Onay Bekliyor'), w.Price
        FROM WorkshopEnrollments e
        JOIN Workshops w ON e.WorkshopId = w.Id
        WHERE e.UserEmail = @p1
        ORDER BY e.CreatedAt DESC`

	rows, err := db.Query(query, email)
	if err != nil {
		log.Printf("❌ user-enrollments SQL Hatası: %v", err)
		http.Error(w, "Sorgu hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var enrollments []map[string]interface{}
	for rows.Next() {
		var id, pCount int
		var price float64
		var title, rDate, cAt, location, aDates, status string
		err := rows.Scan(&id, &title, &pCount, &rDate, &cAt, &location, &aDates, &status, &price)
		if err != nil {
			log.Printf("❌ user-enrollments Scan Hatası: %v", err)
			continue
		}

		enrollments = append(enrollments, map[string]interface{}{
			"id":               id,
			"workshopTitle":    title,
			"participantCount": pCount,
			"reservedDate":     rDate,
			"createdAt":        cAt,
			"location":         location,
			"availableDates":   aDates,
			"status":           status,
			"workshopPrice":    price,
		})
	}

	if enrollments == nil {
		enrollments = []map[string]interface{}{}
	}
	json.NewEncoder(w).Encode(enrollments)
}

func getSellerWorkshopsOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	sellerEmail := r.URL.Query().Get("email")
	if sellerEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "E-posta parametresi eksik! ❌"})
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "DB Hatası", 500)
		return
	}
	defer db.Close()

	query := `
		SELECT e.Id, w.Title, e.UserEmail, e.ParticipantCount, e.ReservedDate, ISNULL(e.Status, 'Onay Bekliyor')
		FROM WorkshopEnrollments e
		JOIN Workshops w ON e.WorkshopId = w.Id
		JOIN Users u ON w.InstructorID = u.UserID
		WHERE u.Email = @p1
		ORDER BY e.CreatedAt DESC`

	rows, err := db.Query(query, sellerEmail)
	if err != nil {
		log.Printf("❌ seller-workshops-orders SQL Hatası: %v", err)
		http.Error(w, "Sorgu hatası: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var id, participantCount int
		var title, email, reservedDate, status string

		err := rows.Scan(&id, &title, &email, &participantCount, &reservedDate, &status)
		if err != nil {
			log.Printf("❌ seller-workshops-orders Scan Hatası: %v", err)
			continue
		}

		orders = append(orders, map[string]interface{}{
			"id":               id,
			"title":            title,
			"email":            email,
			"participantCount": participantCount,
			"reservedDate":     reservedDate,
			"status":           status,
		})
	}

	if orders == nil {
		orders = []map[string]interface{}{}
	}
	json.NewEncoder(w).Encode(orders)
}

func confirmWorkshopEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		EnrollmentId int `json:"enrollmentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Geçersiz istek! ❌"})
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "DB Bağlantı Hatası", 500)
		return
	}
	defer db.Close()

	query := "UPDATE WorkshopEnrollments SET Status = 'Onaylandı' WHERE Id = @p1"
	_, err = db.Exec(query, req.EnrollmentId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Rezervasyon onaylanırken hata oluştu: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Atölye rezervasyonu başarıyla onaylandı! ✅🎨"})
}

// 🛍️ %100 UYUMLU VE HİBRİT ALAN DESTEKLİ ASUDE SİPARİŞ GETİRME HANDLERI
func getSellerOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	sellerEmail := r.URL.Query().Get("email")
	if sellerEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "E-posta parametresi eksik! ❌"})
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "DB Bağlantı Hatası: "+err.Error(), 500)
		return
	}
	defer db.Close()

	query := `
		SELECT p.Id, a.Title, p.UserEmail, p.PurchasePrice, p.Status 
		FROM ArtworkPurchases p
		JOIN Artworks a ON p.ArtworkId = a.Id
		JOIN Artists art ON p.SellerID = art.ArtistID
		JOIN Users u ON art.UserID = u.UserID
		WHERE u.Email = @p1`

	rows, err := db.Query(query, sellerEmail)
	if err != nil {
		http.Error(w, "SQL Sorgu Hatası: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var id int
		var title, email, status string
		var price float64

		err := rows.Scan(&id, &title, &email, &price, &status)
		if err != nil {
			http.Error(w, "Veri Okuma Hatası: "+err.Error(), 500)
			return
		}

		orders = append(orders, map[string]interface{}{
			"id":            id,
			"purchaseId":    id,
			"title":         title,
			"artworkTitle":  title,
			"email":         email,
			"buyerEmail":    email,
			"price":         price,
			"purchasePrice": price,
			"status":        status,
		})
	}

	if orders == nil {
		orders = []map[string]interface{}{}
	}

	json.NewEncoder(w).Encode(orders)
}

func getArtistHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	artistName := r.URL.Query().Get("name")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	var artistID int
	var name, biography string
	var artworksCount, workshopsCount int

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
	err := db.QueryRow(query, artistName).Scan(&artistID, &name, &biography, &artworksCount, &workshopsCount)

	if err != nil {
		http.Error(w, "Sanatçı bulunamadı", http.StatusNotFound)
		return
	}

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
			var isCampaign, isSold bool

			err := rows.Scan(&id, &title, &price, &imageUrl, &category, &isCampaign, &discountRate, &isSold)
			if err != nil {
				fmt.Println("Sanatçı Sayfası Eser Scan Hatası:", err)
				continue
			}

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

	response := map[string]interface{}{
		"name":           name,
		"biography":      biography,
		"artworksCount":  artworksCount,
		"workshopsCount": workshopsCount,
		"artworks":       artworksList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	var req WorkshopRequest
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var userID int
	db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", req.Email).Scan(&userID)
	db.Exec("INSERT INTO Workshops (Title, Description, InstructorID, Location, Capacity, Price, ImageUrl, AvailableDates) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)", req.Title, req.Description, userID, req.Location, req.Capacity, req.Price, req.ImageUrl, req.AvailableDates)
	w.WriteHeader(201)
}

func deleteWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WorkshopID int    `json:"workshopId"`
		Email      string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		req.WorkshopID = id
	}
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("DELETE FROM WorkshopEnrollments WHERE WorkshopId = @p1", req.WorkshopID)
	db.Exec("DELETE FROM Workshops WHERE Id = @p1", req.WorkshopID)
	w.WriteHeader(200)
}

func getMyWorkshopsHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var userID int
	db.QueryRow("SELECT UserID FROM Users WHERE Email = @p1", email).Scan(&userID)
	rows, _ := db.Query("SELECT Id, Title, Description, Location, Capacity, Price, ImageUrl, AvailableDates FROM Workshops WHERE InstructorID = @p1", userID)
	defer rows.Close()
	var workshops []Workshop
	for rows.Next() {
		var ws Workshop
		rows.Scan(&ws.Id, &ws.Title, &ws.Description, &ws.Location, &ws.Capacity, &ws.Price, &ws.ImageUrl, &ws.AvailableDates)
		workshops = append(workshops, ws)
	}
	json.NewEncoder(w).Encode(workshops)
}

func createTicketHandler(w http.ResponseWriter, r *http.Request) {
	var t SupportTicket
	json.NewDecoder(r.Body).Decode(&t)
	userID := r.Context().Value(userIDKey).(int)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("INSERT INTO SupportTickets (UserID, Subject, Message, SupportType) VALUES (@p1, @p2, @p3, @p4)", userID, t.Subject, t.Message, t.SupportType)
	w.WriteHeader(201)
}

func getUserTicketsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	rows, _ := db.Query("SELECT TicketID, Subject, Message, SupportType, Status, CreatedAt, UpdatedAt FROM SupportTickets WHERE UserID = @p1 ORDER BY CreatedAt DESC", userID)
	defer rows.Close()
	var tickets []SupportTicket
	for rows.Next() {
		var t SupportTicket
		t.UserID = userID
		rows.Scan(&t.TicketID, &t.Subject, &t.Message, &t.SupportType, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		tickets = append(tickets, t)
	}
	if tickets == nil {
		tickets = []SupportTicket{}
	}
	json.NewEncoder(w).Encode(tickets)
}

func getAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	rows, _ := db.Query("SELECT TicketID, UserID, Subject, Message, SupportType, Status, CreatedAt, UpdatedAt FROM SupportTickets ORDER BY CreatedAt DESC")
	defer rows.Close()
	var tickets []SupportTicket
	for rows.Next() {
		var t SupportTicket
		rows.Scan(&t.TicketID, &t.UserID, &t.Subject, &t.Message, &t.SupportType, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		tickets = append(tickets, t)
	}
	if tickets == nil {
		tickets = []SupportTicket{}
	}
	json.NewEncoder(w).Encode(tickets)
}

func sendTicketMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg SupportMessage
	json.NewDecoder(r.Body).Decode(&msg)
	userID := r.Context().Value(userIDKey).(int)

	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	var status string
	err := db.QueryRow("SELECT Status FROM SupportTickets WHERE TicketID = @p1", msg.TicketID).Scan(&status)
	if err != nil {
		http.Error(w, "Bilet bulunamadı", 404)
		return
	}
	if status == "Closed" || status == "Çözüldü" {
		http.Error(w, "Çözülmüş bir talebe mesaj gönderilemez.", 403)
		return
	}

	db.Exec("INSERT INTO SupportMessages (TicketID, SenderID, Message) VALUES (@p1, @p2, @p3)", msg.TicketID, userID, msg.Message)
	db.Exec("UPDATE SupportTickets SET UpdatedAt = GETDATE(), Status = 'Açık' WHERE TicketID = @p1", msg.TicketID)
	w.WriteHeader(201)
}

func getTicketMessagesHandler(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("ticketId")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	query := `SELECT m.MessageID, m.TicketID, m.SenderID, u.FirstName + ' ' + u.LastName, u.UserRole, m.Message, m.CreatedAt FROM SupportMessages m JOIN Users u ON m.SenderID = u.UserID WHERE m.TicketID = @p1 ORDER BY m.CreatedAt ASC`
	rows, _ := db.Query(query, ticketID)
	defer rows.Close()
	var msgs []SupportMessage
	for rows.Next() {
		var m SupportMessage
		rows.Scan(&m.MessageID, &m.TicketID, &m.SenderID, &m.SenderName, &m.SenderRole, &m.Message, &m.CreatedAt)
		msgs = append(msgs, m)
	}
	if msgs == nil {
		msgs = []SupportMessage{}
	}
	json.NewEncoder(w).Encode(msgs)
}

func updateTicketStatusHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TicketID int    `json:"ticketId"`
		Status   string `json:"status"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("UPDATE SupportTickets SET Status = @p1, UpdatedAt = GETDATE() WHERE TicketID = @p2", req.Status, req.TicketID)
	w.WriteHeader(200)
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	var c Comment
	json.NewDecoder(r.Body).Decode(&c)
	userID := r.Context().Value(userIDKey).(int)
	isVerified := false

	switch c.TargetType {
	case "Artwork":
		isVerified = checkPurchase(userID, c.TargetID)
	case "Workshop":
		isVerified = checkEnrollment(userID, c.TargetID)
		if !isVerified {
			http.Error(w, "Bu atölyeye yorum yapabilmek için önce katılmanız gerekmektedir.", 403)
			return
		}
	}

	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("INSERT INTO Comments (UserID, TargetID, TargetType, CommentText, Rating, IsVerified) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)", userID, c.TargetID, c.TargetType, c.CommentText, c.Rating, isVerified)
	w.WriteHeader(201)
}

func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	targetID := r.URL.Query().Get("targetId")
	targetType := r.URL.Query().Get("targetType")
	sortBy := r.URL.Query().Get("sort")

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("DB Bağlantı Hatası:", err)
		json.NewEncoder(w).Encode([]Comment{})
		return
	}
	defer db.Close()

	orderClause := "ORDER BY c.CreatedAt DESC"
	switch sortBy {
	case "highest":
		orderClause = "ORDER BY c.Rating DESC, c.CreatedAt DESC"
	case "most_helpful":
		orderClause = "ORDER BY (c.Upvotes - ISNULL(c.Downvotes, 0)) DESC, c.CreatedAt DESC"
	}

	query := fmt.Sprintf(`
		SELECT 
			c.CommentID, c.UserID, u.FirstName + ' ' + u.LastName, c.CommentText, c.Rating, 
			c.Upvotes, ISNULL(c.Downvotes, 0), c.IsVerified, c.CreatedAt, 
			(SELECT TOP 1 cr.ReplyText FROM CommentReplies cr WHERE cr.CommentID = c.CommentID ORDER BY cr.CreatedAt DESC) as AdminReply,
			(SELECT TOP 1 ru.FirstName + ' ' + ru.LastName FROM CommentReplies cr JOIN Users ru ON cr.UserID = ru.UserID WHERE cr.CommentID = c.CommentID ORDER BY cr.CreatedAt DESC) as ReplierName,
			(SELECT TOP 1 ru.UserRole FROM CommentReplies cr JOIN Users ru ON cr.UserID = ru.UserID WHERE cr.CommentID = c.CommentID ORDER BY cr.CreatedAt DESC) as ReplierRole
		FROM Comments c 
		JOIN Users u ON c.UserID = u.UserID 
		WHERE c.TargetID = @p1 AND c.TargetType = @p2 %s`, orderClause)

	rows, err := db.Query(query, targetID, targetType)
	if err != nil {
		fmt.Println("Yorum Getirme Sorgu Hatası:", err)
		json.NewEncoder(w).Encode([]Comment{})
		return
	}
	if rows == nil {
		json.NewEncoder(w).Encode([]Comment{})
		return
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		var reply, replierName, replierRole sql.NullString
		err := rows.Scan(
			&c.CommentID, &c.UserID, &c.UserName, &c.CommentText, &c.Rating,
			&c.Upvotes, &c.Downvotes, &c.IsVerified, &c.CreatedAt,
			&reply, &replierName, &replierRole,
		)
		if err != nil {
			fmt.Println("Yorum Scan Hatası:", err)
			continue
		}
		if reply.Valid {
			c.AdminReply = &reply.String
		}
		if replierName.Valid {
			c.ReplierName = &replierName.String
		}
		if replierRole.Valid {
			c.ReplierRole = &replierRole.String
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func voteCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CommentID int    `json:"commentId"`
		VoteType  string `json:"voteType"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value(userIDKey).(int)

	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	var existingVote string
	err := db.QueryRow("SELECT VoteType FROM CommentVotes WHERE UserID = @p1 AND CommentID = @p2", userID, req.CommentID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		db.Exec("INSERT INTO CommentVotes (UserID, CommentID, VoteType) VALUES (@p1, @p2, @p3)", userID, req.CommentID, req.VoteType)
		if req.VoteType == "Up" {
			db.Exec("UPDATE Comments SET Upvotes = Upvotes + 1 WHERE CommentID = @p1", req.CommentID)
		} else {
			db.Exec("UPDATE Comments SET Downvotes = Downvotes + 1 WHERE CommentID = @p1", req.CommentID)
		}
	} else if existingVote == req.VoteType {
		db.Exec("DELETE FROM CommentVotes WHERE UserID = @p1 AND CommentID = @p2", userID, req.CommentID)
		if req.VoteType == "Up" {
			db.Exec("UPDATE Comments SET Upvotes = Upvotes - 1 WHERE CommentID = @p1", req.CommentID)
		} else {
			db.Exec("UPDATE Comments SET Downvotes = Downvotes - 1 WHERE CommentID = @p1", req.CommentID)
		}
	} else {
		db.Exec("UPDATE CommentVotes SET VoteType = @p1 WHERE UserID = @p2 AND CommentID = @p3", req.VoteType, userID, req.CommentID)
		if req.VoteType == "Up" {
			db.Exec("UPDATE Comments SET Upvotes = Upvotes + 1, Downvotes = Downvotes - 1 WHERE CommentID = @p1", req.CommentID)
		} else {
			db.Exec("UPDATE Comments SET Downvotes = Downvotes + 1, Upvotes = Upvotes - 1 WHERE CommentID = @p1", req.CommentID)
		}
	}
	w.WriteHeader(200)
}

func addCommentReplyHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CommentID int    `json:"commentId"`
		ReplyText string `json:"replyText"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value(userIDKey).(int)
	role := r.Context().Value(userRoleKey).(string)

	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	if role != "Admin" && role != "Instructor" {
		http.Error(w, "Yalnızca galeri yöneticileri veya eğitmenler yorumlara yanıt verebilir.", 403)
		return
	}

	db.Exec("INSERT INTO CommentReplies (CommentID, UserID, ReplyText) VALUES (@p1, @p2, @p3)", req.CommentID, userID, req.ReplyText)
	w.WriteHeader(201)
}

func logInteractionHandler(w http.ResponseWriter, r *http.Request) {
	var logReq InteractionLog
	json.NewDecoder(r.Body).Decode(&logReq)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	if logReq.InteractionType == "Like" {
		// 🔥 Enes'in getirdiği 401 Giriş Yap zırhı eklendi kanka
		if logReq.UserID == nil || *logReq.UserID == 0 {
			http.Error(w, "Beğenmek için giriş yapmalısınız.", 401)
			return
		}

		var exists int
		db.QueryRow("SELECT COUNT(*) FROM InteractionLogs WHERE UserID = @p1 AND TargetID = @p2 AND TargetType = @p3 AND InteractionType = 'Like'", logReq.UserID, logReq.TargetID, logReq.TargetType).Scan(&exists)
		if exists > 0 {
			db.Exec("DELETE FROM InteractionLogs WHERE UserID = @p1 AND TargetID = @p2 AND TargetType = @p3 AND InteractionType = 'Like'", logReq.UserID, logReq.TargetID, logReq.TargetType)
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(map[string]string{"message": "Beğeni kaldırıldı"})
			return
		}
	}

	db.Exec("INSERT INTO InteractionLogs (UserID, TargetID, TargetType, InteractionType) VALUES (@p1, @p2, @p3, @p4)", logReq.UserID, logReq.TargetID, logReq.TargetType, logReq.InteractionType)
	w.WriteHeader(200)
}

func getEntityStatsHandler(w http.ResponseWriter, r *http.Request) {
	targetID := r.URL.Query().Get("targetId")
	targetType := r.URL.Query().Get("targetType")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var viewCount, likeCount, commentCount int
	var avgRating float64
	db.QueryRow("SELECT COUNT(*) FROM InteractionLogs WHERE TargetID = @p1 AND TargetType = @p2 AND InteractionType = 'View'", targetID, targetType).Scan(&viewCount)
	db.QueryRow("SELECT COUNT(*) FROM InteractionLogs WHERE TargetID = @p1 AND TargetType = @p2 AND InteractionType = 'Like'", targetID, targetType).Scan(&likeCount)
	db.QueryRow("SELECT COUNT(*) FROM Comments WHERE TargetID = @p1 AND TargetType = @p2", targetID, targetType).Scan(&commentCount)
	db.QueryRow("SELECT ISNULL(AVG(CAST(Rating AS FLOAT)), 0) FROM Comments WHERE TargetID = @p1 AND TargetType = @p2 AND Rating > 0", targetID, targetType).Scan(&avgRating)
	reservations := 0
	occupancy := 0.0
	if targetType == "Workshop" {
		db.QueryRow("SELECT ISNULL(SUM(ParticipantCount), 0) FROM WorkshopEnrollments WHERE WorkshopId = @p1", targetID).Scan(&reservations)
		var capacity int
		db.QueryRow("SELECT Capacity FROM Workshops WHERE Id = @p1", targetID).Scan(&capacity)
		if capacity > 0 {
			occupancy = (float64(reservations) / float64(capacity)) * 100
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"views": viewCount, "likes": likeCount, "comments": commentCount, "avgRating": avgRating, "reservations": reservations, "occupancy": occupancy})
}

func adminDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var totalArtworks, totalWorkshops, totalUsers, totalTickets int
	db.QueryRow("SELECT COUNT(*) FROM Artworks").Scan(&totalArtworks)
	db.QueryRow("SELECT COUNT(*) FROM Workshops").Scan(&totalWorkshops)
	db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&totalUsers)
	db.QueryRow("SELECT COUNT(*) FROM SupportTickets WHERE Status != 'Çözüldü'").Scan(&totalTickets)
	json.NewEncoder(w).Encode(map[string]interface{}{"totalArtworks": totalArtworks, "totalWorkshops": totalWorkshops, "totalUsers": totalUsers, "activeTickets": totalTickets})
}

func getPopularArtworksHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	query := `
		SELECT TOP 5 
			a.Id, a.Title, ar.ArtistName, a.Category, a.ImageUrl,
			(SELECT COUNT(*) FROM InteractionLogs i WHERE i.TargetID = a.Id AND i.TargetType = 'Artwork' AND i.InteractionType = 'View') as ViewCount
		FROM Artworks a
		JOIN Artists ar ON a.ArtistID = ar.ArtistID
		ORDER BY ViewCount DESC`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Veri çekilemedi", 500)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, viewCount int
		var title, artist, category, imageUrl string
		rows.Scan(&id, &title, &artist, &category, &imageUrl, &viewCount)
		results = append(results, map[string]interface{}{
			"id": id, "title": title, "artist": artist, "category": category, "imageUrl": imageUrl, "views": viewCount,
		})
	}
	json.NewEncoder(w).Encode(results)
}

func saveComparisonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok || userID <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req struct {
		Title      string `json:"title"`
		TargetType string `json:"targetType"`
		TargetIDs  string `json:"targetIds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 🚀 İŞTE ÇÖZÜM: CreatedAt sütununa GETDATE() fonksiyonunu gönderiyoruz
	query := `INSERT INTO Comparisons (UserID, Title, TargetType, TargetIDs, CreatedAt) 
              VALUES (@p1, @p2, @p3, @p4, GETDATE())`

	_, err = db.Exec(query,
		sql.Named("p1", userID),
		sql.Named("p2", req.Title),
		sql.Named("p3", req.TargetType),
		sql.Named("p4", req.TargetIDs),
	)

	if err != nil {
		// Eğer hala hata verirse, hatayı log'a yazıyoruz ki görebilelim
		log.Printf("SQL HATASI: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "SQL Hatası: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Analiz başarıyla kaydedildi!"})
}
func getComparisonsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Sütun sırasını veritabanındaki tablo ile birebir eşledik (ComparisonID, UserID, Title, TargetType, TargetIDs, CreatedAt)
	query := `SELECT ComparisonID, UserID, Title, TargetType, TargetIDs, CreatedAt 
              FROM Comparisons 
              WHERE UserID = @p1 
              ORDER BY CreatedAt DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("DB Sorgu Hatası: %v", err)
		http.Error(w, "Veri çekme hatası", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comparisons []Comparison
	for rows.Next() {
		var c Comparison
		// Scan sırası SQL sorgusundakiyle aynı olmalı!
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.TargetType, &c.TargetIDs, &c.CreatedAt); err != nil {
			log.Printf("Scan hatası: %v", err)
			continue
		}
		comparisons = append(comparisons, c)
	}

	// Eğer hiç kayıt yoksa boş array dön, hata dönme
	if comparisons == nil {
		comparisons = []Comparison{}
	}

	json.NewEncoder(w).Encode(comparisons)
}
func updateComparisonTitleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ComparisonID int    `json:"comparisonId"`
		Title        string `json:"title"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value(userIDKey).(int)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("UPDATE Comparisons SET Title = @p1 WHERE ComparisonID = @p2 AND UserID = @p3", req.Title, req.ComparisonID, userID)
	w.WriteHeader(200)
}

func deleteComparisonHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	userID := r.Context().Value(userIDKey).(int)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("DELETE FROM Comparisons WHERE ComparisonID = @p1 AND UserID = @p2", id, userID)
	w.WriteHeader(200)
}

func getArtistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

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

func updateBalanceHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email  string  `tax:"email"`
		Amount float64 `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		http.Error(w, "DB Hatası", 500)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Users SET Balance = ISNULL(Balance, 0) + @p1 WHERE Email = @p2", data.Amount, data.Email)
	if err != nil {
		http.Error(w, "SQL Hatası", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bakiye başarıyla güncellendi!"})
}

func checkCouponHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	db, _ := sql.Open("sqlserver", connString)
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
	_, err := db.Exec("UPDATE Artworks SET IsCampaign = 0, DiscountRate = 0")
	if err != nil {
		log.Println("Sıfırlama hatası:", err)
		return err
	}

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

	countToDiscount := int(math.Round(float64(total) * 0.4))
	if countToDiscount == 0 && total > 0 {
		countToDiscount = 1
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })

	for i := 0; i < countToDiscount; i++ {
		randomRate := r.Intn(41) + 10
		query := "UPDATE Artworks SET [DiscountRate] = @p1, [IsCampaign] = 1 WHERE CAST(Id AS INT) = @p2"
		_, err := db.Exec(query, sql.Named("p1", randomRate), sql.Named("p2", ids[i]))
		if err != nil {
			log.Printf("❌ ID %d güncellenemedi: %v", ids[i], err)
		}
	}
	return nil
}

func SeedCoupons(db *sql.DB) error {
	coupons := []struct {
		Code   string
		Amount float64
	}{
		{"SANAT100", 100.00},
		{"ILK500", 500.00},
	}

	for _, c := range coupons {
		var exists bool
		checkQuery := "SELECT CASE WHEN EXISTS (SELECT 1 FROM Coupons WHERE Code = @p1) THEN 1 ELSE 0 END"
		err := db.QueryRow(checkQuery, c.Code).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			insertQuery := "INSERT INTO Coupons (Code, DiscountAmount, IsActive) VALUES (@p1, @p2, 1)"
			_, err = db.Exec(insertQuery, c.Code, c.Amount)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/artworks", getArtworksHandler)
	mux.HandleFunc("/workshops", getWorkshopsHandler)
	mux.HandleFunc("/artist", getArtistHandler)
	mux.HandleFunc("/artists", getArtistsHandler)

	mux.HandleFunc("/profile", getUserProfileHandler)
	mux.HandleFunc("/profile/update", updateProfileHandler)
	mux.HandleFunc("/profile/change-password", changePasswordHandler)
	mux.HandleFunc("/favorites/add", addFavoriteHandler)
	mux.HandleFunc("/favorites/remove", deleteFavoriteHandler)
	mux.HandleFunc("/favorites/check", checkFavoriteHandler)
	mux.HandleFunc("/favorites/list", getUserFavoritesHandler)
	mux.HandleFunc("/workshops/enroll", enrollWorkshopHandler)
	mux.HandleFunc("/user-enrollments", getUserEnrollmentsHandler)
	mux.HandleFunc("/update-enrollment", updateEnrollmentHandler)
	mux.HandleFunc("/delete-enrollment", deleteEnrollmentHandler)

	// Hem senin slug yapını hem Enes'in düz rotasını tek adreste erittik kanka!
	mux.HandleFunc("/buy-artwork/", buyArtworkHandler)
	mux.HandleFunc("/artworks/buy", buyArtworkHandler)

	mux.HandleFunc("/user-purchases", getUserPurchasesHandler)
	mux.HandleFunc("/add-artwork", addArtworkHandler)
	mux.HandleFunc("/delete-artwork", deleteArtworkHandler)
	mux.HandleFunc("/add-workshop", addWorkshopHandler)
	mux.HandleFunc("/delete-workshop", deleteWorkshopHandler)
	mux.HandleFunc("/my-workshops", getMyWorkshopsHandler)

	mux.HandleFunc("/comments", getCommentsHandler)
	mux.HandleFunc("/log-interaction", logInteractionHandler)
	mux.HandleFunc("/entity-stats", getEntityStatsHandler)

	mux.HandleFunc("/tickets", isAuth(getUserTicketsHandler))
	mux.HandleFunc("/tickets/create", isAuth(createTicketHandler))
	mux.HandleFunc("/tickets/messages", isAuth(getTicketMessagesHandler))
	mux.HandleFunc("/tickets/messages/send", isAuth(sendTicketMessageHandler))
	mux.HandleFunc("/comments/add", isAuth(addCommentHandler))
	mux.HandleFunc("/comments/vote", isAuth(voteCommentHandler))
	mux.HandleFunc("/comments/reply", isAuth(addCommentReplyHandler))

	mux.HandleFunc("/comparisons", isAuth(getComparisonsHandler))
	mux.HandleFunc("/comparisons/save", isAuth(saveComparisonHandler))
	mux.HandleFunc("/comparisons/update", isAuth(updateComparisonTitleHandler))
	mux.HandleFunc("/comparisons/delete", isAuth(deleteComparisonHandler))

	mux.HandleFunc("/admin/tickets", isAdmin(getAllTicketsHandler))
	mux.HandleFunc("/admin/tickets/update", isAdmin(updateTicketStatusHandler))
	mux.HandleFunc("/admin/dashboard-stats", isAdmin(adminDashboardStatsHandler))
	mux.HandleFunc("/admin/popular-artworks", isAdmin(getPopularArtworksHandler))

	mux.HandleFunc("/profile/update-balance", updateBalanceHandler)
	mux.HandleFunc("/check-coupon", checkCouponHandler)
	mux.HandleFunc("/confirm-sale", confirmArtworksSaleHandler)
	mux.HandleFunc("/seller-orders", getSellerOrdersHandler)
	mux.HandleFunc("/seller-workshops-orders", getSellerWorkshopsOrdersHandler)
	mux.HandleFunc("/confirm-workshop-enrollment", confirmWorkshopEnrollmentHandler)
	mux.HandleFunc("/enroll-workshop-payment", enrollWorkshopWithPaymentHandler)

	fmt.Println("Günlük kampanyalar veritabanına işleniyor...")
	ApplyRandomCampaigns(db)

	fmt.Println("Sistem kuponları kontrol ediliyor...")
	SeedCoupons(db)

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

	fmt.Println("Server 8080 portunda çalışıyor...")
	err = http.ListenAndServe(":8080", finalHandler)
	if err != nil {
		log.Fatal(err)
	}
}
