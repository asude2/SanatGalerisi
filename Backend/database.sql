USE SanatProjesi;
GO

IF OBJECT_ID('WorkshopEnrollments', 'U') IS NOT NULL DROP TABLE WorkshopEnrollments;
IF OBJECT_ID('ArtworkPurchases', 'U') IS NOT NULL DROP TABLE ArtworkPurchases;
IF OBJECT_ID('Favorites', 'U') IS NOT NULL DROP TABLE Favorites;
IF OBJECT_ID('Workshops', 'U') IS NOT NULL DROP TABLE Workshops;
IF OBJECT_ID('Artworks', 'U') IS NOT NULL DROP TABLE Artworks;
IF OBJECT_ID('Artists', 'U') IS NOT NULL DROP TABLE Artists;
IF OBJECT_ID('Users', 'U') IS NOT NULL DROP TABLE Users;
IF OBJECT_ID('Coupons', 'U') IS NOT NULL DROP TABLE Coupons;
GO

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

CREATE TABLE Artists (
    ArtistID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Biography NVARCHAR(MAX),
    Nationality NVARCHAR(100),
    ArtistName NVARCHAR(200),
    CONSTRAINT FK_Artist_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);

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
    SellerID INT NULL,
    PurchasePrice DECIMAL(18, 2) NOT NULL, 
    Status NVARCHAR(50) DEFAULT 'Hazırlanıyor', 
    PaymentMethod NVARCHAR(50) DEFAULT 'Kredi Kartı',
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Purchase_Artwork FOREIGN KEY (ArtworkId) REFERENCES Artworks(Id),
    CONSTRAINT FK_Purchase_Seller FOREIGN KEY (SellerID) REFERENCES Users(UserID)
);

CREATE TABLE Coupons (
    Id INT PRIMARY KEY IDENTITY(1,1),
    Code NVARCHAR(50) UNIQUE NOT NULL,
    DiscountAmount DECIMAL(18, 2) NOT NULL,
    IsActive BIT DEFAULT 1
);
GO





--satın alınanları sıfırlama
USE [SanatProjesi];
GO
-- 1. Tüm eserleri 'Satılmadı' (0) olarak işaretle
UPDATE [dbo].[Artworks] 
SET IsSold = 0;

DELETE FROM [dbo].[ArtworkPurchases];
-- 3. Kullanıcıların 'Son alınan kategori' bilgisini temizle (Öneri bandını test etmek için)
UPDATE [dbo].[Users] 
SET LastPurchasedCategory = NULL;
GO
--











