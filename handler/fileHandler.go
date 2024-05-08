package handler

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/TeamWAF/woorizip-attachment/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 액세스 키 ID와 시크릿 액세스 키 설정
const (
	accessKeyID     = "AKIATCKASJX3T3QQLSBO"
	secretAccessKey = "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"
	region          = "ap-northeast-2"
	bucketName      = "woorizip-attachment"
)

// AWS S3 클라이언트 생성
var svc *s3.S3

func init() {
	// 세션 생성
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// S3 클라이언트 생성
	svc = s3.New(sess)

	// debug
	log.Println("S3 client created")
	log.Println(svc.Endpoint)
}

// GenerateRandomString 주어진 길이의 랜덤한 문자열을 생성합니다.
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		log.Fatalf("Failed to generate random string: %v", err)
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

// UploadFileToS3 파일을 S3에 업로드하고 업로드한 파일의 키(경로)를 반환합니다.
func UploadFileToS3(file *multipart.FileHeader) (string, error) {
	// 파일 오픈
	fileContent, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer fileContent.Close()

	// 파일 이름 생성
	fileExtension := filepath.Ext(file.Filename)
	fileName := GenerateRandomString(30) + fileExtension

	// S3에 파일 업로드
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("attachments/" + fileName),
		Body:   fileContent,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return fileName, nil
}

// UploadFile 함수는 파일을 업로드합니다.
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName, err := UploadFileToS3(file)
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// DB에 파일 정보 저장
	dbInterface, exists := c.Get("db")
	if !exists {
		log.Println("Database connection not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	db, ok := dbInterface.(*gorm.DB)
	if !ok {
		log.Println("Database connection found in context is not of type *gorm.DB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	attachment := models.Attachment{
		ID:        uuid.New().String(),
		OwnerID:   uuid.New().String(),
		Path:      fileName, // S3에 업로드된 파일의 키
		Extension: filepath.Ext(fileName),
	}

	if err := db.Create(&attachment).Error; err != nil {
		log.Printf("Error saving attachment info to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save attachment info. Server error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "url": "https://test.teamwaf.app/attachment/" + attachment.ID})
}

// DownloadFileFromS3 S3로부터 파일을 다운로드합니다.
func DownloadFileFromS3(c *gin.Context, id string) error {
	dbInterface, exists := c.Get("db")
	if !exists {
		return fmt.Errorf("database connection not found in context")
	}

	db, ok := dbInterface.(*gorm.DB)
	if !ok {
		return fmt.Errorf("database connection found in context is not of type *gorm.DB")
	}

	var attachment models.Attachment
	if err := db.Where("id = ?", id).First(&attachment).Error; err != nil {
		return fmt.Errorf("error finding attachment in database: %v", err)
	}

	// S3에서 파일 다운로드
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("attachments/" + attachment.Path),
	})
	if err != nil {
		return fmt.Errorf("failed to download file from S3: %v", err)
	}
	defer resp.Body.Close()

	// 파일의 내용을 읽어 바이트 슬라이스로 변환
	fileBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	// 파일 확장자를 기반으로 MIME 타입 결정
	contentType := http.DetectContentType(fileBytes)

	// Content-Disposition 헤더를 설정하여 파일 이름 지정
	c.Header("Content-Disposition", "inline; filename="+filepath.Base(attachment.Path))
	// Content-Type 헤더를 설정하여 MIME 타입 지정
	c.Header("Content-Type", contentType)
	// HTTP 상태 코드를 설정하여 응답 시작
	c.Writer.WriteHeader(http.StatusOK)
	// 파일 내용을 클라이언트에 바이트 슬라이스로 보냄
	c.Writer.Write(fileBytes)

	return nil
}

// DownloadFile 함수는 파일을 다운로드합니다.
func DownloadFile(c *gin.Context) {
	id := c.Param("id")
	if err := DownloadFileFromS3(c, id); err != nil {
		log.Printf("Error downloading file: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
	}
}
