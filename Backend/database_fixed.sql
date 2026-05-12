USE SanatProjesi;
GO

-- Önce bağımlı tabloları sil (foreign key sırasına göre)
IF OBJECT_ID('CommentHelpfulVotes', 'U') IS NOT NULL DROP TABLE CommentHelpfulVotes;
IF OBJECT_ID('CommentReplies', 'U') IS NOT NULL DROP TABLE CommentReplies;
IF OBJECT_ID('Comments', 'U') IS NOT NULL DROP TABLE Comments;
IF OBJECT_ID('SupportMessages', 'U') IS NOT NULL DROP TABLE SupportMessages;
IF OBJECT_ID('SupportTickets', 'U') IS NOT NULL DROP TABLE SupportTickets;
IF OBJECT_ID('ArtworkPurchases', 'U') IS NOT NULL DROP TABLE ArtworkPurchases;
IF OBJECT_ID('WorkshopEnrollments', 'U') IS NOT NULL DROP TABLE WorkshopEnrollments;
IF OBJECT_ID('Favorites', 'U') IS NOT NULL DROP TABLE Favorites;
IF OBJECT_ID('Workshops', 'U') IS NOT NULL DROP TABLE Workshops;
IF OBJECT_ID('Artworks', 'U') IS NOT NULL DROP TABLE Artworks;
IF OBJECT_ID('Artists', 'U') IS NOT NULL DROP TABLE Artists;
IF OBJECT_ID('Users', 'U') IS NOT NULL DROP TABLE Users;

GO

-- KULLANICILAR (en temel tablo, diğerleri buna bağımlı)
CREATE TABLE Users (
    UserID INT PRIMARY KEY IDENTITY(1,1),
    FirstName NVARCHAR(100) NOT NULL,
    LastName NVARCHAR(100) NOT NULL,
    Email NVARCHAR(100) UNIQUE NOT NULL,
    Password NVARCHAR(255) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE()
);

CREATE TABLE Artists (
    ArtistID INT PRIMARY KEY IDENTITY(1,1),
    FullName NVARCHAR(200) NOT NULL,
    Biography NVARCHAR(MAX),
    Nationality NVARCHAR(100)
);

CREATE TABLE Artworks (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Title NVARCHAR(200) NOT NULL,
    ArtistID INT NOT NULL,
    Artist NVARCHAR(200),
    Price DECIMAL(18, 2),
    ImageUrl NVARCHAR(MAX),
    Description NVARCHAR(MAX),
    ViewsCount INT DEFAULT 0,
    CreatedAt DATETIME DEFAULT GETDATE(),
    FOREIGN KEY (ArtistID) REFERENCES Artists(ArtistID)
);

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
    ViewsCount INT DEFAULT 0,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Workshop_Instructor FOREIGN KEY (InstructorID) REFERENCES Users(UserID)
);

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

CREATE TABLE SupportTickets (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    Subject NVARCHAR(200) NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    Status NVARCHAR(50) DEFAULT 'Acik', 
    CreatedAt DATETIME DEFAULT GETDATE()
);

CREATE TABLE SupportMessages (
    Id INT PRIMARY KEY IDENTITY(1,1),
    TicketId INT NOT NULL,
    SenderEmail NVARCHAR(100) NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_SupportMessage_Ticket FOREIGN KEY (TicketId) REFERENCES SupportTickets(Id)
);

CREATE TABLE Comments (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    TargetType NVARCHAR(50) NOT NULL, 
    TargetId INT NOT NULL,
    Content NVARCHAR(MAX) NOT NULL,
    Rating INT, 
    CreatedAt DATETIME DEFAULT GETDATE()
);

CREATE TABLE CommentReplies (
    Id INT PRIMARY KEY IDENTITY(1,1),
    CommentId INT NOT NULL,
    UserEmail NVARCHAR(100) NOT NULL, 
    ReplyContent NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Reply_Comment FOREIGN KEY (CommentId) REFERENCES Comments(Id)
);

