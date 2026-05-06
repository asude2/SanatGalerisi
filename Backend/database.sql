USE SanatProjesi;
GO

IF OBJECT_ID('Artworks', 'U') IS NOT NULL DROP TABLE Artworks;
IF OBJECT_ID('Artists', 'U') IS NOT NULL DROP TABLE Artists;
IF OBJECT_ID('Workshops', 'U') IS NOT NULL DROP TABLE Workshops;
IF OBJECT_ID('Favorites', 'U') IS NOT NULL DROP TABLE Favorites;
GO

CREATE TABLE Artists (
    ArtistID INT PRIMARY KEY IDENTITY(1,1),
    FullName NVARCHAR(200) NOT NULL,
    Biography NVARCHAR(MAX),
    Nationality NVARCHAR(100)
);
GO

CREATE TABLE Artworks (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Title NVARCHAR(200) NOT NULL,
    ArtistID INT NOT NULL, 
    Artist NVARCHAR(200),  
    Price DECIMAL(18, 2),
    ImageUrl NVARCHAR(MAX),
    Description NVARCHAR(MAX),
    CreatedAt DATETIME DEFAULT GETDATE(),
    FOREIGN KEY (ArtistID) REFERENCES Artists(ArtistID)
);
GO

CREATE TABLE Workshops (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Title NVARCHAR(200) NOT NULL,
    Description NVARCHAR(MAX),
    InstructorID INT NOT NULL, 
    Date DATETIME NOT NULL,
    Location NVARCHAR(200),
    Capacity INT,
    Price DECIMAL(18, 2),
    ImageUrl NVARCHAR(MAX),
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Workshop_Instructor FOREIGN KEY (InstructorID) REFERENCES Users(UserID)
);
GO

CREATE TABLE Favorites (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL, 
    ArtworkId INT NOT NULL,        
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT UC_UserFavorite UNIQUE (UserEmail, ArtworkId) -- Bir kullanıcı aynı eseri iki kez favorileyemesin
);
GO

-- 3. Sanatçı Verilerini Ekle
INSERT INTO Artists (FullName, Biography, Nationality)
VALUES 
('Vincent van Gogh', 'Vincent van Gogh (1853–1890), Batı sanat tarihinin en ünlü ve etkili isimlerinden biridir. Hollandalı ard izlenimci bir ressamdır. On yıldan biraz fazla bir süre içinde, yaklaşık 860 yağlı boya tablonun da dahil olduğu 2.100''den fazla sanat eseri üretti.', 'Hollandalı'),
('Salvador Dali', 'Salvador Dalí (1904–1989), İspanyol sürrealist ressamdır. Dalí, sadece resimleriyle değil, aynı zamanda heykeltıraşlık, fotoğrafçılık ve film yapımcılığı gibi pek çok alanda verdiği eserlerle de tanınır.', 'İspanyol');
GO

-- 4. Eserleri Ekle (Vincent van Gogh için 1, Salvador Dali için 2 ID'sini kullanıyoruz)
INSERT INTO Artworks (Title, ArtistID, Artist, Price, ImageUrl, Description)
VALUES 
('Yıldızlı Gece', 1, 'Vincent van Gogh', 15000.00, 'https://upload.wikimedia.org/wikipedia/commons/thumb/e/ea/Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg/1280px-Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg', 'Van Gogh''un 1889 yılında akıl hastanesindeki odasının penceresinden görünen gece görünüşünü betimleyen tablodur.'),

('Arles''daki Yatak Odası', 1, 'Vincent van Gogh', 12000.00, 'https://upload.wikimedia.org/wikipedia/commons/thumb/7/76/Vincent_van_Gogh_-_De_slaapkamer_-_Google_Art_Project.jpg/1280px-Vincent_van_Gogh_-_De_slaapkamer_-_Google_Art_Project.jpg', 'Sanatçının Arles''daki "Sarı Ev" içindeki yatak odasını tasvir eder.'),

('Ayçiçekleri', 1, 'Vincent van Gogh', 25000.00, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRxh7asZrhZ_HIc8wGB7hVRtRQ6--A_JNm0qw&s', 'Van Gogh''un dostu Gauguin''i etkilemek için yaptığı en ünlü serilerinden biridir.'),

('Belleğin Azmi', 2, 'Salvador Dali', 22000.00, 'https://upload.wikimedia.org/wikipedia/en/d/dd/The_Persistence_of_Memory.jpg', 'Sürrealizmin en bilinen eserlerinden biri olan bu tablo, eriyen cep saatlerini konu alır.');
GO


INSERT INTO Workshops (Title, Description, InstructorID, Date, Location, Capacity, Price, ImageUrl)
VALUES 
('Yağlı Boya Başlangıç Atölyesi', 'Temel teknikleri öğreneceğiniz 3 saatlik eğitim.', 1, '2026-06-15 14:00:00', 'Atölye A', 15, 450.00, 'https://images.unsplash.com/photo-1513364776144-60967b0f800f'),
('Sürrealizm Paneli', 'Modern sanat söyleşisi.', 1, '2026-06-20 18:30:00', 'Konferans Salonu', 50, 0.00, 'https://images.unsplash.com/photo-1460661419201-fd4cecdf8a8b');
GO





