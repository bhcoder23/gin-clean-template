package request

// RegisterReq -.
type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginReq -.
type LoginReq struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
