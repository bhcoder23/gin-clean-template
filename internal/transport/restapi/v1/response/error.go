package response

// Error -.
type Error struct {
	Error ErrorBody `json:"error"`
} // @name v1.Error

// ErrorBody is the stable client-facing error envelope.
type ErrorBody struct {
	Code      string `example:"TASK_NOT_FOUND" json:"code"`
	Message   string `example:"task not found"  json:"message"`
	RequestID string `example:"request-id"      json:"request_id,omitempty"`
} // @name v1.ErrorBody
