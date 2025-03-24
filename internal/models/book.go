package models

import (
	"gorm.io/gorm"
)

// Book đại diện cho cuốn sách
type Book struct {
	gorm.Model
	Title       string    `json:"title"`       // Tiêu đề của sách
	Description string    `json:"description"` // Mô tả của sách
	Slug        string    `json:"slug"`        // Đường dẫn thân thiện
	ShelfID     uint      `json:"shelf_id"`    // Khóa ngoại liên kết đến Shelf
	Shelf       Shelf     `gorm:"foreignKey:ShelfID"`
	Chapters    []Chapter `json:"chapters"`                                             // Danh sách chương của sách
	Tags        []Tag     `gorm:"polymorphic:Entity;polymorphicValue:book" json:"tags"` // Tags liên kết với sách
	Restricted  bool      `json:"restricted"`                                           // Trường kiểm soát quyền truy cập
	CreatedBy   uint      `json:"created_by"`                                           // ID của người tạo sách
	UpdatedBy   uint      `json:"updated_by"`                                           // ID của người cập nhật sách
}

// Chapter đại diện cho chương của một cuốn sách
type Chapter struct {
	gorm.Model
	Title      string `json:"title"`      // Tiêu đề chương
	Order      int    `json:"order"`      // Thứ tự sắp xếp chương
	BookID     uint   `json:"book_id"`    // Khóa ngoại liên kết đến Book
	Pages      []Page `json:"pages"`      // Danh sách trang của chương
	Restricted bool   `json:"restricted"` // Quyền truy cập chương
}

// Page đại diện cho trang trong một chương hoặc sách
type Page struct {
	gorm.Model
	Title      string  `json:"title"`                    // Tiêu đề trang
	Slug       string  `json:"slug"`                     // Đường dẫn thân thiện
	Content    string  `json:"content" gorm:"type:text"` // Nội dung trang (có thể là markdown, HTML,...)
	Order      int     `json:"order"`                    // Thứ tự sắp xếp trang trong chương
	ChapterID  uint    `json:"chapter_id"`               // Khóa ngoại liên kết đến Chapter
	Chapter    Chapter `gorm:"foreignKey:ChapterID"`
	Restricted bool    `json:"restricted"` // Quyền truy cập trang
}

// Shelf đại diện cho kệ chứa sách, dùng để phân loại sách theo chủ đề hoặc danh mục
type Shelf struct {
	gorm.Model
	Name        string `json:"name"`                                                  // Tên kệ
	Description string `json:"description"`                                           // Mô tả kệ
	Order       int    `json:"order"`                                                 // Thứ tự hiển thị của kệ
	Books       []Book `json:"books"`                                                 // Danh sách sách trong kệ
	Tags        []Tag  `gorm:"polymorphic:Entity;polymorphicValue:shelf" json:"tags"` // Tags liên kết với kệ
	CreatedBy   uint   `json:"created_by"`                                            // ID của người tạo kệ
	UpdatedBy   uint   `json:"updated_by"`                                            // ID của người cập nhật kệ
}

// Tag dùng để gắn nhãn mô tả, từ khóa cho các entity (Book, Chapter, Page,...)
type Tag struct {
	gorm.Model
	EntityID   uint   `json:"entity_id"`   // ID của entity được gắn tag
	EntityType string `json:"entity_type"` // Loại của entity (book, chapter, page, ...)
	Name       string `json:"name"`        // Tên của tag
	Value      string `json:"value"`       // Giá trị của tag
	Order      int    `json:"order"`       // Thứ tự sắp xếp nếu cần
}
