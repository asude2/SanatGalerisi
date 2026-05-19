USE SanatProjesi;
GO

--- =================================================================================
--- 1. TABLOLARI SILME AŞAMASI (Bağımlılık sırasına göre en çocuk tablodan başlar kanka)
--- =================================================================================
IF OBJECT_ID('CommentVotes', 'U') IS NOT NULL DROP TABLE CommentVotes;
IF OBJECT_ID('CommentReplies', 'U') IS NOT NULL DROP TABLE CommentReplies;
IF OBJECT_ID('Comments', 'U') IS NOT NULL DROP TABLE Comments;
IF OBJECT_ID('SupportMessages', 'U') IS NOT NULL DROP TABLE SupportMessages;
IF OBJECT_ID('SupportTickets', 'U') IS NOT NULL DROP TABLE SupportTickets;
IF OBJECT_ID('InteractionLogs', 'U') IS NOT NULL DROP TABLE InteractionLogs;

IF OBJECT_ID('WorkshopEnrollments', 'U') IS NOT NULL DROP TABLE WorkshopEnrollments;
IF OBJECT_ID('ArtworkPurchases', 'U') IS NOT NULL DROP TABLE ArtworkPurchases;
IF OBJECT_ID('Favorites', 'U') IS NOT NULL DROP TABLE Favorites;
IF OBJECT_ID('Workshops', 'U') IS NOT NULL DROP TABLE Workshops;
IF OBJECT_ID('Artworks', 'U') IS NOT NULL DROP TABLE Artworks;
IF OBJECT_ID('Artists', 'U') IS NOT NULL DROP TABLE Artists;
IF OBJECT_ID('Users', 'U') IS NOT NULL DROP TABLE Users;
IF OBJECT_ID('Coupons', 'U') IS NOT NULL DROP TABLE Coupons;
GO

--- =================================================================================
--- 2. ANA TABLOLARI OLUŞTURMA AŞAMASI
--- =================================================================================

-- USERS
CREATE TABLE Users (
    UserID INT PRIMARY KEY IDENTITY(1,1),
    FirstName NVARCHAR(100) NOT NULL,
    LastName NVARCHAR(100) NOT NULL,
    Email NVARCHAR(100) UNIQUE NOT NULL,
    Password NVARCHAR(255) NOT NULL,
    UserRole NVARCHAR(20) DEFAULT 'User',
    Balance DECIMAL(18, 2) DEFAULT 0.00,
    LastPurchasedCategory NVARCHAR(100) NULL,
    CreatedAt DATETIME DEFAULT GETDATE()
);

