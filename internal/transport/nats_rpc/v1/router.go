package v1

import (
	"github.com/bhcoder23/gin-clean-template/internal/usecase"
	"github.com/bhcoder23/gin-clean-template/pkg/jwt"
	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/bhcoder23/gin-clean-template/pkg/nats/nats_rpc/server"
	"github.com/go-playground/validator/v10"
)

// NewRoutes -.
func NewRoutes(routes map[string]server.CallHandler, n usecase.Notification, u usecase.User, tk usecase.Task, j *jwt.Manager, l logger.Interface) {
	r := &V1{n: n, u: u, tk: tk, j: j, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	routes["v1.auth.register"] = r.register()
	routes["v1.auth.login"] = r.login()

	routes["v1.notification.list"] = r.listNotifications()
	routes["v1.notification.markRead"] = r.markNotificationRead()

	routes["v1.task.create"] = r.createTask()
	routes["v1.task.get"] = r.getTask()
	routes["v1.task.list"] = r.listTasks()
	routes["v1.task.update"] = r.updateTask()
	routes["v1.task.transition"] = r.transitionTask()
	routes["v1.task.delete"] = r.deleteTask()
}
