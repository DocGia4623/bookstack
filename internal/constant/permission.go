package constant

// User Permissions
const (
	ReadUser   = "read:user"
	WriteUser  = "write:user"
	DeleteUser = "delete:user"
)

// Admin Permissions
const (
	ManageUsers = "manage:users"
	ManageRoles = "manage:roles"
)

// Shipper Permissions
const (
	ReceiveOrder      = "receive:order"
	UpdateOrderStatus = "update:order:status"
)
