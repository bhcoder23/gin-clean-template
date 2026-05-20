package v1

import "github.com/bhcoder23/gin-clean-template/pkg/rabbitmq/rmq_rpc/server"

func (r *V1) registerAuthRoutes(routes map[string]server.CallHandler) {
	routes["v1.auth.register"] = r.register()
	routes["v1.auth.login"] = r.login()
}
