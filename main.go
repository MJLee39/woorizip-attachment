package main

import (
	"fmt"
	"log"

	"github.com/TeamWAF/woorizip-attachment/models"
	"github.com/TeamWAF/woorizip-attachment/router"
	"github.com/TeamWAF/woorizip-attachment/utils"
)

func main() {

	db, err := utils.GetDatabase()
	if err != nil {
		fmt.Println(err)
	}

	// db 테이블 삭제후 재생성
	// db.Migrator().DropTable(&models.Attachment{})
	db.AutoMigrate(&models.Attachment{})

	// 라우터 초기화
	r := router.InitRouter()

	// template 파일 로드
	r.LoadHTMLGlob("tmpl/*.html")

	// 인증서 파일 경로 설정
	certFile := "certificate.pem"
	certChainFile := "certificate_chain.pem"

	// HTTPS 서버 설정
	err = r.RunTLS(":443", certFile, certChainFile)
	if err != nil {
		log.Fatal("Error starting HTTPS server:", err)
	}

}
