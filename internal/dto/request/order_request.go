package request

type OrderDetailRequest struct {
	BookID   uint `json:"book_id" binding:"required"`  // ID của sách
	Quantity int  `json:"quantity" binding:"required"` // Số lượng sách đặt
}

type OrderRequest struct {
	OrderDetails []OrderDetailRequest `json:"order_details" binding:"required"` // Danh sách sách trong đơn
	Address      string               `json:"address"`                          // Địa chỉ giao hàng
	Phone        string               `json:"phone"`                            // Số điện thoại
}
