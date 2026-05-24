# 快速启动指南

## 前置要求

- Go 1.22+
- Node.js 20+
- npm 或 yarn

## 后端服务启动

### 1. 配置文件

配置文件已经包含豆包 API Key，位置：`config.yaml`

```yaml
llm:
  api_key: "0dd386ea-afd4-45f1-a7da-efb3dc2e6df7"
```

### 2. 启动后端服务

```bash
# 方式一：使用启动脚本（推荐）
chmod +x start.sh
./start.sh

# 方式二：直接运行
go run cmd/server/main.go

# 方式三：构建后运行
go build -o bin/server cmd/server/main.go
./bin/server
```

服务将在 `http://localhost:8080` 启动

### 3. 验证后端服务

```bash
# 健康检查
curl http://localhost:8080/health

# 创建会话
curl -X POST http://localhost:8080/api/session

# 发送消息（SSE 流）
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "<session-id>",
    "message": "五一假期想去海边玩，预算5000左右，从北京出发"
  }'
```

## 前端应用启动

### 1. 安装依赖

```bash
cd web
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

前端将在 `http://localhost:5173` 启动

### 3. 构建生产版本

```bash
npm run build
```

## 完整流程

### 终端 1：启动后端服务

```bash
./start.sh
```

输出示例：
```
=== 旅行规划 Agent 启动脚本 ===

✓ Go 版本: go1.22.0
✓ 配置文件已找到

正在构建服务...
✓ 构建成功

正在启动服务...
服务地址: http://localhost:8080
健康检查: http://localhost:8080/health
前端地址: http://localhost:5173

按 Ctrl+C 停止服务
```

### 终端 2：启动前端应用

```bash
cd web
npm run dev
```

输出示例：
```
  VITE v5.3.0  ready in 123 ms

  ➜  Local:   http://localhost:5173/
  ➜  press h to show help
```

### 3. 打开浏览器

访问 `http://localhost:5173`，开始使用旅行规划 Agent

## 配置说明

### config.yaml 主要配置项

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `server.port` | 后端服务端口 | 8080 |
| `server.cors_origins` | CORS 允许的源 | http://localhost:5173 |
| `llm.api_key` | 豆包 API Key | 已配置 |
| `llm.model` | 使用的模型 | doubao-seed-2-0-lite-260428 |
| `tools.*.provider_type` | 数据源类型 | mock（模拟数据） |
| `log.level` | 日志级别 | info |

### 修改配置

编辑 `config.yaml` 文件：

```yaml
# 修改后端端口
server:
  port: 9000

# 修改日志级别为 debug
log:
  level: "debug"

# 修改前端 CORS 源
server:
  cors_origins:
    - "http://localhost:3000"
```

## 常见问题

### Q: 后端启动失败，提示 "config.yaml not found"

A: 确保在项目根目录运行启动脚本。如果没有 config.yaml，复制 config.example.yaml：

```bash
cp config.example.yaml config.yaml
```

### Q: 前端无法连接到后端

A: 检查：
1. 后端服务是否在 http://localhost:8080 运行
2. 前端 CORS 配置是否包含 http://localhost:5173
3. 浏览器控制台是否有错误信息

### Q: 豆包 API 返回错误

A: 检查：
1. API Key 是否正确配置在 config.yaml
2. 网络连接是否正常
3. API 配额是否充足

### Q: 如何查看详细日志

A: 修改 config.yaml 中的日志级别：

```yaml
log:
  level: "debug"  # 改为 debug
```

## 项目结构

```
travel-agent/
├── cmd/server/main.go           # 后端入口
├── internal/                    # 后端核心模块
│   ├── server/                  # HTTP 服务器
│   ├── engine/                  # Agent 引擎
│   ├── tool/                    # 工具层
│   ├── provider/                # 数据源层
│   ├── llm/                     # LLM 客户端
│   └── session/                 # 会话存储
├── web/                         # 前端应用
│   ├── src/
│   │   ├── components/          # React 组件
│   │   ├── hooks/               # 自定义 hooks
│   │   ├── stores/              # Zustand 状态管理
│   │   └── api/                 # API 客户端
│   └── package.json
├── config.yaml                  # 配置文件（已配置 API Key）
├── config.example.yaml          # 配置文件示例
├── start.sh                     # 启动脚本
└── README.md                    # 项目文档
```

## 下一步

- 查看 [README.md](./README.md) 了解项目详情
- 查看 [设计文档](./sdd/changes/travel-agent-20260517/design/design.md) 了解架构
- 查看 [测试方案](./sdd/changes/travel-agent-20260517/test/design.md) 了解测试策略
