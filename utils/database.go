package utils

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gorm (mysql) 싱글톤 객체 생성

type Database struct {
	DB *gorm.DB // 데이터베이스 연결
}

var (
	dsn        string
	dbOnce     sync.Once
	dbInstance *Database
)

func init() {

	dsn = "teamwaf:dntmdrhkaudtjd12@tcp(testaws.teamwaf.app:3306)/teamwafdb?charset=utf8mb4&parseTime=True&loc=Local"
}

// InitDatabase Database 구조체를 초기화합니다.
func InitDatabase() *Database {
	dbOnce.Do(func() {
		dbInstance = &Database{}
	})
	return dbInstance
}

// GetDatabase 싱글톤 Database 객체를 반환합니다.
func GetDatabase() (*gorm.DB, error) {
	db := InitDatabase()

	// 데이터베이스 설정
	if db.DB == nil {
		var err error

		if db.DB, err = connectToDatabase(dsn); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %v", err)
		}

	}

	return db.DB, nil
}

// connectToDatabase 주어진 DSN에 대한 데이터베이스 연결을 생성합니다.
func connectToDatabase(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}
