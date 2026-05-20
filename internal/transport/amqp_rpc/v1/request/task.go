package request

// CreateTaskReq -.
type CreateTaskReq struct {
	Title       string `json:"title"       validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// GetTaskReq -.
type GetTaskReq struct {
	ID string `json:"id" validate:"required"`
}

// ListTasksReq -.
type ListTasksReq struct {
	Status string `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	Query  string `json:"query"  validate:"max=255"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// UpdateTaskReq -.
type UpdateTaskReq struct {
	ID          string `json:"id"          validate:"required"`
	Title       string `json:"title"       validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// TransitionTaskReq -.
type TransitionTaskReq struct {
	ID     string `json:"id"     validate:"required"`
	Status string `json:"status" validate:"required,oneof=todo in_progress done"`
}

// DeleteTaskReq -.
type DeleteTaskReq struct {
	ID string `json:"id" validate:"required"`
}
