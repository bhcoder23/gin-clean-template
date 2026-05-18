package request

// Translate -.
type Translate struct {
	Source      string `example:"auto"                               json:"source"      validate:"required"`
	Destination string `example:"en"                                 json:"destination" validate:"required"`
	Original    string `example:"текст для перевода"                 json:"original"    validate:"required"`
} // @name v1.Translate
