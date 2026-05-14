# 🎨 Sanat Galerisi — Geliştirme İş Akışı

Bu belge, projenin geri kalan modüllerinin (10, 12, 13, 14, 15, 16) geliştirilme sürecini takip etmek için kullanılır.
Her maddeyi tamamlandığında `- [x]` olarak işaretle.

---

## ✅ Mevcut Durum (Tamamlananlar)

- [x] Kullanıcı kaydı ve girişi (JWT + bcrypt)
- [x] Kullanıcı rolleri: `User` / `Instructor`
- [x] Eser listeleme ve detay sayfası
- [x] Atölye listeleme ve detay sayfası
- [x] Sanatçı listeleme ve detay sayfası
- [x] Eser satın alma
- [x] Atölye kayıt (enrollment)
- [x] Favorilere ekleme / kaldırma
- [x] Instructor: Eser ekleme / silme
- [x] Instructor: Atölye ekleme / silme
- [x] Profil sayfası (satın almalar + atölyeler)

---

## 📋 Madde 10 — Müşteri Destek

### Backend
- [x] `SupportTickets` tablosu mevcut
- [x] `POST /support/create` — Destek talebi oluşturma
- [x] `GET /support/tickets?email=` — Kullanıcının taleplerini listeleme
- [ ] `SupportMessages` tablosu üzerinden mesaj zinciri (canlı destek)
- [ ] `POST /support/reply` — Yönetici yanıtı ekleme
- [ ] `GET /support/messages?ticketId=` — Mesajları getirme

### Frontend
- [x] `SupportView.vue` — İletişim formu ve talep oluşturma
- [x] Kullanıcının kendi taleplerini listeleme
- [ ] Destek talebi detay sayfası (mesaj zinciri)
- [ ] Talep durumu badge'i (Beklemede / Yanıtlandı / Kapatıldı)
- [ ] Yönetici yanıtını görüntüleme arayüzü

---

## 📋 Madde 12 — Yorum Ekleme

### Backend
- [x] `Comments` tablosu mevcut (`TargetType`, `TargetId`, `UserEmail`, `Content`)
- [x] `POST /comments/add` — Yorum ekleme
- [x] `GET /comments/list?targetType=&targetId=` — Yorum listeleme
- [ ] Etkinliğe katılım kontrolü (yalnızca kayıtlı kullanıcı yorum yapabilsin)
- [ ] Satın alınan eser için yorum izni kontrolü

