package response

type OrderResponse struct {
	OrderID     uint                  `json:"order_id"`
	UserID      uint                  `json:"user_id"`
	Status      string                `json:"status"`
	TotalPrice  float64               `json:"total_price"`
	Address     string                `json:"address"` // Địa chỉ giao hàng
	Phone       string                `json:"phone"`   // Số điện thoại
	ShiperID    uint                  `json:"shiper_id"`
	OrderDetail []OrderDetailResponse `json:"order_detail"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

type OrderDetailResponse struct {
	Book     BookOrderResponse `json:"book"`
	Quantity int               `json:"quantity"`
	Price    float64           `json:"price"`
}

type BookOrderResponse struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	Slug      string  `json:"slug"`
	CreatedBy string  `json:"created_by"`
}
