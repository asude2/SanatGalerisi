USE SanatProjesi;
GO

-- 1. Comments tablosuna Downvotes sütunu ekle (Eğer yoksa)
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('Comments') AND name = 'Downvotes')
BEGIN
    ALTER TABLE Comments ADD Downvotes INT DEFAULT 0;
END
GO

-- 2. Eski CommentUpvotes tablosunu daha genel bir CommentVotes tablosuna dönüştürmek için güncelle
-- Önce eski tabloyu silelim (test aşamasında olduğumuz için sorun olmaz)
IF OBJECT_ID('CommentUpvotes', 'U') IS NOT NULL DROP TABLE CommentUpvotes;
GO

IF OBJECT_ID('CommentVotes', 'U') IS NULL 
CREATE TABLE CommentVotes (
    UserID INT NOT NULL,
    CommentID INT NOT NULL,
    VoteType NVARCHAR(10) NOT NULL, -- 'Up' veya 'Down'
    CreatedAt DATETIME DEFAULT GETDATE(),
    PRIMARY KEY (UserID, CommentID),
    CONSTRAINT FK_Vote_User FOREIGN KEY (UserID) REFERENCES Users(UserID),
    CONSTRAINT FK_Vote_Comment FOREIGN KEY (CommentID) REFERENCES Comments(CommentID)
);
GO
