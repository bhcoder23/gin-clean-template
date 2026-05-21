# Gin 整洁模板

[🇺🇸 English](README.md)

通用的 Go 后端整洁架构模板，由 `bhcoder23` 维护。

[![License](https://img.shields.io/badge/License-MIT-success)](LICENSE)
[![Maintainer](https://img.shields.io/badge/Maintainer-bhcoder23-1f6feb)](https://github.com/bhcoder23)

[![Web Framework](https://img.shields.io/badge/Gin-Web%20Framework-blue)](https://github.com/gin-gonic/gin)
[![API Documentation](https://img.shields.io/badge/Swagger-API%20Documentation-blue)](https://github.com/swaggo/swag)
[![Validation](https://img.shields.io/badge/Validator-Data%20Integrity-blue)](https://github.com/go-playground/validator)
[![JSON Handling](https://img.shields.io/badge/Go--JSON-Fast%20Serialization-blue)](https://github.com/goccy/go-json)
[![Persistence](https://img.shields.io/badge/sqlc-Type--safe%20SQL-blue)](https://sqlc.dev/)
[![Database Migrations](https://img.shields.io/badge/Migrations-Seamless%20Schema%20Updates-blue)](https://github.com/golang-migrate/migrate)
[![Logging](https://img.shields.io/badge/ZeroLog-Structured%20Logging-blue)](https://github.com/rs/zerolog)
[![Metrics](https://img.shields.io/badge/Prometheus-Metrics%20Integration-blue)](https://github.com/prometheus/client_golang)
[![Testing](https://img.shields.io/badge/Testify-Testing%20Framework-blue)](https://github.com/stretchr/testify)
[![Mocking](https://img.shields.io/badge/Mock-Mocking%20Library-blue)](https://go.uber.org/mock)

## 综述

这个仓库是 `bhcoder23` 维护的 Gin 后端脚手架。

这个模板关注的是项目离开 demo 阶段以后仍然可维护的组织方式：

- 让 domain 和 usecase 独立于 HTTP、gRPC、消息队列和数据库驱动
- 让 transport adapter 保持薄、明确、可替换
- 显式处理事务、错误码、request ID、日志、trace 和 outbox 边界
- 让派生项目可以轻松删除不需要的可选 adapter

参考的原始项目（MIT 协议）：
- [evrone/go-clean-template](https://github.com/evrone/go-clean-template)

此模板是一个应用进程，外挂多种传输适配器：

- AMQP RPC（基于 RabbitMQ 作为传输）
- NATS RPC（基于 NATS 作为传输）
- gRPC（基于 protobuf 的 [gRPC](https://grpc.io/) 框架）
- REST API（基于 [Gin](https://github.com/gin-gonic/gin) 框架）

默认本地开发路径只启动 HTTP。其他 transport 仍然作为可选 adapter 保留，派生项目可以按需打开，不必一开始就背上所有依赖。

模板包含三个领域，用于演示多服务架构。
它们是脚手架示例领域，并不是必须保留的产品边界：

- **用户认证** — 注册、登录、基于 JWT 的授权
- **任务管理** — CRUD 操作，支持状态转换（todo、in_progress、done）
- **通知中心** — 任务活动通知与已读追踪

这些示例领域可以通过四种传输协议（REST、gRPC、AMQP RPC、NATS RPC）暴露，但派生项目通常应该删掉不需要的 adapter。

## 内容

- [从这里开始](#从这里开始)
- [演示链路](#演示链路)
- [领域](#领域)
- [快速开始](#快速开始)
- [工程架构](#工程架构)
- [新增业务模块](docs/add-new-module.md)
- [依赖注入](#依赖注入)
- [整洁架构](#整洁架构)

## 从这里开始

第一次体验建议走 HTTP-first 路径，这样更容易看清脚手架骨架，也方便后续裁剪：

```sh
# 启动 PostgreSQL，走默认 HTTP-first 本地路径
make compose-up

# 执行迁移并启动当前启用的 transport
make run
```

如果要一次查看所有演示 adapter，先启动可选 broker，再使用 `make run-all-transports`。

```sh
make compose-up-adapters
make run-all-transports
```

服务起来后，最快理解这个脚手架的方式，就是顺着一条完整的 REST 业务链路走一遍。

## 演示链路

注册用户：

```sh
curl -s http://127.0.0.1:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"johndoe","email":"john@example.com","password":"secret123"}'
```

登录并取回 JWT：

```sh
TOKEN=$(
  curl -s http://127.0.0.1:8080/v1/auth/login \
    -H 'Content-Type: application/json' \
    -d '{"email":"john@example.com","password":"secret123"}' | jq -r '.token'
)
```

读取当前用户资料：

```sh
curl -s http://127.0.0.1:8080/v1/user/profile \
  -H "Authorization: Bearer $TOKEN"
```

创建任务：

```sh
curl -s http://127.0.0.1:8080/v1/tasks \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Ship the scaffold","description":"Exercise the happy path"}'
```

查看任务列表：

```sh
curl -s 'http://127.0.0.1:8080/v1/tasks?limit=10&offset=0' \
  -H "Authorization: Bearer $TOKEN"
```

查看任务流转产生的未读通知：

```sh
curl -s 'http://127.0.0.1:8080/v1/notifications?unread_only=true&limit=10&offset=0' \
  -H "Authorization: Bearer $TOKEN"
```

## 领域

模板包含三个完整实现的领域，并分别演示如何挂接到这些传输适配器上。

### 用户认证

注册、登录和基于 JWT 的授权。

| 操作   | REST                     | gRPC                     |
|------|--------------------------|--------------------------|
| 注册   | `POST /v1/auth/register` | `AuthService/Register`   |
| 登录   | `POST /v1/auth/login`    | `AuthService/Login`      |
| 获取资料 | `GET /v1/user/profile`   | `AuthService/GetProfile` |

- 密码使用 bcrypt 加密
- JWT 令牌支持可配置的过期时间
- 所有传输协议均有认证中间件

### 任务管理

CRUD 操作，支持状态机。

| 操作   | REST                         | gRPC                         |
|------|------------------------------|------------------------------|
| 创建   | `POST /v1/tasks`             | `TaskService/CreateTask`     |
| 列表   | `GET /v1/tasks`              | `TaskService/ListTasks`      |
| 获取   | `GET /v1/tasks/:id`          | `TaskService/GetTask`        |
| 更新   | `PUT /v1/tasks/:id`          | `TaskService/UpdateTask`     |
| 状态转换 | `PATCH /v1/tasks/:id/status` | `TaskService/TransitionTask` |
| 删除   | `DELETE /v1/tasks/:id`       | `TaskService/DeleteTask`     |

- 状态转换：`todo` → `in_progress` → `done`（以及 `in_progress` → `todo`）
- 支持 `limit`/`offset` 分页和可选状态过滤
- 任务绑定到已认证的用户

### 通知中心

任务活动通知会持久化到 PostgreSQL，并通过所有 transport 暴露出来。

| 操作 | REST | gRPC |
|----|------|------|
| 列表 | `GET /v1/notifications` | `NotificationService/ListNotifications` |
| 标记已读 | `PATCH /v1/notifications/:id/read` | `NotificationService/MarkNotificationRead` |

- 任务创建和状态流转都会生成通知
- 支持 `unread_only=true` 过滤未读通知
- 已读通知会记录 `read_at`

## 快速开始

### 本地开发

Docker 不是必选项。`.env.example` 默认只开启 HTTP；gRPC、RabbitMQ RPC、NATS RPC 都是可选项。Docker Compose 演示栈会在需要完整 adapter 集合时显式打开这些开关。

```sh
# 默认 HTTP-first 路径只需要 PostgreSQL
make compose-up
# 执行迁移并启动应用
make run
```

如果你想忽略当前 `.env`，强制把所有演示 transport 都跑起来，可以直接执行：

```sh
make compose-up-adapters
make run-all-transports
```

### 集成测试

```sh
# DB、app + migrations、integration tests.
# 这些测试需要 integration build tag，通常由 Jenkins 或其他流水线执行。
make compose-up-integration-test
```

### 带反向代理的完整 Docker 栈

```sh
make compose-up-all
```

完整演示栈中的服务：

- AMQP RPC:
  - URL: `amqp://guest:guest@127.0.0.1:5672/`
  - Client Exchange: `rpc_client`
  - Server Exchange: `rpc_server`
- NATS RPC:
  - URL: `nats://guest:guest@127.0.0.1:4222/`
  - Server Exchange: `rpc_server`
- REST API:
  - http://app.lvh.me/healthz | http://127.0.0.1:8080/healthz
  - http://app.lvh.me/readyz | http://127.0.0.1:8080/readyz
  - http://app.lvh.me/metrics | http://127.0.0.1:8080/metrics
  - http://app.lvh.me/swagger | http://127.0.0.1:8080/swagger
- gRPC:
  - URL: `tcp://grpc.lvh.me:8081` | `tcp://127.0.0.1:8081`
  - [v1/auth.proto](docs/proto/v1/auth.proto)
  - [v1/task.proto](docs/proto/v1/task.proto)
  - [v1/notification.proto](docs/proto/v1/notification.proto)
- PostgreSQL:
  - `postgres://user:myAwEsOm3pa55@w0rd@127.0.0.1:5432/db`
- RabbitMQ:
  - http://rabbitmq.lvh.me | http://127.0.0.1:15672
  - Credentials: `guest` / `guest`
- NATS monitoring:
  - http://nats.lvh.me | http://127.0.0.1:8222/
  - Credentials: `guest` / `guest`

## 工程架构

### `cmd/app/main.go`

配置和日志能力初始化。主要的功能在 `internal/app/app.go`

### `config`

12-Factor推荐将应用的配置存储于 环境变量 中（ `env vars`, `env` ）。环境变量可以非常方便地在不同的部署间做修改，却不动一行代码；
与配置文件不同，不小心把它们签入代码库的概率微乎其微；与一些传统的解决配置问题的机制（比如 Java 的属性配置文件）相比，
环境变量与语言和系统无关。

配置：[config.go](config/config.go)

例如：[.env.example](.env.example)

默认本地 transport 开关：

- `HTTP_ENABLED=true`
- `GRPC_ENABLED=false`
- `RMQ_ENABLED=false`
- `NATS_ENABLED=false`

`APP_ENV=production` 会开启额外防呆：Swagger 必须关闭，示例 JWT 密钥必须替换。

请求关联 ID 是基础脚手架能力：

- HTTP 读取并写回 `X-Request-ID`。
- gRPC、AMQP RPC、NATS RPC 使用 `x-request-id` metadata/header 传递。
- REST 错误响应包含 `request_id`，方便把客户端错误和服务端日志串起来。

可选 trace 默认关闭：

- `TRACE_ENABLED=false`
- `TRACE_EXPORTER=stdout`
- `TRACE_SERVICE_NAME=gin-clean-template`

stdout exporter 是一个真实可运行的绑定，方便本地验证链路。派生项目可以替换成 OTLP/collector 接入，不需要改 handler 或 usecase。

[docker-compose.yml](docker-compose.yml) 使用 `env` 变量配置服务。

### `docs`

Swagger 文档。由  [swag](https://github.com/swaggo/swag) 库自动生成
你不需要自己修改任何内容

[Add a New Business Module](docs/add-new-module.md) 说明如何在不破坏脚手架边界的前提下新增产品业务代码。

#### `docs/proto`

Protobuf 文件。它们用于为 gRPC 服务生成 Go 代码。
这些 proto 文件也用于生成 gRPC 服务的文档。
您不需要自己修改任何内容。

### `integration-test`

集成测试
它会在应用容器旁启动独立的容器
普通 `go test ./...` 不会执行这批测试；直接运行时需要 `integration` build tag。

## `internal/app`

运行时组装放在这里。`Run` 会创建基础设施、repository、usecase、可选 outbox relay，以及当前启用的 transport server。

当前是一个应用进程，可以按配置启动多个 adapter：

- `HTTP_ENABLED`
- `GRPC_ENABLED`
- `RMQ_ENABLED`
- `NATS_ENABLED`

应用通过 root context 统一协调关闭。如果 wiring 规模继续变大，可以在应用边界引入 [wire](https://github.com/google/wire) 这类 DI 生成器；模板默认保留显式构造函数，便于理解依赖关系。

`migrate.go` 文件用于数据库迁移。只有使用 _migrate_ build tag 时才会编译进应用：

```sh
go run -tags migrate ./cmd/app
```

### `internal/transport`

入口适配器层。模板包含 4 种可选 transport：

- AMQP RPC（基于 RabbitMQ 作为传输）
- NATS RPC（基于 NATS 作为传输）
- gRPC（基于 protobuf 的 [gRPC](https://grpc.io/) 框架）
- REST API（基于 [Gin](https://github.com/gin-gonic/gin) 框架）

服务器路由器以相同的风格编写：

- handler 按应用领域分组
- 版本路由依赖通过 `RouterDeps` struct 组织，避免函数签名不断变长
- 路由分组在版本包里显式注册
- 业务逻辑接口注入到 router controller struct 中，handler 只负责调用

#### `internal/transport/amqp_rpc`

基于 RabbitMQ 的 AMQP request/reply adapter。路由在 `amqp_rpc/v1/routes.go` 中注册；认证绑定、请求校验、错误映射和 request ID 传播都留在这个 adapter 内部。

#### `internal/transport/grpc`

基于 `docs/proto/v1` 生成代码的 gRPC adapter。稳定应用错误会映射成 gRPC status，并通过 `google.rpc.ErrorInfo.reason` 携带客户端可识别的错误码。

#### `internal/transport/nats_rpc`

NATS request/reply adapter。它和 AMQP RPC 保持相同的 controller、routes、request、response 布局，方便比较、替换或删除可选 transport。

#### `internal/transport/restapi`

Gin REST adapter。`/healthz`、`/readyz`、`/metrics`、`/swagger/*any` 这类运行时端点在 `internal/transport/restapi/router.go` 注册；版本化业务 API 放在 `internal/transport/restapi/v1` 下。

REST request DTO 放在 `v1/request`，response DTO 放在 `v1/response`，路由分组在 `v1/routes.go` 显式注册。Swagger 注释放在 handler 附近，并通过 [swag](https://github.com/swaggo/swag) 生成。

### `internal/domain`

核心领域模型以及直接属于模型本身的规则。
这一层适合放实体、枚举、值对象和领域错误，并保持其独立于 transport 和存储实现。

### `internal/usecase`

应用层业务逻辑。

- usecase 实现按领域放在子包中，例如 `internal/usecase/task`
- 应用层契约集中在 `internal/usecase/contracts.go`
- usecase 边界以内主要流转 `internal/domain` 值

`usecase` 通过 `internal/usecase/contracts.go` 中定义的抽象依赖外部能力。
持久化实现、传输适配器和可复用技术包通过依赖注入接入 usecase
（阅读 [依赖注入](#依赖注入)）。

#### `internal/infra/persistence`

PostgreSQL 等持久化仓储的具体实现，供 usecase 通过契约调用。

### `pkg`

可复用的技术组件放在这里。这里不应该放应用业务规则。

当前示例包括：

- `pkg/httpserver` 和 `pkg/grpcserver` server wrapper
- `pkg/postgres` 连接池和事务 executor helper
- `pkg/logger` zerolog adapter
- `pkg/requestid` request/correlation ID helper
- `pkg/rabbitmq` 和 `pkg/nats` request/reply 基础组件
- `pkg/observability` 可选 OpenTelemetry 初始化

## 依赖注入

模板使用显式构造函数注入。目标不是隐藏依赖，而是让边界一眼可见。

正常依赖方向是：

- `internal/app` 创建具体基础设施和 usecase 实现
- `internal/usecase` 定义自己需要的契约
- `internal/infra/...` 实现 repository、outbox storage 等出站契约
- `internal/transport/...` 消费入站 usecase 契约

例如 task usecase 接收 repository port 和可选事务 port：

```go
task.New(taskRepo, notificationRepo, transactor)
```

测试使用基于 `internal/usecase/contracts.go` 生成的 mock：

```sh
make mock
```

如果 wiring 变得很大，可以在应用边界引入 [wire](https://github.com/google/wire) 或其他 DI 生成器，但不要改变 contracts 和依赖方向。

## 整洁架构

### 当前规则

当前仓库采用的是务实的 ports-and-adapters 布局：

- `internal/domain` 放 framework-free 的领域模型、领域错误、枚举和模型级规则
- `internal/usecase` 放应用工作流，以及这些工作流需要的契约
- `internal/infra` 放具体出站 adapter，目前包括 PostgreSQL persistence 和 outbox storage
- `internal/transport` 放入站 adapter，目前包括 REST、gRPC、AMQP RPC 和 NATS RPC
- `pkg` 放可复用技术组件，不放产品业务规则

核心依赖规则很简单：内层不导入外层。

允许的例子：

- transport -> usecase contract -> usecase implementation
- usecase implementation -> usecase contract -> infra implementation
- infra implementation -> domain conversion

![Clean Architecture](docs/img/layers-1.png)

不要把 framework 类型、request DTO、persistence row、broker client 或数据库句柄放进 `internal/domain` 或 usecase contracts。

### 边界示例

一个需要数据库数据的 HTTP 请求，流向应该是：

```
    HTTP > usecase
           usecase > persistence contract
           usecase < persistence contract
    HTTP < usecase
```

符号 > 和 < 表示通过接口跨越边界。

![Example](docs/img/example-http-db.png)

如果工作流还需要发布事件，也应保持边界显式：

```
    HTTP > usecase
           usecase > persistence contract
           usecase < persistence contract
           usecase > external integration contract
           usecase < external integration contract
           usecase > RPC
           usecase < RPC
           usecase > persistence contract
           usecase < persistence contract
    HTTP < usecase
```

![Example](docs/img/layers-2.png)

### Domain、DTO 和 Persistence Row

domain value 是业务代码的语言。transport 和 persistence adapter 在边界处转换：

- REST request/response struct 放在 `internal/transport/restapi/v1/request` 和 `response`
- AMQP/NATS request/response struct 放在各自 transport 目录
- gRPC protobuf message 放在 `docs/proto/v1`
- sqlc row 放在 `internal/infra/persistence/sqlc`
- row 和 domain value 之间的转换由 repository 实现负责

### 事务和持久化

对于跨 repository 写入，持久层提供了轻量事务模板，而不是要求使用 ORM：

- `persistence.NewRepositories(pg)` 创建基于普通连接池的 repository。
- `persistence.NewTransactor(pg).WithinTx(ctx, fn)` 创建基于同一个 `pgx` 事务的 repository。
- repository 依赖最小的 `postgres.Executor` 接口，所以同一套 repository 可以跑在连接池或事务上。
- demo repository 的 SQL 放在 `internal/infra/persistence/sql/*.sql`，并通过 `sqlc` 生成到 `internal/infra/persistence/sqlc`。

这是脚手架扩展点。简单的单 repository demo usecase 可以直接调用 repository；需要多表原子更新的流程应使用 `WithinTx`，同时不把 `pgx.Tx` 泄漏到 `internal/usecase`。task 示例在写 task 和 notification 时已经演示这个事务边界。`sqlc` 用于减少手写 scan 样板，repository 仍然负责 domain 转换和稳定错误映射。

### 错误

REST 错误使用稳定 envelope：

```json
{
  "error": {
    "code": "TASK_NOT_FOUND",
    "message": "task not found",
    "request_id": "..."
  }
}
```

映射集中在 `internal/apperror`，思路接近 Kratos 的 `code`/`reason` 拆分：协议状态码仍然表达 HTTP/gRPC 语义，字符串 `code` 作为稳定的客户端错误原因。REST 响应通过 `error.code` 暴露；gRPC 响应通过 `google.rpc.ErrorInfo.reason` 携带；AMQP RPC 和 NATS RPC 使用它作为 RPC status code。示例中的 domain error 放在对应模型文件旁边；跨 REST/gRPC/AMQP/NATS 的错误映射和预期错误日志分级都由 `apperror` 统一处理，避免 transport 层复制多套错误工具。

### Outbox

事件发布提供了接近生产形态的 transactional outbox 示例：

- 通过 migration 创建 `outbox_events` 表
- `outbox.Store` 支持事务内写入、原子 claim pending 事件和 stale publishing lock 回收
- `outbox.Relay` 支持重试、lock timeout、单次 publish timeout 和失败记录
- `outbox.NATSPublisher` 是默认的具体 publisher 绑定，并做 client-side flush

默认通过 `OUTBOX_ENABLED=false` 关闭。启用 `OUTBOX_PUBLISHER=nats` 后，relay 会发布到 `OUTBOX_SUBJECT_PREFIX + "." + event_type`。当业务需要 DB + outbox 一致性时，应在同一个 `WithinTx` 回调里通过事务 `RepoProvider` 暴露的 `OutboxStore` port 写业务数据和 outbox 事件。Core NATS publish + flush 只能确认客户端把消息交给服务端连接，不是 durable broker ack；如果业务事件必须抵抗 broker 侧故障，应替换为 JetStream、Kafka、RabbitMQ confirms 或其他 durable publisher。

## 相似的工程

- [https://github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [https://github.com/zhashkevych/courses-backend](https://github.com/zhashkevych/courses-backend)

## 可能有用的链接

- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/)