-- ARTISTS
CREATE TABLE Artists (
    ArtistID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Biography NVARCHAR(MAX),
    Nationality NVARCHAR(100),
    ArtistName NVARCHAR(200),
    CONSTRAINT FK_Artist_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

-- ARTWORKS
CREATE TABLE Artworks (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Title NVARCHAR(200) NOT NULL,
    ArtistID INT NOT NULL, 
    Price DECIMAL(18, 2),
    ImageUrl NVARCHAR(MAX),
    Description NVARCHAR(MAX),
    Category NVARCHAR(100),
    IsCampaign BIT DEFAULT 0,
    DiscountRate INT DEFAULT 0,
    IsSold BIT DEFAULT 0,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Artwork_Artist FOREIGN KEY (ArtistID) REFERENCES Artists(ArtistID)
);

-- WORKSHOPS
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

-- FAVORITES
CREATE TABLE Favorites (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL, 
    ArtworkId INT NOT NULL,                
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT UC_UserFavorite UNIQUE (UserEmail, ArtworkId)
);

-- WORKSHOP ENROLLMENTS
CREATE TABLE WorkshopEnrollments (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    WorkshopId INT NOT NULL,
    ParticipantCount INT DEFAULT 1, 
    ReservedDate NVARCHAR(100) NOT NULL,   
    CreatedAt DATETIME DEFAULT GETDATE(),
    Status NVARCHAR(50) DEFAULT 'Onay Bekliyor',
    CONSTRAINT FK_Enrollment_Workshop FOREIGN KEY (WorkshopId) REFERENCES Workshops(Id)
);

-- ARTWORK PURCHASES
CREATE TABLE ArtworkPurchases (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    ArtworkId INT NOT NULL,
    SellerID INT NULL,
    PurchasePrice DECIMAL(18, 2) NOT NULL, 
    Status NVARCHAR(50) DEFAULT 'Hazırlanıyor', 
    PaymentMethod NVARCHAR(50) DEFAULT 'Kredi Kartı',
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Purchase_Artwork FOREIGN KEY (ArtworkId) REFERENCES Artworks(Id),
    CONSTRAINT FK_Purchase_Seller FOREIGN KEY (SellerID) REFERENCES Users(UserID)
);

-- COUPONS
CREATE TABLE Coupons (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Code NVARCHAR(50) UNIQUE NOT NULL,
    DiscountAmount DECIMAL(18, 2) NOT NULL,
    IsActive BIT DEFAULT 1
);
GO

--- =================================================================================
--- 3. YENİ SİSTEMLER (MIGRATION DOSYALARININ ENTEGRASYONU)
--- =================================================================================

-- Destek Sistemi (Support System)
CREATE TABLE SupportTickets (
    TicketID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Subject NVARCHAR(200) NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    SupportType NVARCHAR(100),
    Status NVARCHAR(50) DEFAULT 'Açık', -- Açık, Beklemede, Çözüldü
    CreatedAt DATETIME DEFAULT GETDATE(),
    UpdatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Ticket_User FOREIGN KEY (UserID) REFERENCES dbo.Users(UserID)
);

CREATE TABLE SupportMessages (
    MessageID INT PRIMARY KEY IDENTITY(1,1),
    TicketID INT NOT NULL,
    SenderID INT NOT NULL, -- Kim gönderdi (Kullanıcı veya Admin)
    Message NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Message_Ticket FOREIGN KEY (TicketID) REFERENCES SupportTickets(TicketID),
    CONSTRAINT FK_Message_Sender FOREIGN KEY (SenderID) REFERENCES dbo.Users(UserID)
);

-- Yorum Sistemi (Comments) - Faz 3'teki Downvotes buraya entegre edildi kanka.
CREATE TABLE Comments (
    CommentID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    TargetID INT NOT NULL,
    TargetType NVARCHAR(50) NOT NULL, -- 'Artwork' veya 'Workshop'
    CommentText NVARCHAR(MAX) NOT NULL,
    Rating INT DEFAULT 0, -- 1-5 arası puanlama
    Upvotes INT DEFAULT 0, -- Faydalı buldum sayısı
    Downvotes INT DEFAULT 0, -- Faydalı bulmadım sayısı (Faz 3 Güncellemesi)
    IsVerified BIT DEFAULT 0, -- Doğrulanmış Alıcı / Katılımcı
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Comment_User FOREIGN KEY (UserID) REFERENCES dbo.Users(UserID)
);

-- Yorumlara Admin/Yönetici Yanıtları
CREATE TABLE CommentReplies (
    ReplyID INT PRIMARY KEY IDENTITY(1,1),
    CommentID INT NOT NULL,
    UserID INT NOT NULL, -- Yanıt veren Yönetici/Sorumlu
    ReplyText NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Reply_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID),
    CONSTRAINT FK_Reply_User FOREIGN KEY (UserID) REFERENCES dbo.Users(UserID)
);

-- Oylama Sistemi (Faz 3'teki gelişmiş CommentVotes yapısı doğrudan kuruldu)
CREATE TABLE CommentVotes (
    UserID INT NOT NULL,
    CommentID INT NOT NULL,
    VoteType NVARCHAR(10) NOT NULL, -- 'Up' veya 'Down'
    CreatedAt DATETIME DEFAULT GETDATE(),
    PRIMARY KEY (UserID, CommentID),
    CONSTRAINT FK_Vote_User FOREIGN KEY (UserID) REFERENCES Users(UserID),
    CONSTRAINT FK_Vote_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID)
);

-- İstatistikler için Etkileşim Logları (Görüntüleme vs.)
CREATE TABLE InteractionLogs (
    LogID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NULL, -- Giriş yapmamış kullanıcılar için NULL olabilir
    TargetID INT NOT NULL,
    TargetType NVARCHAR(50) NOT NULL, -- 'Artwork' veya 'Workshop'
    InteractionType NVARCHAR(50) NOT NULL, -- 'View', 'Like' vb.
    CreatedAt DATETIME DEFAULT GETDATE()
);
GO

--- =================================================================================
--- 4. TEST VE SIFIRLAMA SORGULARI (DATA MANIPULATION)
--- =================================================================================

-- Örnek Kupon Ekleme (Çakışma olmasın diye kontrol eklendi)
IF NOT EXISTS (SELECT 1 FROM Coupons WHERE Code = 'SANAT100')
BEGIN
    INSERT INTO Coupons (Code, DiscountAmount, IsActive) VALUES ('SANAT100', 100.00, 1);
END
GO

-- Satın Alınanları Sıfırlama Mantığı
UPDATE [dbo].[Artworks] SET IsSold = 0;
DELETE FROM [dbo].[ArtworkPurchases];
UPDATE [dbo].[Users] SET LastPurchasedCategory = NULL;
GO

-- Atölye Rezervasyonlarını Sıfırlama Mantığı
DELETE FROM [dbo].[WorkshopEnrollments];
DBCC CHECKIDENT ('[dbo].[WorkshopEnrollments]', RESEED, 0);
PRINT 'Atölye rezervasyonları, başvuru geçmişi ve yeni sistemler başarıyla hazırlandı! 🚀🎨';
GO