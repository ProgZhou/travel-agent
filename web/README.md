# 旅行规划 Agent - 前端应用

基于 React 18 + TypeScript + Zustand 的对话式旅行规划前端应用。

## 技术栈

- **框架**: React 18.3+
- **语言**: TypeScript 5.4+
- **构建工具**: Vite 5.3+
- **状态管理**: Zustand 4.5+
- **样式**: 原生 CSS（移动端优先设计）

## 功能特性

- ✅ 对话式 UI，支持流式文本渲染
- ✅ SSE（Server-Sent Events）实时通信
- ✅ 目的地推荐卡片展示
- ✅ 行程时间线可视化
- ✅ 预订引导流程
- ✅ 响应式布局（移动端优先，375px 基准）
- ✅ 完整的 TypeScript 类型定义

## 快速开始

### 1. 安装依赖

```bash
cd web
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:5173 启动。

**注意**: 前端需要后端 API 服务运行在 `http://localhost:8080`。Vite 已配置代理，所有 `/api` 请求会自动转发到后端。

### 3. 构建生产版本

```bash
npm run build
```

构建产物将输出到 `dist/` 目录。

### 4. 预览生产构建

```bash
npm run preview
```

## 项目结构

```
web/
├── src/
│   ├── components/          # UI 组件
│   │   ├── StatusBar.tsx    # 状态栏
│   │   ├── ChatView.tsx     # 对话视图
│   │   ├── MessageBubble.tsx # 消息气泡
│   │   ├── InputBar.tsx     # 输入框
│   │   ├── DestinationCard.tsx # 目的地推荐卡片
│   │   ├── ItineraryTimeline.tsx # 行程时间线
│   │   ├── BookingProgress.tsx # 预订进度
│   │   └── CompleteSummary.tsx # 完成汇总
│   ├── hooks/               # 自定义 Hooks
│   │   ├── useSSE.ts        # SSE 连接管理
│   │   └── useSession.ts    # 会话生命周期
│   ├── stores/              # Zustand 状态管理
│   │   └── chatStore.ts     # 聊天状态
│   ├── types/               # TypeScript 类型定义
│   │   └── index.ts
│   ├── api/                 # API 客户端
│   │   └── client.ts
│   ├── App.tsx              # 根组件
│   ├── App.css              # 全局样式
│   └── main.tsx             # 入口文件
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## 核心设计

### 状态管理

使用 Zustand 管理全局状态：

- `sessionId`: 当前会话 ID
- `phase`: 当前阶段（COLLECTING / RECOMMENDING / PLANNING / BOOKING / COMPLETED）
- `messages`: 消息列表
- `isStreaming`: 是否正在接收流式响应

### SSE 通信

`useSSE` hook 实现了基于 Fetch API 的 SSE 客户端，支持：

- POST 请求发送 SSE
- 事件类型：thinking / text / tool_call / card / phase_change / done / error
- 自动重连和错误处理

### 消息模型

每条消息包含多个内容块（ContentBlock）：

- `text`: 文本内容
- `thinking`: 思考过程
- `tool_call`: 工具调用状态
- `card`: 结构化卡片（推荐/行程/预订）

### 响应式设计

- 移动端优先（375px 基准宽度）
- 桌面端最大宽度 420px 居中显示
- 所有交互元素适配触摸操作

## API 接口

前端调用以下后端接口：

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 创建会话 | POST | /api/session | 创建新会话 |
| 获取会话 | GET | /api/session/:id | 获取会话状态 |
| 删除会话 | DELETE | /api/session/:id | 删除会话 |
| 发送消息 | POST | /api/chat | 发送消息（SSE 流） |

详细接口协议见 `sdd/changes/travel-agent-20260517/design/specs.md`。

## 开发说明

### 添加新组件

1. 在 `src/components/` 创建 `.tsx` 和 `.css` 文件
2. 使用 TypeScript 严格模式，避免 `any` 类型
3. 为交互元素添加 `aria-label` 等无障碍属性
4. 样式遵循 8px 基础网格

### 修改状态管理

编辑 `src/stores/chatStore.ts`，所有状态变更通过 actions 进行。

### 调试 SSE

打开浏览器开发者工具 Network 面板，筛选 EventStream 类型查看 SSE 事件流。

## 常见问题

**Q: 前端启动后无法连接后端？**

A: 确保后端服务运行在 `http://localhost:8080`，检查 `vite.config.ts` 中的 proxy 配置。

**Q: SSE 连接中断？**

A: 检查后端日志，确认 `/api/chat` 接口返回 `Content-Type: text/event-stream`。

**Q: 样式不生效？**

A: 确认 CSS 文件已在对应 `.tsx` 文件中导入。

## 参考文档

- [产品需求文档](../sdd/changes/travel-agent-20260517/product/product.md)
- [架构设计文档](../sdd/changes/travel-agent-20260517/design/design.md)
- [详细规格文档](../sdd/changes/travel-agent-20260517/design/specs.md)
- [交互原型](../sdd/changes/travel-agent-20260517/product/prototype.html)

## License

MIT
