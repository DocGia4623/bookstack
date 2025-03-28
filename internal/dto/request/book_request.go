package request

type CompleteBookCreateRequest struct {
	BookCreateRequest  BookCreateRequest    `json:"book"`
	BookChapterRequest []BookChapterRequest `json:"chapters"`
	PageRequest        []PageRequest        `json:"pages"`
}

type BookCreateRequest struct {
	Title       string       `json:"title" binding:"required"`     // Tiêu đề của sách (bắt buộc)
	Description string       `json:"description"`                  // Mô tả của sách
	Slug        string       `json:"slug"`                         // Đường dẫn thân thiện
	ShelveID    uint         `json:"shelve_id" binding:"required"` // ID của kệ sách chứa nó (bắt buộc)
	Restricted  bool         `json:"restricted"`                   // Trạng thái kiểm soát quyền truy cập
	CreatedBy   uint         `json:"created_by"`                   // ID của người tạo sách
	Price       float64      `json:"price"`
	Tags        []TagRequest `json:"tags"` // Danh sách tag
}

type TagRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ShelveCreateRequest struct {
	Name        string   `json:"name" binding:"required"`       // Tên kệ (bắt buộc)
	Description string   `json:"description"`                   // Mô tả kệ
	Order       int      `json:"order"`                         // Thứ tự hiển thị của kệ
	Tags        []string `json:"tags"`                          // Danh sách tag (chỉ lấy tên tag)
	CreatedBy   uint     `json:"created_by" binding:"required"` // ID người tạo kệ (bắt buộc)
}

type BookChapterRequest struct {
	Title string `json:"title"` // Tiêu đề chương
	Order int    `json:"order"` // Thứ tự sắp xếp chương
}

type PageRequest struct {
	Title     string `json:"title" binding:"required"` // Tiêu đề trang (bắt buộc)
	Slug      string `json:"slug"`                     // Đường dẫn thân thiện (có thể tự động tạo nếu rỗng)
	Content   string `json:"content"`                  // Nội dung trang (markdown, HTML,...)
	Order     int    `json:"order"`                    // Thứ tự sắp xếp trang
	ChapterId int    `json:"chapter_id"`
}

type CommentRequest struct {
}
