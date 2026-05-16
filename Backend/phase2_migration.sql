USE SanatProjesi;
GO

-- Destek Sistemi (Support System)
IF OBJECT_ID('SupportTickets', 'U') IS NULL 
CREATE TABLE SupportTickets (
    TicketID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    Subject NVARCHAR(200) NOT NULL,
    Message NVARCHAR(MAX) NOT NULL,
    SupportType NVARCHAR(100),
    Status NVARCHAR(50) DEFAULT 'Açık', -- Açık, Beklemede, Çözüldü
    CreatedAt DATETIME DEFAULT GETDATE(),
    UpdatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Ticket_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);
GO

IF OBJECT_ID('SupportMessages', 'U') IS NULL 
CREATE TABLE SupportMessages (
    MessageID INT PRIMARY KEY IDENTITY(1,1),
    TicketID INT NOT NULL,
    SenderID INT NOT NULL, -- Kim gönderdi (Kullanıcı veya Admin)
    Message NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Message_Ticket FOREIGN KEY (TicketID) REFERENCES SupportTickets(TicketID),
    CONSTRAINT FK_Message_Sender FOREIGN KEY (SenderID) REFERENCES Users(UserID)
);
GO

-- Yorum Sistemi (Comments)
IF OBJECT_ID('Comments', 'U') IS NULL 
CREATE TABLE Comments (
    CommentID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NOT NULL,
    TargetID INT NOT NULL,
    TargetType NVARCHAR(50) NOT NULL, -- 'Artwork' veya 'Workshop'
    CommentText NVARCHAR(MAX) NOT NULL,
    Rating INT DEFAULT 0, -- 1-5 arası puanlama
    Upvotes INT DEFAULT 0, -- Faydalı buldum sayısı
    IsVerified BIT DEFAULT 0, -- Doğrulanmış Alıcı / Katılımcı
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Comment_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);
GO

-- Yorumlara Admin/Yönetici Yanıtları
IF OBJECT_ID('CommentReplies', 'U') IS NULL 
CREATE TABLE CommentReplies (
    ReplyID INT PRIMARY KEY IDENTITY(1,1),
    CommentID INT NOT NULL,
    UserID INT NOT NULL, -- Yanıt veren Yönetici/Sorumlu
    ReplyText NVARCHAR(MAX) NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_Reply_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID),
    CONSTRAINT FK_Reply_User FOREIGN KEY (UserID) REFERENCES Users(UserID)
);
GO

-- Yorum Beğeni/Faydalı Bulma (Upvotes)
IF OBJECT_ID('CommentUpvotes', 'U') IS NULL 
CREATE TABLE CommentUpvotes (
    UserID INT NOT NULL,
    CommentID INT NOT NULL,
    CreatedAt DATETIME DEFAULT GETDATE(),
    PRIMARY KEY (UserID, CommentID),
    CONSTRAINT FK_Upvote_User FOREIGN KEY (UserID) REFERENCES Users(UserID),
    CONSTRAINT FK_Upvote_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID)
);
GO

-- İstatistikler için Etkileşim Logları (Görüntüleme vs.)
IF OBJECT_ID('InteractionLogs', 'U') IS NULL 
CREATE TABLE InteractionLogs (
    LogID INT PRIMARY KEY IDENTITY(1,1),
    UserID INT NULL, -- Giriş yapmamış kullanıcılar için NULL olabilir
    TargetID INT NOT NULL,
    TargetType NVARCHAR(50) NOT NULL, -- 'Artwork' veya 'Workshop'
    InteractionType NVARCHAR(50) NOT NULL, -- 'View', 'Like' vb.
    CreatedAt DATETIME DEFAULT GETDATE()
);
GO
