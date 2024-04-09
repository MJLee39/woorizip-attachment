package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DatabaseMiddleware 함수는 각 HTTP 요청의 컨텍스트에 데이터베이스 연결을 추가합니다.
func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 컨텍스트에 데이터베이스 연결 추가
		c.Set("db", db)
		c.Next()
	}
}
