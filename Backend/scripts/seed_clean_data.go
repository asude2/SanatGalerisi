package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

var connString = "server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;TrustServerCertificate=true;"

func main() {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Bağlantı hatası:", err)
	}
	defer db.Close()

	fmt.Println("🚀 Veritabanı baştan aşağı yenileniyor...")
	tables := []string{"InteractionLogs", "CommentVotes", "CommentReplies", "Comments", "SupportMessages", "SupportTickets", "Comparisons", "WorkshopEnrollments", "ArtworkPurchases", "Favorites", "Artworks", "Workshops", "Artists", "Users"}
	for _, table := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		db.Exec(fmt.Sprintf("DBCC CHECKIDENT ('%s', RESEED, 0)", table))
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	
	// --- 1. KULLANICILAR ---
	fmt.Println("👤 Kullanıcılar oluşturuluyor...")
	var adminID, userID int
	db.QueryRow("INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) OUTPUT INSERTED.UserID VALUES (@p1, @p2, @p3, @p4, @p5)", 
		"Enes", "Admin", "admin@sanat.com", string(hashedPassword), "Admin").Scan(&adminID)
	db.QueryRow("INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) OUTPUT INSERTED.UserID VALUES (@p1, @p2, @p3, @p4, @p5)", 
		"Ahmet", "Sanatsever", "user@sanat.com", string(hashedPassword), "User").Scan(&userID)
	
	// Sanatçılar
	artists := []struct{First, Last, Email, Bio string}{
		{"Leonardo", "Da Vinci", "leo@sanat.com", "Rönesans dehası, çok yönlü sanatçı ve mucit."},
		{"Vincent", "Van Gogh", "vincent@sanat.com", "Duyguların ressamı, post-izlenimciliğin öncüsü."},
		{"Frida", "Kahlo", "frida@sanat.com", "Sürrealist Meksikalı sanatçı, otoportre kraliçesi."},
		{"Osman", "Hamdi Bey", "osman@sanat.com", "Türk ressam, arkeolog ve müzeci."},
	}
	artistIDs := make(map[string]int)
	instructorIDs := make(map[string]int)
	for _, a := range artists {
		var uID, artID int
		db.QueryRow("INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) OUTPUT INSERTED.UserID VALUES (@p1, @p2, @p3, @p4, @p5)", 
			a.First, a.Last, a.Email, string(hashedPassword), "Instructor").Scan(&uID)
		db.QueryRow("INSERT INTO Artists (UserID, Biography, ArtistName, Nationality) OUTPUT INSERTED.ArtistID VALUES (@p1, @p2, @p3, @p4)", 
			uID, a.Bio, a.First + " " + a.Last, "Global").Scan(&artID)
		artistIDs[a.First + " " + a.Last] = artID
		instructorIDs[a.First + " " + a.Last] = uID
	}

	// --- 2. ESERLER VE ATÖLYELER ---
	fmt.Println("🎨 Sanat Eserleri ve Atölyeler ekleniyor...")
	artworkList := []struct{Title, Artist, Cat, Img string; Price float64}{
		{"Mona Lisa", "Leonardo Da Vinci", "Rönesans", "https://images.unsplash.com/photo-1544642899-f0d6e5f6ed6f?w=800", 5000000},
		{"Yıldızlı Gece", "Vincent Van Gogh", "İzlenimci", "https://images.unsplash.com/photo-1579783902614-a3fb3927b6a5?w=800", 3500000},
		{"Belleğin Azmi", "Leonardo Da Vinci", "Sürrealizm", "https://images.unsplash.com/photo-1577083552431-6e5fd01aa342?w=800", 2800000},
		{"Kaplumbağa Terbiyecisi", "Osman Hamdi Bey", "Klasik", "https://images.unsplash.com/photo-1578301978693-85fa9c0320b9?w=800", 4200000},
		{"İnci Küpeli Kız", "Leonardo Da Vinci", "Barok", "https://images.unsplash.com/photo-1582555172866-f73bb12a2ab3?w=800", 1500000},
	}
	artIDs := []int{}
	for _, art := range artworkList {
		var aID int
		db.QueryRow("INSERT INTO Artworks (Title, ArtistID, Price, ImageUrl, Category, Description) OUTPUT INSERTED.Id VALUES (@p1, @p2, @p3, @p4, @p5, @p6)", 
			art.Title, artistIDs[art.Artist], art.Price, art.Img, art.Cat, art.Title + " eserinin profesyonel replikası.").Scan(&aID)
		artIDs = append(artIDs, aID)
	}

	// Atölyeler
	workshopList := []struct{Title, Desc, Instructor, Loc, Dates string; Cap int; Price float64}{
		{"Yağlı Boya Teknikleri", "Profesyonel yağlı boya eğitimi.", "Leonardo Da Vinci", "Salon A", "Pazartesi 18:00", 15, 1200},
		{"Modern Sanat Analizi", "Teorik inceleme ve tartışma.", "Vincent Van Gogh", "Konferans Salonu", "Cumartesi 14:00", 30, 750},
		{"Frida Stili Portre", "Otoportre ve sembolizm.", "Frida Kahlo", "Atölye Katı 2", "Salı 19:00", 12, 1500},
	}
	wIDs := []int{}
	for _, w := range workshopList {
		var wID int
		db.QueryRow("INSERT INTO Workshops (Title, Description, InstructorID, Location, Capacity, Price, ImageUrl, AvailableDates) OUTPUT INSERTED.Id VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)", 
			w.Title, w.Desc, instructorIDs[w.Instructor], w.Loc, w.Cap, w.Price, "https://images.unsplash.com/photo-1513364776144-60967b0f800f?w=800", w.Dates).Scan(&wID)
		wIDs = append(wIDs, wID)
	}

	// --- 3. İŞLEMLER (SATIN ALMA, FAVORİ) ---
	fmt.Println("🛒 İşlem verileri (Favori, Satın Alma, Kayıt) ekleniyor...")
	db.Exec("INSERT INTO Favorites (UserEmail, ArtworkId) VALUES (@p1, @p2)", "user@sanat.com", artIDs[0])
	db.Exec("INSERT INTO Favorites (UserEmail, ArtworkId) VALUES (@p1, @p2)", "user@sanat.com", artIDs[1])
	db.Exec("INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, PurchasePrice, Status) VALUES (@p1, @p2, @p3, @p4)", "user@sanat.com", artIDs[0], 5000000, "Tamamlandı")
	db.Exec("INSERT INTO ArtworkPurchases (UserEmail, ArtworkId, PurchasePrice, Status) VALUES (@p1, @p2, @p3, @p4)", "user@sanat.com", artIDs[1], 3500000, "Tamamlandı")
	db.Exec("INSERT INTO WorkshopEnrollments (UserEmail, WorkshopId, ParticipantCount, ReservedDate) VALUES (@p1, @p2, @p3, @p4)", "user@sanat.com", wIDs[0], 1, "2026-06-15")

	// --- 4. SOSYAL ETKİLEŞİM (YORUM, BEĞENİ) ---
	fmt.Println("💬 Yorumlar, Yanıtlar ve Beğeniler ekleniyor...")
	var c1, c2 int
	db.QueryRow("INSERT INTO Comments (UserID, TargetID, TargetType, CommentText, Rating, IsVerified) OUTPUT INSERTED.CommentID VALUES (@p1, @p2, @p3, @p4, @p5, 1)", userID, artIDs[0], "Artwork", "Büyüleyici bir eser, detaylar harika.", 5).Scan(&c1)
	db.Exec("INSERT INTO CommentReplies (CommentID, UserID, ReplyText) VALUES (@p1, @p2, @p3)", c1, adminID, "Yorumunuz için teşekkürler Ahmet Bey!")
	db.QueryRow("INSERT INTO Comments (UserID, TargetID, TargetType, CommentText, Rating, IsVerified) OUTPUT INSERTED.CommentID VALUES (@p1, @p2, @p3, @p4, @p5, 1)", userID, artIDs[1], "Artwork", "Fırça darbeleri çok etkileyici.", 4).Scan(&c2)
	db.Exec("INSERT INTO CommentReplies (CommentID, UserID, ReplyText) VALUES (@p1, @p2, @p3)", c2, instructorIDs["Vincent Van Gogh"], "Ruhumun bir parçası o tabloda, beğenmenize sevindim.")

	db.Exec("INSERT INTO InteractionLogs (UserID, TargetID, TargetType, InteractionType) VALUES (@p1, @p2, @p3, @p4)", userID, artIDs[0], "Artwork", "View")
	db.Exec("INSERT INTO InteractionLogs (UserID, TargetID, TargetType, InteractionType) VALUES (@p1, @p2, @p3, @p4)", userID, artIDs[0], "Artwork", "Like")
	db.Exec("INSERT INTO InteractionLogs (UserID, TargetID, TargetType, InteractionType) VALUES (@p1, @p2, @p3, @p4)", userID, artIDs[1], "Artwork", "View")

	// --- 5. ANALİZLER VE DESTEK ---
	fmt.Println("⚖️ Karşılaştırma Analizleri ve Destek Talepleri ekleniyor...")
	db.Exec("INSERT INTO Comparisons (UserID, Title, TargetType, TargetIDs) VALUES (@p1, @p2, @p3, @p4)", 
		userID, "Salon Duvarı İçin Seçenekler", "Artwork", fmt.Sprintf("%d,%d,%d", artIDs[0], artIDs[1], artIDs[4]))
	db.Exec("INSERT INTO Comparisons (UserID, Title, TargetType, TargetIDs) VALUES (@p1, @p2, @p3, @p4)", 
		userID, "Klasik vs Modern", "Artwork", fmt.Sprintf("%d,%d", artIDs[3], artIDs[2]))

	var ticketID int
	db.QueryRow("INSERT INTO SupportTickets (UserID, Subject, Message, SupportType, Status) OUTPUT INSERTED.TicketID VALUES (@p1, @p2, @p3, @p4, @p5)", 
		userID, "Eser Sigortası Hakkında", "Mona Lisa için sigorta belgesi alabilir miyim?", "Diğer", "Beklemede").Scan(&ticketID)
	db.Exec("INSERT INTO SupportMessages (TicketID, SenderID, Message) VALUES (@p1, @p2, @p3)", ticketID, adminID, "Tabii ki, faturanızla birlikte e-posta adresinize gönderilmiştir.")
	db.Exec("INSERT INTO SupportMessages (TicketID, SenderID, Message) VALUES (@p1, @p2, @p3)", ticketID, userID, "Aldım, teşekkür ederim.")

	fmt.Println("✅ TÜM VERİLER BAŞARIYLA YÜKLENDİ!")
	fmt.Println("🔑 GİRİŞ BİLGİLERİ:")
	fmt.Println("Standart Kullanıcı: user@sanat.com / 123456")
	fmt.Println("Yönetici (Admin): admin@sanat.com / 123456")
	fmt.Println("Sanatçı (Da Vinci): leo@sanat.com / 123456")
}
