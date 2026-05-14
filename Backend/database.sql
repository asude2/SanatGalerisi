USE SanatProjesi;
GO

-- 1. ADIM: Eski tabloları bağımlılık sırasına göre sil (Varsa)
IF OBJECT_ID('CommentHelpfulVotes', 'U') IS NOT NULL DROP TABLE CommentHelpfulVotes;
IF OBJECT_ID('CommentReplies', 'U') IS NOT NULL DROP TABLE CommentReplies;
IF OBJECT_ID('Comments', 'U') IS NOT NULL DROP TABLE Comments;
IF OBJECT_ID('SupportMessages', 'U') IS NOT NULL DROP TABLE SupportMessages;
IF OBJECT_ID('SupportTickets', 'U') IS NOT NULL DROP TABLE SupportTickets;
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

-- 7. ADIM: Yorumlar ve Etkileşimler
CREATE TABLE Comments (
    Id INT PRIMARY KEY IDENTITY(1,1),
    TargetType NVARCHAR(50) NOT NULL, -- 'Artwork' veya 'Workshop'
    TargetId INT NOT NULL,
    UserEmail NVARCHAR(100) NOT NULL,
    UserName NVARCHAR(200),
    Content NVARCHAR(MAX) NOT NULL,
    Rating INT DEFAULT 0,
    HelpfulCount INT DEFAULT 0,
    CreatedAt DATETIME DEFAULT GETDATE()
);

CREATE TABLE CommentReplies (
    Id INT PRIMARY KEY IDENTITY(1,1),
    CommentId INT NOT NULL,
    UserEmail NVARCHAR(100) NOT NULL,
    UserName NVARCHAR(200),
    Content NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Reply_Comment FOREIGN KEY (CommentId) REFERENCES Comments(Id)
);

CREATE TABLE CommentHelpfulVotes (
    Id INT PRIMARY KEY IDENTITY(1,1),
    CommentId INT NOT NULL,
    UserEmail NVARCHAR(100) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Vote_Comment FOREIGN KEY (CommentId) REFERENCES Comments(Id),
    CONSTRAINT UC_UserVote UNIQUE (CommentId, UserEmail)
);

-- Destek Talepleri
CREATE TABLE SupportTickets (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    Subject NVARCHAR(200) NOT NULL,
    Status NVARCHAR(50) DEFAULT 'Açık',
    CreatedAt DATETIME DEFAULT GETDATE()
);

CREATE TABLE SupportMessages (
    Id INT PRIMARY KEY IDENTITY(1,1),
    TicketId INT NOT NULL,
    SenderEmail NVARCHAR(100) NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Message_Ticket FOREIGN KEY (TicketId) REFERENCES SupportTickets(Id)
);
GO




SELECT a.Title, r.ArtistName, a.Category, a.Price 
FROM Artworks a 
JOIN Artists r ON a.ArtistID = r.ArtistID;