CREATE TABLE CommentHelpfulVotes (
    Id INT PRIMARY KEY IDENTITY(1,1),
    CommentId INT NOT NULL,
    UserEmail NVARCHAR(100) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT UC_UserCommentVote UNIQUE (UserEmail, CommentId),
    CONSTRAINT FK_Vote_Comment FOREIGN KEY (CommentId) REFERENCES Comments(Id)
);

GO

-- =====================
-- ÖRNEK VERİLER
-- =====================

-- Atölye eğitmeni olarak kullanılacak demo kullanıcı
INSERT INTO Users (FirstName, LastName, Email, Password)
VALUES ('Admin', 'Galeri', 'admin@sanatgalerisi.com', 'admin_placeholder');

INSERT INTO Artists (FullName, Biography, Nationality)
VALUES
('Vincent van Gogh', 'Vincent van Gogh (1853–1890), Batı sanat tarihinin en ünlü ve etkili isimlerinden biridir. Hollandalı ard izlenimci bir ressamdır. On yıldan biraz fazla bir süre içinde, yaklaşık 860 yağlı boya tablonun da dahil olduğu 2.100''den fazla sanat eseri üretti.', 'Hollandalı'),
('Salvador Dali', 'Salvador Dalí (1904–1989), İspanyol sürrealist ressamdır. Dalí, sadece resimleriyle değil, aynı zamanda heykeltıraşlık, fotoğrafçılık ve film yapımcılığı gibi pek çok alanda verdiği eserlerle de tanınır.', 'İspanyol');

INSERT INTO Workshops (Title, Description, InstructorID, AvailableDates, Location, Capacity, Price, ImageUrl)
VALUES
('Yağlı Boya Başlangıç Atölyesi', 'Temel teknikler eğitimi.', 1, '15 Haz 10:00, 22 Haz 14:00, 29 Haz 18:00', 'Atölye A', 15, 450.00, 'https://images.unsplash.com/photo-1513364776144-60967b0f800f'),
('Sürrealizm Paneli', 'Modern sanat söyleşisi.', 1, '20 Haz 11:30, 27 Haz 16:00', 'Konferans Salonu', 50, 0.00, 'https://images.unsplash.com/photo-1460661419201-fd4cecdf8a8b');

INSERT INTO Artworks (Title, ArtistID, Artist, Price, ImageUrl, Description)
VALUES
('Yıldızlı Gece', 1, 'Vincent van Gogh', 15000.00, 'https://upload.wikimedia.org/wikipedia/commons/thumb/e/ea/Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg/1280px-Van_Gogh_-_Starry_Night_-_Google_Art_Project.jpg', 'Van Gogh''un 1889 yılında akıl hastanesindeki odasının penceresinden görünen gece görünüşünü betimleyen tablodur.'),
('Arles''daki Yatak Odası', 1, 'Vincent van Gogh', 12000.00, 'https://upload.wikimedia.org/wikipedia/commons/thumb/7/76/Vincent_van_Gogh_-_De_slaapkamer_-_Google_Art_Project.jpg/1280px-Vincent_van_Gogh_-_De_slaapkamer_-_Google_Art_Project.jpg', 'Sanatçının Arles''daki "Sarı Ev" içindeki yatak odasını tasvir eder.'),
('Ayçiçekleri', 1, 'Vincent van Gogh', 25000.00, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRxh7asZrhZ_HIc8wGB7hVRtRQ6--A_JNm0qw&s', 'Van Gogh''un dostu Gauguin''i etkilemek için yaptığı en ünlü serilerinden biridir.'),
('Belleğin Azmi', 2, 'Salvador Dali', 22000.00, 'https://upload.wikimedia.org/wikipedia/en/d/dd/The_Persistence_of_Memory.jpg', 'Sürrealizmin en bilinen eserlerinden biri olan bu tablo, eriyen cep saatlerini konu alır.');

GO

-- Kontrol
SELECT 'Users' as Tablo, COUNT(*) as Kayit FROM Users
UNION ALL SELECT 'Artists', COUNT(*) FROM Artists
UNION ALL SELECT 'Artworks', COUNT(*) FROM Artworks
UNION ALL SELECT 'Workshops', COUNT(*) FROM Workshops;
