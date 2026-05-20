package request

// CreateTaskReq -.
type CreateTaskReq struct {
	Title       string `example:"My task"          json:"title"       validate:"required,max=255"`
	Description string `example:"Task description" json:"description" validate:"max=1000"`
} // @name v1.CreateTaskReq

// UpdateTaskReq -.
type UpdateTaskReq struct {
	Title       string `example:"Updated task"        json:"title"       validate:"required,max=255"`
	Description string `example:"Updated description" json:"description" validate:"max=1000"`
} // @name v1.UpdateTaskReq

// TransitionTaskReq -.
type TransitionTaskReq struct {
	Status string `example:"in_progress" json:"status" validate:"required,oneof=todo in_progress done"`
} // @name v1.TransitionTaskReq
