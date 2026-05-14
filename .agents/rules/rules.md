---
trigger: always_on
glob:
description: Sanat Galerisi projesi için geliştirme kuralları
---

# 🎨 Sanat Galerisi — Proje Geliştirme Kuralları

## 🏗️ Proje Mimarisi

Bu proje iki katmandan oluşur:
- **Backend:** Go (Golang) — `Backend/main.go`, port `8080`
- **Frontend:** Vue 3 (Composition API + `<script setup>`) — `src/`, port `5173` veya `5175`
- **Veritabanı:** Microsoft SQL Server (MSSQL) — `SanatProjesi` veritabanı, Windows Authentication

---

## 🗄️ Veritabanı Kuralları

### Tablo ve Kolon İsimleri
- Artists tablosunda sanatçı adı kolonu `FullName`'dir — `ArtistName` **kullanılmaz**.
- Artworks tablosunda `Category` kolonu **yoktur** — sorgularda kullanılmaz.
- Workshops tablosunda `InstructorID`, `Users.UserID`'ye foreign key ile bağlıdır.
- Artworks tablosunda hem `ArtistID` (FK) hem `Artist` (string) kolonu vardır.

### Bağlantı Dizisi (Connection String)
```
server=localhost;database=SanatProjesi;trusted_connection=yes;encrypt=disable;
```

### Silme Sırası (Foreign Key Bağımlılıkları)
Veri silinirken şu sıra takip edilmelidir:
1. `CommentHelpfulVotes` → `CommentReplies` → `Comments`
2. `SupportMessages` → `SupportTickets`
3. `Favorites` → `ArtworkPurchases` → `WorkshopEnrollments`
4. `Workshops` → `Artworks` → `Artists`
5. `Users`

---

## ⚙️ Backend Kuralları (Go)

### Genel Yapı
- Tüm handler'lar `main.go` içinde tanımlanır.
- Her handler başına CORS header'ları eklenir (`Access-Control-Allow-Origin: *`).
- OPTIONS istekleri her handler'da ayrıca karşılanır.
- Hata durumunda `http.Error(w, "...", statusCode)` kullanılır.

### Kimlik Doğrulama
- JWT token kullanılır. Secret key: `cok_gizli_anahtar_123`
- Token'da `email`, `firstName`, `role`, `userId` alanları bulunur.
- Token süresi: 24 saat.
- Şifreler `bcrypt` ile hash'lenir — düz metin şifre **saklanmaz**.

### Kullanıcı Rolleri
- `User` → Standart kullanıcı (eser/atölye görüntüleme, yorum, satın alma)
- `Instructor` → Eğitmen (AddArtwork, AddWorkshop, kendi atölyelerini yönetme)

### `update_db.go` Dosyası
- Bu dosya `//go:build ignore` tag'ı taşır, normal build'e dahil edilmez.
- Sadece DB migration için `go run update_db.go` ile çalıştırılır.

---

## 🖼️ Frontend Kuralları (Vue 3)

### Kimlik ve Oturum
- Login sonrası `localStorage`'a şunlar yazılır:
  - `userToken` → JWT token
  - `userEmail` → Kullanıcı e-posta
  - `userRole` → `User` veya `Instructor`
  - `userId` → Kullanıcı ID
- `router.beforeEach` ile token yoksa `/login`'e yönlendirilir.

### Bileşenler (Components)
- `CommentSection.vue` → `targetType` (Artwork/Workshop) ve `targetId` prop'u alır.
- `ArtistCard.vue` → Sanatçı listesi için kart bileşeni.
- `ArtworkCard.vue` → Eser listesi için kart bileşeni.

### Sayfa Rotaları
| Rota | Bileşen | Erişim |
|---|---|---|
| `/` | HomeView | Giriş yapan herkes |
| `/login` | LoginView | Herkese açık |
| `/register` | RegisterView | Herkese açık |
| `/profile` | ProfileView | Giriş yapan herkes |
| `/artwork/:id` | ArtworkDetail | Giriş yapan herkes |
| `/workshops` | WorkshopView | Giriş yapan herkes |
| `/workshops/:id` | WorkshopDetail | Giriş yapan herkes |
| `/artists` | ArtistsView | Giriş yapan herkes |
| `/artist/:name` | ArtistDetail | Giriş yapan herkes |
| `/add-artwork` | AddArtwork | Instructor |
| `/add-workshop` | AddWorkshop | Instructor |
| `/support` | SupportView | Giriş yapan herkes |

### API İstekleri
- Backend URL: `http://localhost:8080`
- İstekler `axios` ile yapılır.
- JSON body gönderilirken `Content-Type: application/json` header'ı eklenir.

---

## 🌿 Git & Branch Kuralları

- Ana geliştirme dalı: `enes`
- Takım arkadaşının dalı: `main`
- `main`'den yeni özellikler geldiğinde `git merge origin/main` yapılır.
- Merge sonrası mutlaka `go build` ile backend derleme kontrolü yapılır.
- Commit mesajları Türkçe veya İngilizce olabilir, kısa ve açıklayıcı olmalıdır.

---

## 🔌 API Endpoint Listesi

| Method | Endpoint | Açıklama |
|---|---|---|
| POST | `/register` | Kayıt |
| POST | `/login` | Giriş (JWT döner) |
| GET | `/profile` | Profil bilgisi (`?email=`) |
| PUT | `/profile/update` | Profil güncelleme |
| GET | `/artworks` | Tüm eserler |
| GET | `/artwork` | Tek eser (`?id=`) |
| POST | `/add-artwork` | Eser ekle (Instructor) |
| DELETE | `/delete-artwork` | Eser sil (`?id=`) |
| GET | `/workshops` | Tüm atölyeler |
| POST | `/enroll` | Atölye kaydı |
| GET | `/my-workshops` | Eğitmenin atölyeleri |
| POST | `/add-workshop` | Atölye ekle |
| DELETE | `/delete-workshop` | Atölye sil |
| GET | `/artists` | Tüm sanatçılar |
| GET | `/artist` | Tek sanatçı (`?name=`) |
| POST | `/artworks/buy` | Eser satın alma |
| GET | `/user-purchases` | Satın alınan eserler |
| GET | `/favorites` | Favoriler |
| POST | `/favorites/add` | Favoriye ekle |
| DELETE | `/favorites/remove` | Favoriden çıkar |
| POST | `/support/create` | Destek talebi oluştur |
| GET | `/support/tickets` | Destek talepleri (`?email=`) |
| POST | `/comments/add` | Yorum ekle |
| GET | `/comments/list` | Yorumları listele (`?targetType=&targetId=`) |
| POST | `/comments/helpful` | Faydalı oy (toggle) |
