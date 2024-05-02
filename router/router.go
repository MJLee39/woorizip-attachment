package router

import (
	"fmt"
	"net/http"

	"github.com/TeamWAF/woorizip-attachment/handler"
	"github.com/TeamWAF/woorizip-attachment/middleware"
	"github.com/TeamWAF/woorizip-attachment/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 라우터를 초기화합니다.
func InitRouter() *gin.Engine {

	db, err := utils.GetDatabase()
	if err != nil {
		fmt.Println(err)
	}

	// gin 엔진 생성 모든 로그 출력, 디버그 모드
	r := gin.Default()

	// CORS 미들웨어 설정
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.Use(middleware.DatabaseMiddleware(db))

	// 루트 경로에 대한 핸들러 설정
	r.GET("/attachment", func(c *gin.Context) {
		// 업로드 템플릿을 직접 반환
		c.HTML(http.StatusOK, "upload.tmpl", gin.H{})
	})

	r.POST("/attachment", handler.UploadFile)
	r.GET("/attachment/:id", handler.DownloadFile)
	// v1.DELETE("/attachment/:id", deleteFile)
	r.LoadHTMLGlob("tmpl/*.html")

	return r
}
