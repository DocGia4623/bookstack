package repository

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"fmt"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type BookRepository interface {
	CreateBook(request.BookCreateRequest) (models.Book, error)
}

type BookRepositoryImpl struct {
	DB *gorm.DB
}

func NewBookRepositoryImpl(Db *gorm.DB) BookRepository {
	return &BookRepositoryImpl{
		DB: Db,
	}
}

func (b *BookRepositoryImpl) CreateBook(request request.BookCreateRequest) (models.Book, error) {
	var result models.Book
	// Chuyển request thành entity Book
	err := copier.Copy(&result, request)
	if err != nil {
		return models.Book{}, fmt.Errorf("cant bind request" + err.Error())
	}

	// Thêm danh sách Tags
	for _, tag := range request.Tags {
		result.Tags = append(result.Tags, models.Tag{
			EntityID:   result.ShelveID,
			EntityType: "book",
			Name:       tag.Name,
			Value:      tag.Value,
		})
	}

	// Lưu vào database
	if err := b.DB.Create(&result).Error; err != nil {
		return models.Book{}, err
	}
	return result, nil
}
