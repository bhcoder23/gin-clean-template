package response

import "github.com/bhcoder23/gin-clean-template/internal/domain"

// TaskList -.
type TaskList struct {
	Tasks []domain.Task `json:"tasks"`
	Total int           `json:"total"`
}

// DeleteStatus -.
type DeleteStatus struct {
	Status string `json:"status"`
}
