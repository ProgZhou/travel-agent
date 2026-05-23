# 后端开发任务清单

> 变更：travel-agent-20260517  
> 架构设计：[design.md](../design/design.md) | 详细规格：[specs.md](../design/specs.md)

## 阶段一：项目骨架与基础设施

### BE-01: 项目初始化与目录结构
- **描述**：初始化 Go Module，创建完整目录结构，添加 Makefile
- **涉及文件**：
  - `go.mod`
  - `Makefile`（build/run/test 命令）
  - `cmd/server/main.go`（空 main 函数占位）
  - 所有 `internal/` 子目录创建
- **验收标准**：`go build ./...` 成功，`make run` 可启动空服务
- **预计耗时**：0.5h
- **依赖**：无

### BE-02: 配置管理模块
- **描述**：实现 YAML 配置加载 + 环境变量覆盖机制
- **涉及文件**：
  - `internal/config/config.go`（Config 结构体 + Load 函数）
  - `config.yaml`（配置模板）
- **验收标准**：
  - 可加载 config.yaml
  - 环境变量 `ARK_API_KEY` 覆盖 llm.api_key
  - 缺少必要配置时程序报错退出
- **预计耗时**：1h
- **依赖**：BE-01

### BE-03: 日志模块
- **描述**：基于 `log/slog` 封装结构化日志，支持配置日志级别和格式
- **涉及文件**：
  - `internal/pkg/logger.go`
- **验收标准**：
  - 支持 debug/info/warn/error 级别
  - 日志包含 timestamp、level、msg、结构化字段
  - 通过配置文件控制级别
- **预计耗时**：0.5h
- **依赖**：BE-02

### BE-04: 统一错误处理
- **描述**：定义业务错误类型和错误码常量
- **涉及文件**：
  - `internal/pkg/errors.go`
- **验收标准**：
  - 定义 AppError 结构体（Code, Message, HTTPStatus）
  - 定义所有错误码常量（10001-20002）
  - 提供 NewAppError / Wrap 等辅助方法
- **预计耗时**：0.5h
- **依赖**：BE-01

## 阶段二：会话存储层

### BE-05: Session 数据模型
- **描述**：定义 Session 及所有嵌套数据结构的 Go struct
- **涉及文件**：
  - `internal/session/model.go`（Session, UserRequest, Recommendation, Itinerary, BookingItem, ChatMessage 等全部 struct）
- **验收标准**：
  - 所有 struct 定义与 specs.md 一致
  - JSON tag 正确
  - 枚举类型（Phase, BookingType, BookingStatus）有常量定义
- **预计耗时**：1h
- **依赖**：BE-01

### BE-06: SessionStore 接口与内存实现
- **描述**：定义 SessionStore 接口，实现 MemorySessionStore
- **涉及文件**：
  - `internal/session/store.go`（接口定义）
  - `internal/session/memory.go`（内存实现）
- **验收标准**：
  - Get/Save/Delete 正确工作
  - 使用 sync.RWMutex 保证并发安全
  - Save 时自动更新 UpdatedAt
  - Get 返回 Session 的深拷贝（避免外部修改影响存储）
- **预计耗时**：1h
- **依赖**：BE-05
- **测试**：编写 memory_test.go，覆盖 CRUD + 并发场景

## 阶段三：LLM 调用层

### BE-07: LLM Client 接口定义
- **描述**：定义 LLM 调用的接口和类型
- **涉及文件**：
  - `internal/llm/client.go`（LLMClient 接口）
  - `internal/llm/types.go`（ChatRequest, ChatResponse, StreamChunk 等类型）
- **验收标准**：
  - 接口支持 ChatCompletion（非流式）和 ChatCompletionStream（流式）
  - 类型定义与豆包 API 文档对齐
- **预计耗时**：0.5h
- **依赖**：BE-01

### BE-08: 豆包 API 客户端实现
- **描述**：实现豆包大模型 API 的 HTTP 调用，支持流式和非流式
- **涉及文件**：
  - `internal/llm/doubao.go`
