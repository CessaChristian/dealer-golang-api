package auth

type RegisterRequest struct {
	Name string `json:"name" validate:"required,min=2,max=150"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