### Frontend
- [x] `CommentSection.vue` — Yorum bileşeni (Artwork ve Workshop'ta kullanılıyor)
- [x] Yorum ekleme formu
- [x] Yorumları listeleme
- [ ] Eser yorumu için "Doğrulanmış Alıcı" rozeti
- [ ] Atölye yorumu için "Kayıtlı Katılımcı" rozeti

---

## 📋 Madde 13 — Yorumları Değerlendirme ve Filtreleme

### Backend
- [x] `CommentHelpfulVotes` tablosu mevcut
- [x] `POST /comments/helpful` — Faydalı oy toggle (ekle/kaldır)
- [x] `Rating INT` kolonu `Comments` tablosuna eklendi
- [x] `addCommentHandler` → Rating alanını kaydediyor
- [x] `GET /comments/list` endpoint'ine `sortBy` parametresi eklendi
  - [x] `sortBy=newest` — En yeni
  - [x] `sortBy=helpful` — En faydalı
  - [x] `sortBy=rating` — En yüksek puanlı
- [x] `GET /comments/average` — Ortalama puan endpoint'i

### Frontend
- [x] 👍 Faydalı bulma butonu (toggle, userVoted senkronu düzeltildi)
- [x] ⭐ Yıldız puanlama bileşeni (1-5 yıldız, hover efekti)
- [x] Sıralama butonları (En Yeni / En Faydalı / En Yüksek Puanlı)
- [x] Ortalama puan kartı (başlıkta, backend'den çekiliyor)
- [x] Her yorumda yıldız gösterimi
- [x] Renkli avatar + ratingLabels (Berbat → Mükemmel)

---

## 📋 Madde 14 — Yorumlara Yanıt Verme

### Backend
- [x] `CommentReplies` tablosu mevcut
- [ ] `POST /comments/reply` — Yoruma yanıt ekleme (Instructor/Admin)
- [ ] `GET /comments/replies?commentId=` — Yanıtları getirme

### Frontend
- [ ] Yorumun altında "Yanıt Ver" butonu (yalnızca Instructor/Admin görür)
- [ ] Yanıt ekleme formu (açılır/kapanır)
- [ ] Yanıtların yorumun altında girintili görüntülenmesi

---

## 📋 Madde 15 — Doğrulama ve Güvenilirlik

### Backend
- [x] Yorum eklemek için giriş zorunluluğu (token kontrolü)
- [ ] Atölye yorumu: `WorkshopEnrollments` tablosunda kayıt kontrolü
- [ ] Eser yorumu: `ArtworkPurchases` tablosunda satın alma kontrolü
- [ ] Yanıt ekleme: yalnızca `Instructor` veya Admin rolü

### Frontend
- [x] Router guard ile giriş yapmadan yorum yapılamaz
- [ ] "Doğrulanmış Alıcı" rozeti (satın alınmış eser yorumunda)
- [ ] "Kayıtlı Katılımcı" rozeti (katılınan atölye yorumunda)
- [ ] Yetkisiz kullanıcıya uyarı mesajı

---

## 📋 Madde 16 — İstatistik ve Raporlama

### Backend
- [ ] `Artworks` tablosuna `ViewsCount` sayacı güncelleme (`/artwork?id=` her açıldığında +1)
- [ ] `GET /artwork/stats?id=` — Eser istatistikleri (yorum sayısı, satın alma, ortalama puan)
- [ ] `GET /workshop/stats?id=` — Atölye istatistikleri (kayıt sayısı, doluluk oranı, ortalama puan)
- [ ] `GET /admin/report` — Yönetici özet raporu (tüm istatistikler)

### Frontend
- [ ] Eser detay sayfasında: 👁 Görüntülenme | 💬 Yorum | ⭐ Puan gösterimi
- [ ] Atölye detay sayfasında: 👥 Kayıt sayısı | 📊 Doluluk | ⭐ Puan gösterimi
- [ ] Admin rapor sayfası (`/admin/report`)
  - [ ] En çok görüntülenen eserler
  - [ ] En çok yorum alan eserler
  - [ ] En dolu atölyeler
  - [ ] Toplam satış istatistikleri

---

## 🧪 Madde 4 — Test ve Doğrulama

- [ ] "Yalnızca etkinliğe katılanlar etkinlik yorumu yapabilir" kuralı testi
- [ ] "Yalnızca giriş yapanlar yorum yapabilir" kuralı testi
- [ ] "Doğrulanmış değerlendirme" sistemi satın alınmış eserler üzerinden testi
- [ ] Yorum filtreleme ve sıralama asenkron çalışma testi
- [ ] Sayfaların mobil uyumluluğu (responsive) kontrolü
- [ ] Instructor ile User arasında rol bazlı arayüz farkı kontrolü

---

## 🚀 Madde 5 — Canlıya Alma

- [ ] Tüm `console.log` ve debug çıktılarını temizle
- [ ] Backend için production build (`go build -o sanatgalerisi.exe .`)
- [ ] Frontend için production build (`npm run build`)
- [ ] Genel kod incelemesi ve mimari uygunluk kontrolü

---

## 📊 Genel İlerleme

| Madde | Backend | Frontend | Durum |
|---|---|---|---|
| 10 - Müşteri Destek | %60 | %50 | 🟡 Devam ediyor |
| 12 - Yorum Ekleme | %80 | %80 | 🟡 Neredeyse tamam |
| 13 - Değerlendirme & Filtreleme | %30 | %20 | 🔴 Başlandı |
| 14 - Yorumlara Yanıt | %0 | %0 | 🔴 Başlanmadı |
| 15 - Doğrulama & Güvenilirlik | %40 | %30 | 🔴 Devam ediyor |
| 16 - İstatistik & Raporlama | %0 | %0 | 🔴 Başlanmadı |
