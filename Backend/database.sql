USE SanatProjesi;
GO

-- 1. Kullanıcılar
IF OBJECT_ID('Users', 'U') IS NULL 
CREATE TABLE Users (
    UserID INT PRIMARY KEY IDENTITY(1,1),
    FirstName NVARCHAR(100) NOT NULL,
    LastName NVARCHAR(100) NOT NULL,
    Email NVARCHAR(100) UNIQUE NOT NULL,
    Password NVARCHAR(255) NOT NULL,
    UserRole NVARCHAR(20) DEFAULT 'User',
    CreatedAt DATETIME DEFAULT GETDATE()
);

-- 2. Sanatçılar
IF OBJECT_ID('Artists', 'U') IS NULL 
CREATE TABLE Artists (
    ArtistID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Biography NVARCHAR(MAX),
    Nationality NVARCHAR(100),
    ArtistName NVARCHAR(200),
    CONSTRAINT FK_Artist_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

-- 3. Sanat Eserleri
IF OBJECT_ID('Artworks', 'U') IS NULL 
CREATE TABLE Artworks (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Title NVARCHAR(200) NOT NULL,
    ArtistID INT NOT NULL, 
    Price DECIMAL(18, 2),
    ImageUrl NVARCHAR(MAX),
    Description NVARCHAR(MAX),
    Category NVARCHAR(100),
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Artwork_Artist FOREIGN KEY (ArtistID) REFERENCES Artists(ArtistID)
);

-- 4. Atölyeler
IF OBJECT_ID('Workshops', 'U') IS NULL 
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

-- 5. İşlem Tabloları (Favori, Kayıt, Satın Alma)
IF OBJECT_ID('Favorites', 'U') IS NULL 
CREATE TABLE Favorites (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL, 
    ArtworkId INT NOT NULL,                
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT UC_UserFavorite UNIQUE (UserEmail, ArtworkId)
);

IF OBJECT_ID('WorkshopEnrollments', 'U') IS NULL 
CREATE TABLE WorkshopEnrollments (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    WorkshopId INT NOT NULL,
    ParticipantCount INT DEFAULT 1, 
    ReservedDate NVARCHAR(100) NOT NULL,   
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Enrollment_Workshop FOREIGN KEY (WorkshopId) REFERENCES Workshops(Id)
);

IF OBJECT_ID('ArtworkPurchases', 'U') IS NULL 
CREATE TABLE ArtworkPurchases (
    Id INT PRIMARY KEY IDENTITY(1,1),
    UserEmail NVARCHAR(100) NOT NULL,
    ArtworkId INT NOT NULL,
    PurchasePrice DECIMAL(18, 2) NOT NULL, 
    Status NVARCHAR(50) DEFAULT 'Hazırlanıyor', 
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Purchase_Artwork FOREIGN KEY (ArtworkId) REFERENCES Artworks(Id)
);

-- 6. Destek ve Etkileşim (Yeni Fazlar)
IF OBJECT_ID('SupportTickets', 'U') IS NULL 
CREATE TABLE SupportTickets (
    TicketID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Subject NVARCHAR(200) NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    SupportType NVARCHAR(50),
    Status NVARCHAR(20) DEFAULT 'Open',
    CreatedAt DATETIME DEFAULT GETDATE(),
    UpdatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Ticket_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

IF OBJECT_ID('Comments', 'U') IS NULL 
CREATE TABLE Comments (
    CommentID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    TargetID INT NOT NULL,
    TargetType NVARCHAR(20) NOT NULL,
    CommentText NVARCHAR(MAX) NOT NULL,
    Rating INT CHECK (Rating >= 1 AND Rating <= 5),
    Upvotes INT DEFAULT 0,
    Downvotes INT DEFAULT 0,
    IsVerified BIT DEFAULT 0,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Comment_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

IF OBJECT_ID('CommentVotes', 'U') IS NULL 
CREATE TABLE CommentVotes (
    UserID INT NOT NULL,
    CommentID INT NOT NULL,
    VoteType NVARCHAR(10) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    PRIMARY KEY (UserID, CommentID),
    CONSTRAINT FK_Vote_User FOREIGN KEY (UserID) REFERENCES Users(UserID),
    CONSTRAINT FK_Vote_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID)
);

IF OBJECT_ID('InteractionLogs', 'U') IS NULL 
CREATE TABLE InteractionLogs (
    LogID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT,
    TargetID INT NOT NULL,
    TargetType NVARCHAR(20) NOT NULL,
    InteractionType NVARCHAR(20) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE()
);

-- Like'lar için benzersizlik kısıtlaması (Her kullanıcı bir şeyi 1 kere beğenebilir)
IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'UC_UserInteraction' AND object_id = OBJECT_ID('InteractionLogs'))
BEGIN
    CREATE UNIQUE INDEX UC_UserInteraction 
    ON InteractionLogs(UserID, TargetID, TargetType, InteractionType) 
    WHERE UserID IS NOT NULL AND InteractionType = 'Like';
END

IF OBJECT_ID('Comparisons', 'U') IS NULL 
CREATE TABLE Comparisons (
    ComparisonID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Title NVARCHAR(200),
    TargetType NVARCHAR(50) NOT NULL,
    TargetIDs NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Comparison_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);
GO