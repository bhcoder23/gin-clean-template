package request

// RegisterReq -.
type RegisterReq struct {
	Username string `example:"johndoe"          json:"username" validate:"required,min=3,max=255"`
	Email    string `example:"john@example.com" json:"email"    validate:"required,email"`
	Password string `example:"secret123"        json:"password" validate:"required,min=6"`
} // @name v1.RegisterReq

// LoginReq -.
type LoginReq struct {
	Email    string `example:"john@example.com" json:"email"    validate:"required,email"`
	Password string `example:"secret123"        json:"password" validate:"required"`
} // @name v1.LoginReq
