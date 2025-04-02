package response

type UserResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"name"`
	Email    string `json:"email"`
}

type UserRegisterResponse struct {
	ID             int    `json:"id"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	RememberToken  string `json:"remember_token"`
	EmailConfirmed bool   `json:"email_confirmed"`
	ImageId        int    `json:"image_id"`
}
