USE SanatProjesi;
GO

-- 1. ADIM: Eski tabloları bağımlılık sırasına göre sil (Varsa)
IF OBJECT_ID('WorkshopEnrollments', 'U') IS NOT NULL DROP TABLE WorkshopEnrollments;
IF OBJECT_ID('ArtworkPurchases', 'U') IS NOT NULL DROP TABLE ArtworkPurchases;
IF OBJECT_ID('Favorites', 'U') IS NOT NULL DROP TABLE Favorites;
IF OBJECT_ID('Workshops', 'U') IS NOT NULL DROP TABLE Workshops;
IF OBJECT_ID('Artworks', 'U') IS NOT NULL DROP TABLE Artworks;
IF OBJECT_ID('Artists', 'U') IS NOT NULL DROP TABLE Artists;
IF OBJECT_ID('Users', 'U') IS NOT NULL DROP TABLE Users;
GO

-- 2. ADIM: Kullanıcılar Tablosu (Rol desteğiyle)
CREATE TABLE Users (
    UserID INT PRIMARY KEY IDENTITY(1,1),
    FirstName NVARCHAR(100) NOT NULL,
    LastName NVARCHAR(100) NOT NULL,
    Email NVARCHAR(100) UNIQUE NOT NULL,
    Password NVARCHAR(255) NOT NULL,
    UserRole NVARCHAR(20) DEFAULT 'User', -- 'User' veya 'Instructor'
    CreatedAt DATETIME DEFAULT GETDATE()
);

-- 3. ADIM: Sanatçılar Tablosu (Kullanıcılara bağlı)
CREATE TABLE Artists (
    ArtistID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Biography NVARCHAR(MAX),
    Nationality NVARCHAR(100),
    ArtistName NVARCHAR(200), -- Sanatçı adı (Kolaylık için)
    CONSTRAINT FK_Artist_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

-- 4. ADIM: Sanat Eserleri (Kategori desteği eklendi)
CREATE TABLE Artworks (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Title NVARCHAR(200) NOT NULL,
    ArtistID INT NOT NULL, 
    Price DECIMAL(18, 2),
    ImageUrl NVARCHAR(MAX),
    Description NVARCHAR(MAX),
    Category NVARCHAR(100), -- Frontend'den gelen kategori burada tutulacak
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Artwork_Artist FOREIGN KEY (ArtistID) REFERENCES Artists(ArtistID)
);
ALTER TABLE Artworks ADD IsSold BIT DEFAULT 0;

-- 5. ADIM: Atölyeler (Workshoplar)
CREATE TABLE Workshops (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Title NVARCHAR(200) NOT NULL,
    Description NVARCHAR(MAX),
    InstructorID INT NOT NULL, 
    Location NVARCHAR(200),
    Capacity INT,
    Price DECIMAL(18, 2),
    ImageUrl NVARCHAR(MAX),
    AvailableDates NVARCHAR(MAX),
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Workshop_Instructor FOREIGN KEY (InstructorID) REFERENCES Users(UserID)
);

-- 6. ADIM: Favoriler ve Diğer İşlem Tabloları
CREATE TABLE Favorites (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL, 
    ArtworkId INT NOT NULL,                
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT UC_UserFavorite UNIQUE (UserEmail, ArtworkId)
);

CREATE TABLE WorkshopEnrollments (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    WorkshopId INT NOT NULL,
    ParticipantCount INT DEFAULT 1, 
    ReservedDate NVARCHAR(100) NOT NULL,   
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Enrollment_Workshop FOREIGN KEY (WorkshopId) REFERENCES Workshops(Id)
);

CREATE TABLE ArtworkPurchases (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    ArtworkId INT NOT NULL,
    PurchasePrice DECIMAL(18, 2) NOT NULL, 
    Status NVARCHAR(50) DEFAULT 'Hazırlanıyor', 
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Purchase_Artwork FOREIGN KEY (ArtworkId) REFERENCES Artworks(Id)
);
GO
ALTER TABLE ArtworkPurchases ADD PaymentMethod NVARCHAR(50) DEFAULT 'Kredi Kartı';



-- ! 1. Kullanıcılar tablosuna bakiye sütunu ekle (Eğer yoksa)
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('Users') AND name = 'Balance')
BEGIN
    ALTER TABLE Users ADD Balance DECIMAL(18, 2) DEFAULT 0.00;
END

-- ! 2. Kuponlar tablosunu oluştur
IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID('Coupons') AND type = 'U')
BEGIN
    CREATE TABLE Coupons (
        Id INT PRIMARY KEY IDENTITY(1,1),
        Code NVARCHAR(50) UNIQUE NOT NULL,
        DiscountAmount DECIMAL(18, 2) NOT NULL,
        IsActive BIT DEFAULT 1
    );
END
-- Denemek için bir tane örnek kupon ekleyelim
INSERT INTO Coupons (Code, DiscountAmount, IsActive) VALUES ('SANAT100', 100.00, 1);



ALTER TABLE Artworks ADD IsCampaign BIT DEFAULT 0;
ALTER TABLE Artworks ADD DiscountRate INT DEFAULT 0;

-- SELECT a.Title, r.ArtistName, a.Category, a.Price 
-- FROM Artworks a 
-- JOIN Artists r ON a.ArtistID = r.ArtistID;


-- 1. Eserin kime ait olduğunu bilmemiz lazım (Satıcı)
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[Artworks]') AND name = 'OwnerID')
BEGIN
    ALTER TABLE Artworks ADD OwnerID INT FOREIGN KEY REFERENCES Users(Id);
END
GO

-- 2. Mevcut eserleri test için Asude'ye (diyelim ki ID'si 2) atayalım
UPDATE Artworks SET OwnerID = 2; 
GO



-- Hata veren kısmı düzeltiyoruz: Referans sütun adı 'UserID' olmalı
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[Artworks]') AND name = 'OwnerID')
BEGIN
    ALTER TABLE Artworks ADD OwnerID INT;
    -- Foreign Key'i UserID üzerinden bağlıyoruz
    ALTER TABLE Artworks ADD CONSTRAINT FK_Artworks_Owner FOREIGN KEY (OwnerID) REFERENCES Users(UserID);
END
GO


-- Test verisi (Asude'ye atama)
UPDATE Artworks SET OwnerID = 2; 
GO

-- ArtworkPurchases Tablosuna SellerID: Bir satış gerçekleştiğinde, o satışın hangi satıcıya ait olduğunu doğrudan bu tabloda tutmak sorguları çok hızlandırır.
ALTER TABLE ArtworkPurchases ADD SellerID INT;
ALTER TABLE ArtworkPurchases 
ADD CONSTRAINT FK_Purchase_Seller FOREIGN KEY (SellerID) REFERENCES Users(UserID);
GO



USE SanatProjesi;
GO

-- 1. Eserlerin bir sahibi olduğundan emin olalım (Test için hepsini Asude'ye -ID:2- atıyoruz)
-- Eğer senin Users tablanda Asude'nin ID'si farklıysa 2 yerine onu yaz kanka
UPDATE Artworks SET OwnerID = 2 WHERE OwnerID IS NULL OR OwnerID NOT IN (SELECT UserID FROM Users);
GO

-- 2. Eğer tabloda eski, hatalı kayıtlar kaldıysa onları temizleyelim ki çakışma yapmasın
DELETE FROM ArtworkPurchases WHERE SellerID NOT IN (SELECT UserID FROM Users);
GO

UPDATE Artworks SET OwnerID = 2 WHERE OwnerID IS NULL;

-- Kolon tiplerini görmek için:
EXEC sp_help 'Artworks';