package response

type UserResponse struct {
	FullName string `json:"name"`
	Email    string `json:"email"`
}
