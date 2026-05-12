# Sanat Galerisi - Geliştirme İş Akışı (Workflow)

Bu iş akışı belgesi, projeye eklenecek yeni Müşteri Destek, Yorum, Değerlendirme ve İstatistik/Raporlama modüllerinin (10, 12, 13, 14, 15, 16. maddeler) baştan sona geliştirilme adımlarını tanımlar.

## Adım 1: Veritabanı ve Backend Model Güncellemeleri
- `SupportTickets`, `Comments`, `CommentRatings`, `EventAttendances`, `Purchases` tabloları (veya NoSQL koleksiyonları) oluşturulacak.
- Doğrulama kuralları (sadece katılanların yorum yapması, doğrulanmış satın alımlar vb.) için veri ilişkileri kurulacak.
- Eser ve Etkinlik modellerine istatistik alanları (görüntülenme, ortalama puan, rezervasyon sayısı vb.) eklenecek veya bu verilerin dinamik olarak hesaplanması için gerekli query'ler yazılacak.

## Adım 2: API (Backend) Endpoint'lerinin Geliştirilmesi
- **Destek Modülü (10):** Soru gönderme, canlı destek mesajları, destek talebi durumu sorgulama endpoint'leri.
- **Yorum Modülü (12, 14):** Yorum ekleme, yorumları listeleme ve yönetici/sorumlu yanıtı ekleme endpoint'leri.
- **Değerlendirme Modülü (13):** Yorumları puanlama, faydalı bulma ve filtreleme parametreli (en yeni, en yüksek puanlı vb.) listeleme endpoint'leri.
- **Raporlama Modülü (16):** Yönetici paneli için özet istatistik endpoint'leri.
- *Tüm endpoint'ler `15. Madde` kurallarına (Doğrulama ve Yetkilendirme) uygun olarak korunacaktır.*

## Adım 3: Frontend Geliştirmesi
- **Müşteri Destek (10):** İletişim formu bileşeni, canlı destek modülü/chat arayüzü ve destek talebi takip ekranının kodlanması.
- **Yorumlar ve Değerlendirmeler (12, 13, 14):** Eser ve Etkinlik detay sayfalarına eklenecek yorum sekmesi, yıldız/puan verme sistemi, "Faydalı buldum" butonu ve yönetici cevapları için arayüz bileşenleri oluşturulması.
- **İstatistik ve Raporlama (16):** 
  - Son kullanıcılar için eser/etkinlik detayında istatistiklerin (beğeni, görüntülenme, ortalama puan) gösterimi.
  - Yönetici (Admin) paneli için verilerin görselleştirildiği (grafikler vb.) rapor ekranı.

## Adım 4: Test ve Doğrulama
- "Yalnızca etkinliğe katılanlar etkinlik yorumu yapabilir" kuralının testi.
- "Yalnızca giriş yapanlar yorum yapabilir" kuralının testi.
- "Doğrulanmış değerlendirme" sisteminin satın alınmış eserler üzerinden testi.
- Yorum filtreleme, sıralama seçeneklerinin asenkron olarak doğru çalıştığının doğrulanması.
- Geliştirilen sayfaların responsivite ve UI/UX tasarımı açısından gözden geçirilmesi.

## Adım 5: Kod İnceleme (Code Review) ve Canlıya Alma
- Geliştirilen tüm modüllerin projenin genel mimarisine ve kurallarına uygunluğunun incelenmesi.
- Tespit edilen performans sorunları veya mantıksal hataların giderilmesi ve ardından kodun "production" ortamına alınması.