- **验收标准**：
  - 非流式调用：发送请求，解析完整 JSON 响应
  - 流式调用：发送请求，解析 SSE 流，通过 channel 返回 StreamChunk
  - 支持 Function Calling（tools 参数）
  - 请求超时 30s
  - 重试 2 次，间隔 3s
  - API Key 通过 Authorization Bearer 头传递
- **预计耗时**：3h
- **依赖**：BE-07, BE-02
- **测试**：编写 doubao_test.go（可用 httptest Mock 豆包 API）

## 阶段四：Provider 层（数据源）

### BE-09: Provider 接口定义
- **描述**：定义 FlightProvider、HotelProvider、WeatherProvider 接口
- **涉及文件**：
  - `internal/provider/provider.go`（三个接口 + Request/Result 类型）
- **验收标准**：
  - 接口定义与 specs.md 一致
  - 每个接口一个方法（Search/Query），接受 context.Context
- **预计耗时**：0.5h
- **依赖**：BE-05

### BE-10: MockFlightProvider 实现
- **描述**：实现航班/火车票 Mock 数据源
- **涉及文件**：
  - `internal/provider/flight_mock.go`
- **验收标准**：
  - 内置北京-三亚、北京-厦门、上海-三亚等常见城市对的数据
  - 返回 2-5 条合理的班次数据
  - 价格在合理范围内（经济舱 300-1500 元，高铁 100-800 元）
  - 结果携带 is_mock=true
  - 未覆盖城市返回通用估算数据 + notice
- **预计耗时**：1.5h
- **依赖**：BE-09

### BE-11: MockHotelProvider 实现
- **描述**：实现酒店查询 Mock 数据源
- **涉及文件**：
  - `internal/provider/hotel_mock.go`
- **验收标准**：
  - 内置三亚、厦门、大理等城市的酒店数据
  - 每个城市提供经济/舒适/豪华档各 2-3 家
  - 价格合理（经济 100-300、舒适 250-600、豪华 500-2000）
  - 结果携带 is_mock=true
- **预计耗时**：1h
- **依赖**：BE-09

### BE-12: MockWeatherProvider 实现
- **描述**：实现天气查询 Mock 数据源
- **涉及文件**：
  - `internal/provider/weather_mock.go`
- **验收标准**：
  - 基于城市和月份返回合理历史气候数据
  - 三亚 5 月：26-32°C、晴为主
  - 厦门 5 月：22-28°C、多云偶阵雨
  - 返回每日天气 + 汇总概况
  - 结果携带 is_mock=true
- **预计耗时**：1h
- **依赖**：BE-09

## 阶段五：Tool 层

### BE-13: Tool 接口与 Registry 实现
- **描述**：定义 Tool 接口，实现 ToolRegistry
- **涉及文件**：
  - `internal/tool/tool.go`（Tool 接口 + ToolSchema 类型）
  - `internal/tool/registry.go`（ToolRegistry 实现）
- **验收标准**：
  - Tool 接口包含 Name/Description/Parameters/Execute 方法
  - ToolRegistry 支持 Register/Get/AllSchemas
  - AllSchemas 返回所有 Tool 的 JSON Schema（供 LLM tools 参数）
- **预计耗时**：1h
- **依赖**：BE-09

### BE-14: FlightTool 实现
- **描述**：实现航班搜索 Tool，内含超时/重试/降级逻辑
- **涉及文件**：
  - `internal/tool/flight.go`
- **验收标准**：
  - 返回 flight_search 的 JSON Schema
  - Execute 时调用 FlightProvider.Search
  - 超时 5s + 重试 1 次 + 降级为 Mock
  - 降级时 ToolResult.IsMock=true, Notice="预估数据，仅供参考"
- **预计耗时**：1.5h
- **依赖**：BE-13, BE-10, BE-02

### BE-15: HotelTool 实现
- **描述**：实现酒店搜索 Tool
- **涉及文件**：
  - `internal/tool/hotel.go`
- **验收标准**：同 BE-14，调用 HotelProvider
- **预计耗时**：1h
- **依赖**：BE-13, BE-11, BE-02

### BE-16: WeatherTool 实现
- **描述**：实现天气查询 Tool
- **涉及文件**：
  - `internal/tool/weather.go`
