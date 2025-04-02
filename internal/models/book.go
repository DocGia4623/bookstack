package models

import (
	"bookstack/internal/constant"

	"gorm.io/gorm"
)

// Book đại diện cho cuốn sách
type Book struct {
	gorm.Model
	Price       float64   `json:"price"`       // giá sách vật lý
	Title       string    `json:"title"`       // Tiêu đề của sách
	Description string    `json:"description"` // Mô tả của sách
	Slug        string    `json:"slug"`        // Đường dẫn thân thiện
	ShelveID    uint      `json:"shelve_id"`   // Khóa ngoại liên kết đến Shelf
	Shelve      Shelve    `gorm:"foreignKey:ShelveID"`
	Chapters    []Chapter `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE" json:"chapters"` // Danh sách chương của sách
	Tags        []Tag     `gorm:"polymorphic:Entity;polymorphicValue:book" json:"tags"`          // Tags liên kết với sách
	Comments    []Comment `gorm:"polymorphic:Entity;polymorphicValue:book" json:"comments"`
	Restricted  bool      `json:"restricted"` // Trường kiểm soát quyền truy cập
	CreatedBy   uint      `json:"created_by"` // ID của người tạo sách
	UpdatedBy   uint      `json:"updated_by"` // ID của người cập nhật sách
}

// Chapter đại diện cho chương của một cuốn sách
type Chapter struct {
	gorm.Model
	Title      string `json:"title"`                                                         // Tiêu đề chương
	Order      int    `json:"order"`                                                         // Thứ tự sắp xếp chương
	BookID     uint   `json:"book_id"`                                                       // Khóa ngoại liên kết đến Book
	Pages      []Page `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE" json:"pages"` // Danh sách trang của chương
	Restricted bool   `json:"restricted"`                                                    // Quyền truy cập chương
}

// Page đại diện cho trang trong một chương hoặc sách
type Page struct {
	gorm.Model
	Title         string         `json:"title"`                    // Tiêu đề trang
	Slug          string         `json:"slug"`                     // Đường dẫn thân thiện
	Content       string         `json:"content" gorm:"type:text"` // Nội dung trang (có thể là markdown, HTML,...)
	Order         int            `json:"order"`                    // Thứ tự sắp xếp trang trong chương
	ChapterID     uint           `json:"chapter_id"`               // Khóa ngoại liên kết đến Chapter
	Chapter       Chapter        `gorm:"foreignKey:ChapterID"`
	Restricted    bool           `json:"restricted"` // Quyền truy cập trang
	PageRevisions []PageRevision `gorm:"foreignKey:PageId"`
}

type PageRevision struct {
	gorm.Model
	PageId         int    `json:"page_id"`
	Content        string `json:"content"`
	RevisionNumber int    `json:"revision_number"`
}

// Shelf đại diện cho kệ chứa sách, dùng để phân loại sách theo chủ đề hoặc danh mục
type Shelve struct {
	gorm.Model
	Name        string    `json:"name"`                                                   // Tên kệ
	Description string    `json:"description"`                                            // Mô tả kệ
	Order       int       `json:"order"`                                                  // Thứ tự hiển thị của kệ
	Books       []Book    `gorm:"foreignKey:ShelveID" json:"books"`                       // Danh sách sách trong kệ
	Tags        []Tag     `gorm:"polymorphic:Entity;polymorphicValue:shelve" json:"tags"` // Tags liên kết với kệ
	Comments    []Comment `gorm:"polymorphic:Entity;polymorphicValue:shelve" json:"comments"`
	CreatedBy   uint      `json:"created_by"` // ID của người tạo kệ
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

// Comment của sách và kệ sách

type Comment struct {
	gorm.Model
	EntityID   uint   `json:"entity_id"`   // ID của entity (Book hoặc Shelve)
	EntityType string `json:"entity_type"` // Loại entity (book, shelve)
	Rating     int    `json:"rating"`
	Text       string `json:"text"`
	CreatedBy  int    `json:"created_by"`
	ParentID   int    `json:"parent_id"`
}

// Order - Đơn hàng
type Order struct {
	gorm.Model
	UserID      uint                 `json:"user_id"`       // Người đặt hàng
	TotalPrice  float64              `json:"total_price"`   // Tổng giá trị đơn hàng
	Status      constant.OrderStatus `json:"status"`        // Trạng thái đơn hàng: pending, completed, cancelled
	OrderDetail []OrderDetail        `json:"order_details"` // Danh sách sách trong đơn
	Address     string               `json:"address"`       // Địa chỉ giao hàng
	Phone       string               `json:"phone"`         // Số điện thoại liên hệ
}

// OrderDetail - Chi tiết đơn hàng
type OrderDetail struct {
	gorm.Model
	OrderID  uint    `json:"order_id"`
	Order    Order   `gorm:"foreignKey:OrderID"`
	BookID   uint    `json:"book_id"`
	Book     Book    `gorm:"foreignKey:BookID"`
	Quantity int     `json:"quantity"` // Số lượng sách đặt
	Price    float64 `json:"price"`    // Giá sách tại thời điểm đặt hàng
}
