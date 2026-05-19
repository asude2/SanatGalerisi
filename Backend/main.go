package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

// --- MODELLER (Eski + Yeni) ---

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserRole  string `json:"role"` // 'User', 'Instructor', 'Admin'
	Biography string `json:"biography"`
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
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ArtistID    int     `json:"artistId"`
	Artist      string  `json:"artist"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"imageUrl"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Rating      float64 `json:"rating"`
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
	Rating         float64 `json:"rating"`
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
	ComparisonID int       `json:"comparisonId"`
	UserID       int       `json:"userId"`
	Title        string    `json:"title"`
	TargetType   string    `json:"targetType"`
	TargetIDs    string    `json:"targetIds"`
	CreatedAt    time.Time `json:"createdAt"`
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
		if authHeader == "" {
			http.Error(w, "Yetkilendirme başlığı eksik", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Geçersiz token", http.StatusUnauthorized)
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
		role := r.Context().Value(userRoleKey).(string)
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

// --- HANDLERS (Hepsi) ---

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
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	rows, _ := db.Query(`SELECT aw.Id, aw.Title, aw.ArtistID, aw.Price, aw.ImageUrl, ISNULL(aw.Description, ''), ISNULL(aw.Category, ''), ISNULL(ar.ArtistName, 'Bilinmeyen Sanatçı') FROM Artworks aw LEFT JOIN Artists ar ON aw.ArtistID = ar.ArtistID`)
	defer rows.Close()
	var artworks []Artwork
	for rows.Next() {
		var a Artwork
		rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.Price, &a.ImageUrl, &a.Description, &a.Category, &a.Artist)
		artworks = append(artworks, a)
	}
	json.NewEncoder(w).Encode(artworks)
}

func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var u User
	var biography sql.NullString
	query := `SELECT u.FirstName, u.LastName, u.Email, ISNULL(a.Biography, '') FROM Users u LEFT JOIN Artists a ON u.UserID = a.UserID WHERE u.Email = @p1`
	err := db.QueryRow(query, email).Scan(&u.FirstName, &u.LastName, &u.Email, &biography)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", 404)
		return
	}
	if biography.Valid {
		u.Biography = biography.String
	}
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
	email := r.URL.Query().Get("email")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	query := `SELECT aw.Id, aw.Title, aw.ArtistID, aw.Price, aw.ImageUrl, ISNULL(ar.ArtistName, 'Bilinmeyen Sanatçı') FROM Favorites f JOIN Artworks aw ON f.ArtworkId = aw.Id LEFT JOIN Artists ar ON aw.ArtistID = ar.ArtistID WHERE f.UserEmail = @p1`
	rows, _ := db.Query(query, email)
	defer rows.Close()
	var favs []Artwork
	for rows.Next() {
		var a Artwork
		rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.Price, &a.ImageUrl, &a.Artist)
		favs = append(favs, a)
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

func getUserEnrollmentsHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	query := `SELECT e.Id, w.Title, e.ParticipantCount, e.ReservedDate, e.CreatedAt, w.Location, w.AvailableDates FROM WorkshopEnrollments e JOIN Workshops w ON e.WorkshopId = w.Id WHERE e.UserEmail = @p1`
	rows, _ := db.Query(query, email)
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var id, pCount int
		var title, rDate, cAt, location, aDates string
		rows.Scan(&id, &title, &pCount, &rDate, &cAt, &location, &aDates)
		results = append(results, map[string]interface{}{
			"id": id, "workshopTitle": title, "participantCount": pCount, "reservedDate": rDate, "createdAt": cAt, "location": location, "availableDates": aDates,
		})
	}
	json.NewEncoder(w).Encode(results)
}

func updateEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID               int    `json:"id"`
		ParticipantCount int    `json:"participantCount"`
		ReservedDate     string `json:"reservedDate"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("UPDATE WorkshopEnrollments SET ParticipantCount = @p1, ReservedDate = @p2 WHERE Id = @p3", data.ParticipantCount, data.ReservedDate, data.ID)
	w.WriteHeader(200)
}

func deleteEnrollmentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("DELETE FROM WorkshopEnrollments WHERE Id = @p1", id)
	w.WriteHeader(200)
}

func buyArtworkHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string  `json:"email"`
		ArtworkId int     `json:"artworkId"`
		Price     float64 `json:"price"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var exists int
	db.QueryRow("SELECT COUNT(*) FROM ArtworkPurchases WHERE ArtworkId = @p1", req.ArtworkId).Scan(&exists)
	if exists > 0 {
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(map[string]string{"message": "Eser zaten satılmış!"})
		return
	}
	db.Exec("INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, PurchasePrice) VALUES (@p1, @p2, @p3)", req.Email, req.ArtworkId, req.Price)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Satın alma başarılı!"})
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
	db.Exec("DELETE FROM Favorites WHERE ArtworkId = @p1", req.ArtworkID)
	db.Exec("DELETE FROM ArtworkPurchases WHERE ArtworkId = @p1", req.ArtworkID)
	db.Exec("DELETE FROM Artworks WHERE Id = @p1", req.ArtworkID)
	w.WriteHeader(200)
}

func getArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("name")
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var artistInfo struct {
		Name           string `json:"name"`
		Biography      string `json:"biography"`
		ArtworksCount  int    `json:"artworksCount"`
		WorkshopsCount int    `json:"workshopsCount"`
	}
	query := `SELECT ISNULL(a.ArtistName, ''), ISNULL(a.Biography, ''), COUNT(DISTINCT aw.Id), COUNT(DISTINCT w.Id) FROM Artists a LEFT JOIN Artworks aw ON a.ArtistID = aw.ArtistID LEFT JOIN Users u ON a.UserID = u.UserID LEFT JOIN Workshops w ON u.UserID = w.InstructorID WHERE a.ArtistName = @p1 GROUP BY a.ArtistID, a.ArtistName, a.Biography`
	db.QueryRow(query, artistName).Scan(&artistInfo.Name, &artistInfo.Biography, &artistInfo.ArtworksCount, &artistInfo.WorkshopsCount)
	json.NewEncoder(w).Encode(artistInfo)
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

func getArtistsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	rows, _ := db.Query(`SELECT a.ArtistID, a.ArtistName, a.Biography, a.Nationality, u.Email FROM Artists a JOIN Users u ON a.UserID = u.UserID`)
	defer rows.Close()
	var artists []map[string]interface{}
	for rows.Next() {
		var id int
		var name, bio, nationality, email string
		rows.Scan(&id, &name, &bio, &nationality, &email)
		artists = append(artists, map[string]interface{}{"id": id, "name": name, "biography": bio, "nationality": nationality, "email": email})
	}
	json.NewEncoder(w).Encode(artists)
}

// --- YENİ ÖZELLİK HANDLERLARI ---

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
	if tickets == nil { tickets = []SupportTicket{} }
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
	if tickets == nil { tickets = []SupportTicket{} }
	json.NewEncoder(w).Encode(tickets)
}

func sendTicketMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg SupportMessage
	json.NewDecoder(r.Body).Decode(&msg)
	userID := r.Context().Value(userIDKey).(int)
	
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	// Mantıksal Hata 6 Fix: Bilet durumunu kontrol et
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
	db.Exec("UPDATE SupportTickets SET UpdatedAt = GETDATE(), Status = 'Açık' WHERE TicketID = @p1", msg.TicketID) // Kullanıcı yazınca tekrar açılabilir veya durum güncellenebilir
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
	if msgs == nil { msgs = []SupportMessage{} }
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
	
	if c.TargetType == "Artwork" {
		isVerified = checkPurchase(userID, c.TargetID)
	} else if c.TargetType == "Workshop" {
		isVerified = checkEnrollment(userID, c.TargetID)
		// Mantıksal Hata 2 Fix: Atölye yorumu için katılım şartı
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
	if sortBy == "highest" {
		orderClause = "ORDER BY c.Rating DESC, c.CreatedAt DESC"
	} else if sortBy == "most_helpful" {
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
		VoteType  string `json:"voteType"` // 'Up' veya 'Down'
	}
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value(userIDKey).(int)

	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	// Mevcut oyu kontrol et
	var existingVote string
	err := db.QueryRow("SELECT VoteType FROM CommentVotes WHERE UserID = @p1 AND CommentID = @p2", userID, req.CommentID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// Yeni oy
		db.Exec("INSERT INTO CommentVotes (UserID, CommentID, VoteType) VALUES (@p1, @p2, @p3)", userID, req.CommentID, req.VoteType)
		if req.VoteType == "Up" {
			db.Exec("UPDATE Comments SET Upvotes = Upvotes + 1 WHERE CommentID = @p1", req.CommentID)
		} else {
			db.Exec("UPDATE Comments SET Downvotes = Downvotes + 1 WHERE CommentID = @p1", req.CommentID)
		}
	} else if existingVote == req.VoteType {
		// Oyu geri çek (Toggle off)
		db.Exec("DELETE FROM CommentVotes WHERE UserID = @p1 AND CommentID = @p2", userID, req.CommentID)
		if req.VoteType == "Up" {
			db.Exec("UPDATE Comments SET Upvotes = Upvotes - 1 WHERE CommentID = @p1", req.CommentID)
		} else {
			db.Exec("UPDATE Comments SET Downvotes = Downvotes - 1 WHERE CommentID = @p1", req.CommentID)
		}
	} else {
		// Oyu değiştir
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
	var req struct { CommentID int `json:"commentId"`; ReplyText string `json:"replyText"` }
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value(userIDKey).(int)
	role := r.Context().Value(userRoleKey).(string)

	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()

	// Mantıksal Hata 1 Fix: Sadece Admin ve Instructor yanıt verebilir
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

	// Beğeni işlemi ise:
	if logReq.InteractionType == "Like" {
		if logReq.UserID == nil || *logReq.UserID == 0 {
			http.Error(w, "Beğenmek için giriş yapmalısınız.", 401)
			return
		}

		// Zaten beğenmiş mi kontrol et
		var exists int
		db.QueryRow("SELECT COUNT(*) FROM InteractionLogs WHERE UserID = @p1 AND TargetID = @p2 AND TargetType = @p3 AND InteractionType = 'Like'", logReq.UserID, logReq.TargetID, logReq.TargetType).Scan(&exists)

		if exists > 0 {
			// Varsa SİL (Toggle - Beğeniyi geri al)
			db.Exec("DELETE FROM InteractionLogs WHERE UserID = @p1 AND TargetID = @p2 AND TargetType = @p3 AND InteractionType = 'Like'", logReq.UserID, logReq.TargetID, logReq.TargetType)
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(map[string]string{"message": "Beğeni kaldırıldı"})
			return
		}
		// Yoksa aşağıda INSERT edilecek
	}

	// View (Görüntülenme) her zaman kaydedilir, Like ise sadece yoksa buraya düşer
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
		if capacity > 0 { occupancy = (float64(reservations) / float64(capacity)) * 100 }
	}
	json.NewEncoder(w).Encode(map[string]interface{}{ "views": viewCount, "likes": likeCount, "comments": commentCount, "avgRating": avgRating, "reservations": reservations, "occupancy": occupancy })
}

func adminDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	var totalArtworks, totalWorkshops, totalUsers, totalTickets int
	db.QueryRow("SELECT COUNT(*) FROM Artworks").Scan(&totalArtworks)
	db.QueryRow("SELECT COUNT(*) FROM Workshops").Scan(&totalWorkshops)
	db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&totalUsers)
	db.QueryRow("SELECT COUNT(*) FROM SupportTickets WHERE Status != 'Çözüldü'").Scan(&totalTickets)
	json.NewEncoder(w).Encode(map[string]interface{}{ "totalArtworks": totalArtworks, "totalWorkshops": totalWorkshops, "totalUsers": totalUsers, "activeTickets": totalTickets })
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
	var c Comparison
	json.NewDecoder(r.Body).Decode(&c)
	userID := r.Context().Value(userIDKey).(int)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	db.Exec("INSERT INTO Comparisons (UserID, Title, TargetType, TargetIDs) VALUES (@p1, @p2, @p3, @p4)", userID, c.Title, c.TargetType, c.TargetIDs)
	w.WriteHeader(201)
}

func getComparisonsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	db, _ := sql.Open("sqlserver", connString)
	defer db.Close()
	rows, _ := db.Query("SELECT ComparisonID, Title, TargetType, TargetIDs, CreatedAt FROM Comparisons WHERE UserID = @p1 ORDER BY CreatedAt DESC", userID)
	defer rows.Close()
	var results []Comparison
	for rows.Next() {
		var c Comparison
		c.UserID = userID
		rows.Scan(&c.ComparisonID, &c.Title, &c.TargetType, &c.TargetIDs, &c.CreatedAt)
		results = append(results, c)
	}
	if results == nil { results = []Comparison{} }
	json.NewEncoder(w).Encode(results)
}

func updateComparisonTitleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct { ComparisonID int `json:"comparisonId"`; Title string `json:"title"` }
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

func main() {
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
	http.ListenAndServe(":8080", finalHandler)
}
