package request

// Register -.
type Register struct {
	Username string `example:"johndoe"          json:"username" validate:"required,min=3,max=255"`
	Email    string `example:"john@example.com" json:"email"    validate:"required,email"`
	Password string `example:"secret123"        json:"password" validate:"required,min=6"`
} // @name v1.Register

// Login -.
type Login struct {
	Email    string `example:"john@example.com" json:"email"    validate:"required,email"`
	Password string `example:"secret123"        json:"password" validate:"required"`
} // @name v1.Login
