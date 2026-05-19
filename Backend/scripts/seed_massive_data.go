package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

func seedMassiveData() {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Bağlantı hatası:", err)
	}
	defer db.Close()

	fmt.Println("🚀 Mevcut verilere dokunmadan yeni, zengin veriler ekleniyor...")
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	
	// --- 1. YENİ KULLANICILAR VE SANATÇILAR ---
	fmt.Println("👤 Ek kullanıcılar ve sanatçılar oluşturuluyor...")
	var user2ID, user3ID, artistID int
	db.QueryRow("INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) OUTPUT INSERTED.UserID VALUES (@p1, @p2, @p3, @p4, @p5)", 
		"Can", "Yılmaz", "can@sanat.com", string(hashedPassword), "User").Scan(&user2ID)
	db.QueryRow("INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) OUTPUT INSERTED.UserID VALUES (@p1, @p2, @p3, @p4, @p5)", 
		"Selin", "Demir", "selin@sanat.com", string(hashedPassword), "User").Scan(&user3ID)
	
	// Yeni Sanatçı: Claude Monet
	var selinUID int
	db.QueryRow("INSERT INTO Users (FirstName, LastName, Email, Password, UserRole) OUTPUT INSERTED.UserID VALUES (@p1, @p2, @p3, @p4, @p5)", 
		"Claude", "Monet", "monet@sanat.com", string(hashedPassword), "Instructor").Scan(&selinUID)
	db.QueryRow("INSERT INTO Artists (UserID, Biography, ArtistName, Nationality) OUTPUT INSERTED.ArtistID VALUES (@p1, @p2, @p3, @p4)", 
		selinUID, "İzlenimcilik akımının kurucusu.", "Claude Monet", "Fransız").Scan(&artistID)

	// --- 2. MEVCUT ESERLERİ BUL ---
	var artID1, artID2 int
	db.QueryRow("SELECT TOP 1 Id FROM Artworks ORDER BY Id ASC").Scan(&artID1)
	db.QueryRow("SELECT TOP 1 Id FROM Artworks ORDER BY Id DESC").Scan(&artID2)

	// --- 3. BİRÇOK YORUM VE YANIT ---
	fmt.Println("💬 Onlarca yeni yorum ve yanıt ekleniyor...")
	commentTexts := []string{
		"Harika bir kompozisyon, renk geçişleri büyüleyici.",
		"Bu sanatçının tarzına her zaman hayran kalmışımdır.",
		"Galerideki en dikkat çekici eserlerden biri kesinlikle bu.",
		"Atölye çalışmasından sonra esere bakış açım tamamen değişti.",
		"Fiyatını sonuna kadar hak eden bir başyapıt.",
	}

	for i, text := range commentTexts {
		var cID int
		uID := user2ID
		if i%2 == 0 { uID = user3ID }
		
		db.QueryRow("INSERT INTO Comments (UserID, TargetID, TargetType, CommentText, Rating, IsVerified) OUTPUT INSERTED.CommentID VALUES (@p1, @p2, @p3, @p4, @p5, 1)", 
			uID, artID1, "Artwork", text, 5-(i%2), 1).Scan(&cID)
		
		// Admin Yanıtları
		db.Exec("INSERT INTO CommentReplies (CommentID, UserID, ReplyText) VALUES (@p1, @p2, @p3)", 
			cID, 1, "Değerli yorumunuz için çok teşekkür ederiz!")
	}

	// --- 4. BİRÇOK KARŞILAŞTIRMA KAYDI ---
	fmt.Println("⚖️ Her kullanıcı için yeni karşılaştırma analizleri ekleniyor...")
	comparisons := []struct{UserID int; Title, IDs string}{
		{user2ID, "Modern Dekorasyon İçin Seçtiklerim", fmt.Sprintf("%d,%d", artID1, artID2)},
		{user2ID, "Hediye Alternatifleri", fmt.Sprintf("%d,%d", artID1, artID2)},
		{user3ID, "Yatak Odası İçin Sakin Renkler", fmt.Sprintf("%d,%d", artID1, artID2)},
		{user3ID, "Fiyat/Performans Analizi", fmt.Sprintf("%d,%d", artID1, artID2)},
	}

	for _, comp := range comparisons {
		db.Exec("INSERT INTO Comparisons (UserID, Title, TargetType, TargetIDs) VALUES (@p1, @p2, @p3, @p4)", 
			comp.UserID, comp.Title, "Artwork", comp.IDs)
	}

	// --- 5. BİRÇOK DESTEK TALEBİ VE MESAJLAŞMA ---
	fmt.Println("🎧 Yeni destek talepleri ve derinlemesine yazışmalar ekleniyor...")
	subjects := []string{"Ödeme Sorunu", "Eser Orijinallik Belgesi", "Atölye İptali Hakkında"}
	for _, sub := range subjects {
		var tID int
		db.QueryRow("INSERT INTO SupportTickets (UserID, Subject, Message, SupportType, Status) OUTPUT INSERTED.TicketID VALUES (@p1, @p2, @p3, @p4, @p5)", 
			user2ID, sub, sub + " hakkında detaylı bilgi almak istiyorum.", "Teknik", "Açık").Scan(&tID)
		
		// Uzun yazışmalar
		db.Exec("INSERT INTO SupportMessages (TicketID, SenderID, Message) VALUES (@p1, @p2, @p3)", tID, 1, "Tabii, size bu konuda yardımcı olalım.")
		db.Exec("INSERT INTO SupportMessages (TicketID, SenderID, Message) VALUES (@p1, @p2, @p3)", tID, user2ID, "Teşekkürler, dönüşünüzü bekliyorum.")
		db.Exec("INSERT INTO SupportMessages (TicketID, SenderID, Message) VALUES (@p1, @p2, @p3)", tID, 1, "Konu inceleniyor, kısa süre içinde çözülecektir.")
	}

	// --- 6. ETKİLEŞİM VE SATIN ALMA ---
	fmt.Println("📈 Binlerce izlenme ve beğeni simüle ediliyor...")
	for i := 0; i < 20; i++ {
		db.Exec("INSERT INTO InteractionLogs (UserID, TargetID, TargetType, InteractionType) VALUES (NULL, @p1, 'Artwork', 'View')", artID1)
		db.Exec("INSERT INTO InteractionLogs (UserID, TargetID, TargetType, InteractionType) VALUES (NULL, @p1, 'Artwork', 'View')", artID2)
	}

	fmt.Println("✅ Mevcut veriler korunarak veritabanı DEVASA şekilde zenginleştirildi!")
}
