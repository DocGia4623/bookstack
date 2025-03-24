package request

type UserCreateRequest struct {
	FullName string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	FullName string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
