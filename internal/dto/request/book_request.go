package request

type BookCreateRequest struct {
	Title       string       `json:"title" binding:"required"`     // Tiêu đề của sách (bắt buộc)
	Description string       `json:"description"`                  // Mô tả của sách
	Slug        string       `json:"slug"`                         // Đường dẫn thân thiện
	ShelveID    uint         `json:"shelve_id" binding:"required"` // ID của kệ sách chứa nó (bắt buộc)
	Restricted  bool         `json:"restricted"`                   // Trạng thái kiểm soát quyền truy cập
	CreatedBy   uint         `json:"created_by"`                   // ID của người tạo sách
	Tags        []TagRequest `json:"tags"`                         // Danh sách tag
}

type TagRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