- **验收标准**：同 BE-14，调用 WeatherProvider
- **预计耗时**：1h
- **依赖**：BE-13, BE-12, BE-02

## 阶段六：Agent Engine

### BE-17: Event 类型定义
- **描述**：定义 Engine 对外输出的所有事件类型
- **涉及文件**：
  - `internal/engine/event.go`
- **验收标准**：
  - Event 结构体 + 所有子事件类型（TextEvent, ThinkingEvent, ToolCallEvent, CardEvent, PhaseChangeEvent, ErrorEvent）
  - 提供 NewXxxEvent 构造函数
- **预计耗时**：0.5h
- **依赖**：BE-05

### BE-18: Prompt 管理器
- **描述**：实现各阶段的 System Prompt 模板管理
- **涉及文件**：
  - `internal/engine/prompt/manager.go`
  - `internal/engine/prompt/system.go`（基础人设）
  - `internal/engine/prompt/collecting.go`
  - `internal/engine/prompt/recommending.go`
  - `internal/engine/prompt/planning.go`
  - `internal/engine/prompt/booking.go`
- **验收标准**：
  - PromptManager.BuildSystemPrompt(session) 返回完整 Prompt
  - Prompt = 基础人设 + 阶段指令 + 上下文（动态注入已收集的参数等）
  - 每个阶段的 Prompt 清晰指导 LLM 行为
- **预计耗时**：2h
- **依赖**：BE-05

### BE-19: Planner 阶段编排器
- **描述**：实现 Planner，负责判断阶段流转
- **涉及文件**：
  - `internal/engine/planner.go`
- **验收标准**：
  - Plan(ctx, session, input) 调用 LLM（非流式），返回 PlanResult
  - PlanResult 包含：是否需要流转 + 目标阶段 + 当前阶段任务描述
  - 遵循状态机规则（COLLECTING 只能到 RECOMMENDING，不能跳到 BOOKING）
  - LLM 返回异常时保持当前阶段不变
- **预计耗时**：2h
- **依赖**：BE-07, BE-08, BE-18

### BE-20: Executor ReAct 执行器
- **描述**：实现 Executor，在单个阶段内运行 ReAct 循环
- **涉及文件**：
  - `internal/engine/executor.go`
- **验收标准**：
  - Execute(ctx, session, task) 返回 Event channel
  - ReAct 循环：调用 LLM → 解析响应 → 若 tool_call 则执行工具 → 将结果加入上下文 → 继续循环
  - 流式文本实时推送 TextEvent
  - 多工具并行调用（goroutine）
  - 最大循环 10 次（可配置），超限退出
  - 正确处理 context 取消
- **预计耗时**：4h
- **依赖**：BE-07, BE-08, BE-13, BE-17, BE-18

### BE-21: Engine 入口实现
- **描述**：实现 Engine.Run，编排 Planner + Executor
- **涉及文件**：
  - `internal/engine/engine.go`
- **验收标准**：
  - Run(ctx, session, input) 返回 Event channel
  - 流程：保存用户消息 → Planner 判断 → 执行阶段流转 → Executor 执行 → 保存 Session
  - 阶段变更时推送 PhaseChangeEvent
  - 错误时推送 ErrorEvent + done
  - context 取消时优雅退出
- **预计耗时**：2h
- **依赖**：BE-19, BE-20, BE-06

## 阶段七：HTTP Server

### BE-22: HTTP Server 与路由
- **描述**：设置 HTTP Server，注册路由，配置中间件
- **涉及文件**：
  - `internal/server/server.go`
  - `internal/server/router.go`
  - `internal/server/middleware.go`（CORS, Recovery, 请求日志）
- **验收标准**：
  - 所有 API 路由正确注册
  - CORS 允许配置的前端域名
  - Recovery 中间件捕获 panic
  - 请求日志记录 method/path/status/duration
- **预计耗时**：1.5h
- **依赖**：BE-02, BE-03

### BE-23: Session API Handler
- **描述**：实现 Session CRUD 的 HTTP Handler
- **涉及文件**：
  - `internal/server/handler/session.go`
