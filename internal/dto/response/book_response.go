package response

type BookResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title" binding:"required"`     // Tiêu đề của sách (bắt buộc)
	Description string        `json:"description"`                  // Mô tả của sách
	Slug        string        `json:"slug"`                         // Đường dẫn thân thiện
	ShelveID    uint          `json:"shelve_id" binding:"required"` // ID của kệ sách chứa nó (bắt buộc)
	Restricted  bool          `json:"restricted"`                   // Trạng thái kiểm soát quyền truy cập
	CreatedBy   string        `json:"created_by"`                   // Name của người tạo sách
	Tags        []TagResponse `json:"tags"`
	Shelve      string        `json:"shelve"`
}

type ShelveResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name" binding:"required"`       // Tên kệ (bắt buộc)
	Description string        `json:"description"`                   // Mô tả kệ
	Order       int           `json:"order"`                         // Thứ tự hiển thị của kệ
	Tags        []TagResponse `json:"tags"`                          // Danh sách tag (chỉ lấy tên tag)
	CreatedBy   string        `json:"created_by" binding:"required"` // người tạo kệ (bắt buộc)
}

type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
