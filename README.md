# 旅行规划 Agent

基于 Go + React 的智能旅行规划助手，使用豆包大模型提供对话式旅行规划服务。

## 功能特性

- 🤖 对话式需求收集
- 🌍 智能目的地推荐
- 📅 详细行程规划
- 🔗 预订引导
- 📊 实时流式响应（SSE）
- 🛠️ 工具调用（航班、酒店、天气查询）

## 技术栈

### 后端
- Go 1.22+
- chi 路由框架
- 豆包大模型（火山引擎）
- SSE 流式推送

### 前端
- React 18
- TypeScript
- Zustand 状态管理
- Vite 构建工具

## 快速开始

### 前置要求

- Go 1.22 或更高版本
- Node.js 18 或更高版本
- 豆包 API Key（火山引擎）

### 1. 克隆项目

```bash
git clone <repository-url>
cd travel_agent
```

### 2. 配置环境变量

```bash
# 设置豆包 API Key
export ARK_API_KEY="your-api-key-here"
```

### 3. 启动后端服务

```bash
# 安装依赖
go mod tidy

# 运行服务
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 4. 启动前端（可选）

```bash
cd web
npm install
npm run dev
```

前端将在 `http://localhost:5173` 启动。

## 配置说明

配置文件位于 `config.yaml`，主要配置项：

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  cors_origins:
    - "http://localhost:5173"

llm:
  provider: "doubao"
  api_url: "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
  model: "doubao-seed-2-0-lite-260428"
  max_tokens: 4096
  temperature: 0.7
  timeout: 30s

tools:
  flight:
    provider_type: "mock"  # 当前使用 Mock 数据
  hotel:
    provider_type: "mock"
  weather:
    provider_type: "mock"
```

### 环境变量覆盖

- `ARK_API_KEY`: 豆包 API Key（必须）
- `LOG_LEVEL`: 日志级别（debug/info/warn/error）

## API 接口

### 创建会话

```bash
POST /api/session
```

### 获取会话

```bash
GET /api/session/:id
```

### 发送消息（SSE 流）

```bash
POST /api/chat
Content-Type: application/json

{
  "session_id": "uuid",
  "message": "五一假期想去海边玩，预算5000左右，从北京出发"
}
```

## 项目结构

```
travel-agent/
├── cmd/server/          # 程序入口
├── internal/
│   ├── config/          # 配置管理
│   ├── server/          # HTTP 服务器
│   ├── engine/          # Agent 引擎
│   ├── tool/            # 工具层
│   ├── provider/        # 数据源 Provider
│   ├── llm/             # LLM 客户端
│   ├── session/         # 会话存储
│   └── pkg/             # 公共工具
├── web/                 # 前端项目
├── config.yaml          # 配置文件
└── README.md
```

## 开发说明

### 添加新的 Tool

1. 在 `internal/provider/` 中实现 Provider 接口
2. 在 `internal/tool/` 中创建 Tool 实现
3. 在 `cmd/server/main.go` 中注册 Tool

### 日志级别

- `debug`: LLM 请求/响应详情、Tool 调用详情
- `info`: 请求入口、阶段流转、Tool 调用结果
- `warn`: 重试事件、降级触发
- `error`: 不可恢复错误

## 注意事项

- 当前使用 Mock Provider，返回的数据为预估数据
- 会话数据存储在内存中，重启后丢失
- 单人使用，无认证机制
- 不收集或存储支付信息

## License

MIT
