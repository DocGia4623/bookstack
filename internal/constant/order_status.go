package constant

type OrderStatus int

const (
	Pending    OrderStatus = iota // 0 - Chờ xác nhận
	Confirmed                     // 1 - Đã xác nhận
	Processing                    // 2 - Đang xử lý
	Shipped                       // 3 - Đã giao hàng
	Delivered                     // 4 - Đã giao thành công
	Cancelled                     // 5 - Đã hủy
	Returned                      // 6 - Đã trả hàng
	Failed                        // 7 - Thất bại
)