- **验收标准**：
  - POST /api/session：创建会话，返回 session_id
  - GET /api/session/:id：返回完整 Session JSON
  - DELETE /api/session/:id：删除会话
  - 参数校验（UUID 格式）
  - 错误码正确返回（404 会话不存在）
- **预计耗时**：1h
- **依赖**：BE-06, BE-22, BE-04

### BE-24: Chat SSE Handler
- **描述**：实现 POST /api/chat 的 SSE 流式响应
- **涉及文件**：
  - `internal/server/handler/chat.go`
- **验收标准**：
  - 接收 JSON 请求体（session_id + message）
  - 校验参数（session_id 必填、message 非空且 <= 2000 字符）
  - 设置 SSE 响应头
  - 调用 Engine.Run，将 Event channel 转换为 SSE 事件流
  - 客户端断开时通过 context 取消处理
  - 请求体校验失败返回标准 JSON 错误（非 SSE）
- **预计耗时**：2h
- **依赖**：BE-21, BE-22, BE-04

### BE-25: Health Check Handler
- **描述**：实现健康检查接口
- **涉及文件**：
  - `internal/server/handler/health.go`
- **验收标准**：
  - GET /api/health 返回 `{"status": "ok"}`
- **预计耗时**：0.5h
- **依赖**：BE-22

## 阶段八：启动与集成

### BE-26: 程序入口与依赖组装
- **描述**：在 main.go 中完成所有模块的初始化和依赖注入
- **涉及文件**：
  - `cmd/server/main.go`
- **验收标准**：
  - 加载配置 → 初始化日志 → 创建 LLM Client → 创建 Providers → 创建 Tools → 注册到 Registry → 创建 SessionStore → 创建 Engine → 创建 Server → 启动
  - 优雅关闭（处理 SIGINT/SIGTERM）
  - 启动日志输出端口和配置信息
- **预计耗时**：1h
- **依赖**：BE-22, BE-24, BE-21

### BE-27: 端到端冒烟测试
- **描述**：编写集成测试，验证从 HTTP 请求到 SSE 响应的完整链路
- **涉及文件**：
  - `cmd/server/main_test.go` 或 `internal/server/integration_test.go`
- **验收标准**：
  - 创建会话 → 发送消息 → 接收 SSE 事件 → 验证事件类型和顺序
  - 使用 Mock LLM Client（可选：真实豆包调用标记为 integration tag）
- **预计耗时**：2h
- **依赖**：BE-26

---

## 任务依赖关系图

```
BE-01 ──┬── BE-02 ──── BE-03
        │     │
        │     ├── BE-08
        │     ├── BE-14
        │     ├── BE-15
        │     └── BE-16
        │
        ├── BE-04
        │
        ├── BE-05 ──┬── BE-06 ──── BE-21
        │           │
        │           ├── BE-09 ──┬── BE-10 ──── BE-14
        │           │           ├── BE-11 ──── BE-15
        │           │           ├── BE-12 ──── BE-16
        │           │           └── BE-13 ──┬── BE-14
        │           │                       ├── BE-15
        │           │                       └── BE-16
        │           │
        │           ├── BE-17
        │           └── BE-18 ──┬── BE-19
        │                       └── BE-20
        │
        └── BE-07 ──── BE-08 ──┬── BE-19
                               └── BE-20

BE-19 + BE-20 ──── BE-21 ──── BE-24 ──── BE-26 ──── BE-27
BE-22 ──── BE-23
BE-22 ──── BE-24
BE-22 ──── BE-25
```

## 预估总耗时

| 阶段 | 任务数 | 预估耗时 |
|------|--------|---------|
| 阶段一：项目骨架 | 4 | 2.5h |
| 阶段二：会话存储 | 2 | 2h |
| 阶段三：LLM 调用 | 2 | 3.5h |
| 阶段四：Provider | 4 | 4h |
| 阶段五：Tool | 4 | 4.5h |
| 阶段六：Engine | 5 | 10.5h |
| 阶段七：HTTP Server | 4 | 5h |
| 阶段八：集成 | 2 | 3h |
| **合计** | **27** | **35h** |
