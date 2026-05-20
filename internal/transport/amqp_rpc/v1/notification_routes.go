package v1

import "github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"

func (r *V1) registerNotificationRoutes(routes map[string]server.CallHandler) {
	routes["v1.notification.list"] = r.listNotifications()
	routes["v1.notification.markRead"] = r.markNotificationRead()
}
