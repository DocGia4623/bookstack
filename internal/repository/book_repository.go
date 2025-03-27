package repository

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"fmt"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type BookRepository interface {
	CreateBook(int, request.BookCreateRequest) (models.Book, error)
	GetAllBook() ([]models.Book, error)
	CreateShelve(int, request.ShelveCreateRequest) (models.Shelve, error)
	CreateChapter(request.BookChapterRequest) (models.Chapter, error)
	GetChaptersOfBook(int) ([]models.Chapter, error)
	AddPage(request.PageRequest) (models.Page, error)
	GetPageChapter(int) ([]models.Page, error)
}

type BookRepositoryImpl struct {
	DB *gorm.DB
}

func NewBookRepositoryImpl(Db *gorm.DB) BookRepository {
	return &BookRepositoryImpl{
		DB: Db,
	}
}

func (b *BookRepositoryImpl) GetPageChapter(chapterId int) ([]models.Page, error) {
	var pages []models.Page
	err := b.DB.Where("chapter_id = ?", chapterId).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (b *BookRepositoryImpl) AddPage(request request.PageRequest) (models.Page, error) {
	var page models.Page
	err := copier.Copy(&page, request)
	if err != nil {
		return models.Page{}, nil
	}
	err = b.DB.FirstOrCreate(&page, models.Page{
		Title:     request.Title,
		ChapterID: request.ChapterID,
	}).Error
	if err != nil {
		return models.Page{}, err
	}
	return page, err
}

func (b *BookRepositoryImpl) GetChaptersOfBook(bookID int) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := b.DB.Where("book_id = ?", bookID).Find(&chapters).Error
	if err != nil {
		return nil, err
	}
	return chapters, nil
}
func (b *BookRepositoryImpl) CreateChapter(request request.BookChapterRequest) (models.Chapter, error) {
	var chapter models.Chapter
	err := copier.Copy(&chapter, request)
	if err != nil {
		return models.Chapter{}, err
	}
	// Dùng FirstOrCreate để tìm hoặc tạo mới Chapter
	result := b.DB.FirstOrCreate(&chapter, models.Chapter{
		Title:  request.Title,
		BookID: request.BookID, // Điều kiện để xác định chương có tồn tại hay chưa
	})
	if result.Error != nil {
		return models.Chapter{}, result.Error
	}
	return chapter, nil

}
func (b *BookRepositoryImpl) GetAllBook() ([]models.Book, error) {
	var books []models.Book
	err := b.DB.Find(&books).Error
	if err != nil {
		return []models.Book{}, err
	}
	return books, nil
}

func (b *BookRepositoryImpl) FindTagByName(name string, entityType string) (*models.Tag, error) {
	var tag models.Tag
	result := b.DB.Where("name = ? AND entity_type = ?", name, entityType).Limit(1).Find(&tag)

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Nếu không tìm thấy tag, trả về nil
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &tag, nil
}

func (b *BookRepositoryImpl) CreateShelve(userId int, request request.ShelveCreateRequest) (models.Shelve, error) {
	var result models.Shelve
	err := copier.Copy(&result, request)
	if err != nil {
		return models.Shelve{}, fmt.Errorf("cant bind request: %w", err)
	}
	result.CreatedBy = uint(userId)

	var tags []models.Tag
	tagSet := make(map[string]struct{}) // Dùng map để tránh trùng lặp

	// Xử lý Tags trước khi lưu Shelve
	for _, tagName := range request.Tags {
		if tagName == "" {
			continue // Bỏ qua tag rỗng
		}

		if _, exists := tagSet[tagName]; exists {
			continue // Bỏ qua nếu tag đã được xử lý
		}
		tagSet[tagName] = struct{}{}

		existingTag, err := b.FindTagByName(tagName, "shelve")
		if err != nil {
			return models.Shelve{}, fmt.Errorf("failed to find tag: %w", err)
		}

		if existingTag == nil {
			// Chưa có tag, chuẩn bị tạo mới
			newTag := models.Tag{
				EntityType: "shelve",
				Name:       tagName,
			}
			tags = append(tags, newTag)
		} else {
			tags = append(tags, *existingTag)
		}
	}

	// Lưu Shelve cùng Tags
	result.Tags = tags // Gán danh sách tags vào Shelve
	if err := b.DB.Create(&result).Error; err != nil {
		return models.Shelve{}, fmt.Errorf("failed to create shelve: %w", err)
	}

	// Cập nhật EntityID cho các tag mới tạo
	for i := range tags {
		if tags[i].EntityID == 0 { // Chỉ cập nhật nếu tag chưa có EntityID
			tags[i].EntityID = result.ID
			b.DB.Save(&tags[i]) // Lưu lại tag với EntityID đúng
		}
	}

	return result, nil
}

func (b *BookRepositoryImpl) CreateBook(userId int, request request.BookCreateRequest) (models.Book, error) {
	var result models.Book
	err := copier.Copy(&result, request)
	if err != nil {
		return models.Book{}, fmt.Errorf("cant bind request: %w", err)
	}
	result.CreatedBy = uint(userId)

	var tags []models.Tag
	tagSet := make(map[string]struct{}) // Dùng map để tránh trùng lặp

	// Xử lý Tags trước khi lưu Book
	for _, tagReq := range request.Tags {
		if tagReq.Name == "" {
			continue // Bỏ qua tag rỗng
		}

		if _, exists := tagSet[tagReq.Name]; exists {
			continue // Bỏ qua nếu tag đã được xử lý
		}
		tagSet[tagReq.Name] = struct{}{}

		existingTag, err := b.FindTagByName(tagReq.Name, "book")
		if err != nil {
			return models.Book{}, fmt.Errorf("failed to find tag: %w", err)
		}

		if existingTag == nil {
			// Chưa có tag, chuẩn bị tạo mới
			newTag := models.Tag{
				EntityType: "book",
				Name:       tagReq.Name,
				Value:      tagReq.Value,
			}
			tags = append(tags, newTag)
		} else {
			tags = append(tags, *existingTag)
		}
	}

	// Lưu Book cùng Tags
	result.Tags = tags // Gán danh sách tags vào Book
	if err := b.DB.Create(&result).Error; err != nil {
		return models.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	// Cập nhật EntityID cho các tag mới tạo
	for i := range tags {
		if tags[i].EntityID == 0 { // Chỉ cập nhật nếu tag chưa có EntityID
			tags[i].EntityID = result.ID
			b.DB.Save(&tags[i]) // Lưu lại tag với EntityID đúng
		}
	}

	return result, nil
}
