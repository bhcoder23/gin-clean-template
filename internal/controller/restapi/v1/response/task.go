package response

import "github.com/bhcoder23/gin-clean-template/internal/entity"

// TaskList -.
type TaskList struct {
	Tasks []entity.Task `json:"tasks"`
	Total int           `example:"42" json:"total"`
} // @name v1.TaskList
