package main

import (
	"fmt"

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
	r.Run(":9999")

}
