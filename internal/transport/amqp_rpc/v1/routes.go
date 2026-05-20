package v1

import "github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"

func (r *V1) registerAuthRoutes(routes map[string]server.CallHandler) {
	routes["v1.auth.register"] = r.register()
	routes["v1.auth.login"] = r.login()
}

func (r *V1) registerNotificationRoutes(routes map[string]server.CallHandler) {
	routes["v1.notification.list"] = r.listNotifications()
	routes["v1.notification.markRead"] = r.markNotificationRead()
}

func (r *V1) registerTaskRoutes(routes map[string]server.CallHandler) {
	routes["v1.task.create"] = r.createTask()
	routes["v1.task.get"] = r.getTask()
	routes["v1.task.list"] = r.listTasks()
	routes["v1.task.update"] = r.updateTask()
	routes["v1.task.transition"] = r.transitionTask()
	routes["v1.task.delete"] = r.deleteTask()
}
