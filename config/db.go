package config

import (
	"bookstack/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func ConnectDB(config *Config) *gorm.DB {
	// tạo chuỗi kết nối PostgreSQL
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.PostgresUser,
		config.PostgresPassword,
		config.DBHost,
		config.DBPort,
		config.PostgresDB,
	)
	// kết nối đến cơ sở dữ liệu
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connec to database: %v", err)
	}
	//Migrate
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Book{},
		&models.Chapter{},
		&models.Page{},
		&models.Shelve{},
		&models.Tag{},
		&models.RefreshToken{},
		&models.Permission{},
		&models.RolePermission{},
	}
	for _, model := range modelsToMigrate {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Fatalf("failed to migrate model %T: %v", model, err)
		}
	}
	DB = db
	log.Println("Connected to database successfully")
	return db
}

// ConnectRedis thiết lập kết nối Redis
func ConnectRedis(config *Config) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   config.RedisDB,
	})

	// Kiểm tra kết nối Redis
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	RedisClient = client
	return client
}
func Connect(config *Config) {
	ConnectDB(config)
	ConnectRedis(config)
}
