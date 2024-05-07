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

	// 서버 시작
	if err := r.Run(":9999"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
