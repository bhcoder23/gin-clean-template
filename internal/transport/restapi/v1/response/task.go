package response

import "github.com/bhcoder23/gin-clean-template/internal/domain"

// TaskList -.
type TaskList struct {
	Tasks []domain.Task `json:"tasks"`
	Total int           `example:"42" json:"total"`
} // @name v1.TaskList
