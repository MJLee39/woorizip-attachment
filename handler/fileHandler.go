package handler

import (
	"crypto/rand"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/TeamWAF/woorizip-attachment/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		log.Fatal(err)
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

// UploadFile 파일을 업로드합니다.
func UploadFile(c *gin.Context) {

	dbInterface, exists := c.Get("db")
	// 데이터베이스 연결을 컨텍스트에서 가져오기
	if !exists {
		log.Println("Database connection not found in context")
		return
	}

	// 데이터베이스 연결 형변환 확인
	db, ok := dbInterface.(*gorm.DB)
	if !ok {
		log.Println("Database connection found in context is not of type *gorm.DB")
		return
	}

	// 파일 가져오기
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 파일 이름 생성
	fileExtension := filepath.Ext(file.Filename)
	fileName := GenerateRandomString(30) + fileExtension
	filePath := filepath.Join("attachments", fileName)

	// 파일 저장
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	attachment := models.Attchment{
		OwnerID:   uuid.New().String(),
		Path:      filePath,
		Extension: fileExtension,
	}

	if err := db.Create(&attachment).Error; err != nil {
		log.Printf("Error saving attachment info to database: %v", err)
		c.String(http.StatusInternalServerError, "Failed to save attachment info. Server error.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

// downloadFile 파일을 다운로드합니다.
func DownloadFile(c *gin.Context) {
	// id를 받아 db에서 파일 path를 찾아 해당 파일을 다운로드
	id := c.Param("id")

	dbInterface, exists := c.Get("db")
	// 데이터베이스 연결을 컨텍스트에서 가져오기
	if !exists {
		log.Println("Database connection not found in context")
		return
	}

	// 데이터베이스 연결 형변환 확인
	db, ok := dbInterface.(*gorm.DB)
	if !ok {
		log.Println("Database connection found in context is not of type *gorm.DB")
		return
	}

	var attachment models.Attchment
	if err := db.Where("id = ?", id).First(&attachment).Error; err != nil {
		log.Printf("Error finding attachment in database: %v", err)
		c.String(http.StatusNotFound, "Attachment not found")
		return
	}

	c.File(attachment.Path)

}
