package repository

import (
	"bookstack/internal/dto/request"
	"bookstack/internal/models"
	"fmt"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type BookRepository interface {
	//book
	CreateCompleteBook(int, request.CompleteBookCreateRequest) (models.Book, error)
	CreateBook(int, request.BookCreateRequest) (models.Book, error)
	GetAllBook() ([]models.Book, error)
	UpdateBook(int, request.BookCreateRequest) (models.Book, error)
	DeleteBook(int) error
	//shelve
	CreateShelve(int, request.ShelveCreateRequest) (models.Shelve, error)
	GetShelves() ([]models.Shelve, error)
	DeleteShelve(int) error
	//chapter
	CreateChapter(uint, request.BookChapterRequest) (models.Chapter, error)
	GetChaptersOfBook(int) ([]models.Chapter, error)
	DeleteChapter(int) error
	UpdateChapter(int, request.BookChapterRequest) (models.Chapter, error)
	//page
	AddPage(uint, request.PageRequest) (models.Page, error)
	GetPageChapter(int) ([]models.Page, error)
	DeletePage(int) error
	UpdatePage(int, request.PageRequest) (models.Page, error)
}

type BookRepositoryImpl struct {
	DB *gorm.DB
}

func NewBookRepositoryImpl(Db *gorm.DB) BookRepository {
	return &BookRepositoryImpl{
		DB: Db,
	}
}

func (b *BookRepositoryImpl) UpdateChapter(chapterId int, request request.BookChapterRequest) (models.Chapter, error) {
	var chapter models.Chapter
	err := b.DB.Where("id = ?", chapterId).First(&chapter).Error
	if err != nil {
		return models.Chapter{}, err
	}
	copier.Copy(&chapter, request)
	err = b.DB.Save(&chapter).Error
	if err != nil {
		return models.Chapter{}, err
	}
	return chapter, nil
}

func (b *BookRepositoryImpl) UpdatePage(pageId int, request request.PageRequest) (models.Page, error) {
	var page models.Page

	// Tìm page theo ID
	if err := b.DB.Where("id = ?", pageId).First(&page).Error; err != nil {
		return models.Page{}, err // Trả về lỗi nếu không tìm thấy
	}

	// Cập nhật các trường cụ thể từ request
	updates := map[string]interface{}{}
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.Content != "" {
		updates["content"] = request.Content
	}

	// Nếu không có gì để cập nhật
	if len(updates) == 0 {
		return page, nil
	}

	// Cập nhật dữ liệu
	if err := b.DB.Model(&page).Updates(updates).Error; err != nil {
		return models.Page{}, err
	}

	return page, nil
}

func (b *BookRepositoryImpl) DeleteChapter(chapterId int) error {
	var chapter models.Chapter
	err := b.DB.Where("id = ?", chapterId).Delete(&chapter).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *BookRepositoryImpl) DeletePage(pageId int) error {
	var page models.Page
	err := b.DB.Where("id = ?", pageId).Delete(&page).Error
	if err != nil {
		return err
	}
	return nil
}
func (b *BookRepositoryImpl) DeleteShelve(shelveId int) error {
	var shelve models.Shelve
	err := b.DB.Where("id = ?", shelveId).Delete(&shelve).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *BookRepositoryImpl) UpdateBook(bookId int, request request.BookCreateRequest) (models.Book, error) {
	var book models.Book
	err := b.DB.Where("id = ?", bookId).First(&book).Error
	if err != nil {
		return models.Book{}, err
	}

	// Copy các trường cơ bản, không bao gồm CreatedBy
	book.Title = request.Title
	book.Description = request.Description
	book.Slug = request.Slug
	book.Restricted = request.Restricted
	book.Price = request.Price

	// Chỉ cập nhật ShelveID nếu được cung cấp trong request
	if request.ShelveID > 0 {
		// Kiểm tra xem kệ sách có tồn tại không
		var shelve models.Shelve
		err = b.DB.First(&shelve, request.ShelveID).Error
		if err != nil {
			return models.Book{}, fmt.Errorf("shelve not found: %w", err)
		}
		book.ShelveID = request.ShelveID
	}

	// Cập nhật tags
	var tags []models.Tag
	tagSet := make(map[string]struct{})
	for _, tagReq := range request.Tags {
		if tagReq.Name == "" {
			continue
		}
		if _, exists := tagSet[tagReq.Name]; exists {
			continue
		}
		tagSet[tagReq.Name] = struct{}{}

		existingTag, err := b.FindTagByName(tagReq.Name, "book")
		if err != nil {
			return models.Book{}, fmt.Errorf("failed to find tag: %w", err)
		}

		if existingTag == nil {
			tags = append(tags, models.Tag{
				EntityType: "book",
				Name:       tagReq.Name,
				Value:      tagReq.Value,
			})
		} else {
			tags = append(tags, *existingTag)
		}
	}
	book.Tags = tags

	// Lưu sách với các trường đã cập nhật
	err = b.DB.Model(&book).Updates(map[string]interface{}{
		"title":       book.Title,
		"description": book.Description,
		"slug":        book.Slug,
		"restricted":  book.Restricted,
		"price":       book.Price,
		"shelve_id":   book.ShelveID,
	}).Error
	if err != nil {
		return models.Book{}, err
	}

	// Cập nhật EntityID cho các tag mới
	for i := range tags {
		if tags[i].EntityID == 0 {
			tags[i].EntityID = book.ID
			b.DB.Save(&tags[i])
		}
	}

	return book, nil
}

func (b *BookRepositoryImpl) DeleteBook(bookId int) error {
	var book models.Book
	err := b.DB.Where("id = ?", bookId).Delete(&book).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *BookRepositoryImpl) GetShelves() ([]models.Shelve, error) {
	var shelves []models.Shelve
	err := b.DB.
		Preload("CreatedBy").
		Preload("Books").
		Preload("Tags").
		Find(&shelves).Error
	if err != nil {
		return nil, err
	}
	return shelves, nil
}

func (b *BookRepositoryImpl) CreateCompleteBook(userId int, req request.CompleteBookCreateRequest) (models.Book, error) {
	// Tạo sách
	book, err := b.CreateBook(userId, req.BookCreateRequest)
	if err != nil {
		return models.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	// Liên kết với Shelve nếu có
	if req.BookCreateRequest.ShelveID > 0 {
		var shelve models.Shelve
		if err := b.DB.First(&shelve, req.BookCreateRequest.ShelveID).Error; err != nil {
			return models.Book{}, fmt.Errorf("failed to find shelve: %w", err)
		}
		book.Shelve = shelve
		if err := b.DB.Save(&book).Error; err != nil {
			return models.Book{}, fmt.Errorf("failed to update book shelve: %w", err)
		}
	}

	// Tạo danh sách chương cho sách
	var chapters []models.Chapter
	for _, chapterReq := range req.BookChapterRequest {
		chapter, err := b.CreateChapter(book.ID, chapterReq)
		if err != nil {
			return models.Book{}, fmt.Errorf("failed to create chapter: %w", err)
		}
		chapters = append(chapters, chapter)
	}

	// tạo trang
	for _, pageReq := range req.PageRequest {
		_, err := b.CreatePage(pageReq)
		if err != nil {
			return models.Book{}, fmt.Errorf("cant create page: %w", err)
		}
	}

	// Load lại sách với các chương đã thêm vào
	book.Chapters = chapters
	return book, nil
}

func (b *BookRepositoryImpl) CreatePage(request request.PageRequest) (models.Page, error) {
	var page models.Page
	err := copier.Copy(&page, request)
	if err != nil {
		return models.Page{}, nil
	}
	err = b.DB.FirstOrCreate(&page, models.Page{
		Title:     request.Title,
		ChapterID: uint(request.ChapterId),
	}).Error
	if err != nil {
		return models.Page{}, err
	}
	return page, err
}

func (b *BookRepositoryImpl) GetPageChapter(chapterId int) ([]models.Page, error) {
	var pages []models.Page
	err := b.DB.Where("chapter_id = ?", chapterId).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (b *BookRepositoryImpl) AddPage(chapterId uint, request request.PageRequest) (models.Page, error) {
	var page models.Page
	err := copier.Copy(&page, request)
	if err != nil {
		return models.Page{}, nil
	}
	page.ChapterID = chapterId
	err = b.DB.FirstOrCreate(&page, models.Page{
		Title:     request.Title,
		ChapterID: chapterId,
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
func (b *BookRepositoryImpl) CreateChapter(bookId uint, request request.BookChapterRequest) (models.Chapter, error) {
	var chapter models.Chapter
	err := copier.Copy(&chapter, request)
	if err != nil {
		return models.Chapter{}, err
	}
	chapter.BookID = bookId
	// Dùng FirstOrCreate để tìm hoặc tạo mới Chapter
	result := b.DB.FirstOrCreate(&chapter, models.Chapter{
		Title:  request.Title,
		BookID: bookId, // Điều kiện để xác định chương có tồn tại hay chưa
	})
	if result.Error != nil {
		return models.Chapter{}, result.Error
	}
	return chapter, nil

}
func (b *BookRepositoryImpl) GetAllBook() ([]models.Book, error) {
	var books []models.Book
	err := b.DB.Preload("Chapters").Preload("Tags").Preload("Shelve").Find(&books).Error
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
