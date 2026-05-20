USE SanatProjesi;
GO

--- =================================================================================
--- 1. TABLOLARI SİLME AŞAMASI (Bağımlılık sırasına göre)
--- =================================================================================
IF OBJECT_ID('Comparisons', 'U') IS NOT NULL DROP TABLE Comparisons;
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

-- ARTWORKS (IsCampaign ve DiscountRate buraya dahil kanka!)
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
    Status NVARCHAR(50) DEFAULT N'Onay Bekliyor',
    CONSTRAINT FK_Enrollment_Workshop FOREIGN KEY (WorkshopId) REFERENCES Workshops(Id)
);

-- ARTWORK PURCHASES
CREATE TABLE ArtworkPurchases (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    ArtworkId INT NOT NULL,
    SellerID INT NULL,
    PurchasePrice DECIMAL(18, 2) NOT NULL, 
    Status NVARCHAR(50) DEFAULT N'Hazırlanıyor', 
    PaymentMethod NVARCHAR(50) DEFAULT N'Kredi Kartı',
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
--- 3. YENİ SİSTEMLER VE EKSİK TABLOLAR (Comparisons vs.)
--- =================================================================================

-- Comparisons (Artık en baştaki drop listesinde, sorunsuz kanka)
CREATE TABLE Comparisons (
    ComparisonID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Title NVARCHAR(255) NOT NULL,
    TargetType NVARCHAR(50) NOT NULL,
    TargetIDs NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Comparison_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

-- Destek Sistemi
CREATE TABLE SupportTickets (
    TicketID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Subject NVARCHAR(200) NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    SupportType NVARCHAR(100),
    Status NVARCHAR(50) DEFAULT N'Açık',
    CreatedAt DATETIME DEFAULT GETDATE(),
    UpdatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Ticket_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

CREATE TABLE SupportMessages (
    MessageID INT PRIMARY KEY IDENTITY(1,1),
    TicketID INT NOT NULL,
    SenderID INT NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Message_Ticket FOREIGN KEY (TicketID) REFERENCES SupportTickets(TicketID),
    CONSTRAINT FK_Message_Sender FOREIGN KEY (SenderID) REFERENCES Users(UserID)
);

-- Yorumlar, Oylamalar ve Etkileşimler
CREATE TABLE Comments (
    CommentID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    TargetID INT NOT NULL,
    TargetType NVARCHAR(50) NOT NULL,
    CommentText NVARCHAR(MAX) NOT NULL,
    Rating INT DEFAULT 0,
    Upvotes INT DEFAULT 0,
    Downvotes INT DEFAULT 0,
    IsVerified BIT DEFAULT 0,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Comment_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

CREATE TABLE CommentReplies (
    ReplyID INT PRIMARY KEY IDENTITY(1,1),
    CommentID INT NOT NULL,
    UserID INT NOT NULL,
    ReplyText NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Reply_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID),
    CONSTRAINT FK_Reply_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

CREATE TABLE CommentVotes (
    UserID INT NOT NULL,
    CommentID INT NOT NULL,
    VoteType NVARCHAR(10) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    PRIMARY KEY (UserID, CommentID),
    CONSTRAINT FK_Vote_User FOREIGN KEY (UserID) REFERENCES Users(UserID),
    CONSTRAINT FK_Vote_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID)
);

CREATE TABLE InteractionLogs (
    LogID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NULL,
    TargetID INT NOT NULL,
    TargetType NVARCHAR(50) NOT NULL,
    InteractionType NVARCHAR(50) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE()
);
GO

--- =================================================================================
--- 4. TEST VE SIFIRLAMA SORGULARI
--- =================================================================================

-- Kupon oluşturma
IF NOT EXISTS (SELECT 1 FROM Coupons WHERE Code = 'SANAT100')
BEGIN
    INSERT INTO Coupons (Code, DiscountAmount, IsActive) VALUES ('SANAT100', 100.00, 1);
END
GO

-- Satın alınanları ve rezervasyonları temizleme
UPDATE [dbo].[Artworks] SET IsSold = 0;
DELETE FROM [dbo].[ArtworkPurchases];
UPDATE [dbo].[Users] SET LastPurchasedCategory = NULL;
DELETE FROM [dbo].[WorkshopEnrollments];
DBCC CHECKIDENT ('[dbo].[WorkshopEnrollments]', RESEED, 0);

PRINT 'Tablolar güncellendi, her şey tıkırında kanka! 🚀🎨';
GO

-- Senin istediğin son sorgu:
SELECT Id, Title, IsCampaign, DiscountRate FROM Artworks;