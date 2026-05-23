# 前端开发任务清单

> 变更：travel-agent-20260517  
> 架构设计：[design.md](../design/design.md) | 详细规格：[specs.md](../design/specs.md)

## 阶段一：项目初始化

### FE-01: React + TypeScript 项目初始化
- **描述**：使用 Vite 初始化 React + TypeScript 项目，配置基础开发环境
- **涉及文件**：
  - `web/package.json`
  - `web/tsconfig.json`
  - `web/vite.config.ts`（配置代理到 Go 后端 8080 端口）
  - `web/src/main.tsx`
  - `web/src/App.tsx`
  - `web/index.html`
- **验收标准**：
  - `npm run dev` 正常启动，浏览器显示空白页面
  - Vite 代理配置：`/api/*` 转发到 `http://localhost:8080`
  - TypeScript 严格模式启用
- **预计耗时**：0.5h
- **依赖**：无

### FE-02: 全局样式与设计规范
- **描述**：配置全局 CSS 变量和基础样式，匹配原型设计规范
- **涉及文件**：
  - `web/src/styles/global.css`（CSS 变量、reset、基础排版）
  - `web/src/styles/variables.css`
- **验收标准**：
  - 主色调：#0891b2 / #06b6d4（蓝青色系）
  - 卡片圆角 12px、轻阴影、白底
  - 字体栈：-apple-system, sans-serif
  - 8px 基础网格间距
  - 响应式：max-width 420px 居中，桌面端自适应
- **预计耗时**：0.5h
- **依赖**：FE-01

### FE-03: TypeScript 类型定义
- **描述**：定义前端使用的所有 TypeScript 类型（与后端对齐）
- **涉及文件**：
  - `web/src/types/index.ts`
- **验收标准**：
  - 所有类型定义与 specs.md 第 5 节完全一致
  - 包括：Phase, Session, Recommendation, Itinerary, BookingItem, ChatMessage, ContentBlock, SSE 事件类型等
  - 导出供其他模块使用
- **预计耗时**：0.5h
- **依赖**：FE-01

## 阶段二：API 层与核心 Hooks

### FE-04: API Client
- **描述**：封装后端 API 调用函数
- **涉及文件**：
  - `web/src/api/client.ts`
- **验收标准**：
  - `createSession()` → POST /api/session
  - `getSession(id)` → GET /api/session/:id
  - `deleteSession(id)` → DELETE /api/session/:id
  - 统一错误处理：解析 `{code, message, data}` 格式
  - 非 0 code 抛出自定义错误
- **预计耗时**：0.5h
- **依赖**：FE-03

### FE-05: useSSE Hook
- **描述**：实现 SSE 连接管理 Hook，处理 POST /api/chat 的 SSE 响应
- **涉及文件**：
  - `web/src/hooks/useSSE.ts`
  - `web/src/utils/sse.ts`（SSE 解析工具函数）
- **验收标准**：
  - 使用 `fetch` + `ReadableStream` 实现 POST 请求的 SSE 读取
  - 正确解析 `event:` 和 `data:` 字段
  - 回调函数：`onThinking`, `onText`, `onToolCall`, `onCard`, `onPhaseChange`, `onDone`, `onError`
  - 支持取消请求（AbortController）
  - 连接断开时自动清理
  - 返回 `{send, cancel, isStreaming}` 
- **预计耗时**：2h
- **依赖**：FE-03

### FE-06: useSession Hook
- **描述**：封装会话生命周期管理
- **涉及文件**：
  - `web/src/hooks/useSession.ts`
- **验收标准**：
  - 提供 `createSession()`, `getSession()`, `deleteSession()`
  - 内部调用 API Client
  - 创建成功后自动保存 sessionId 到 Store
- **预计耗时**：0.5h
- **依赖**：FE-04, FE-07

### FE-07: Zustand Store
- **描述**：实现全局状态管理 Store
- **涉及文件**：
  - `web/src/stores/chatStore.ts`
- **验收标准**：
  - 状态字段：sessionId, phase, messages, isConnecting, isStreaming
  - Actions：createSession, sendMessage, appendMessage, updateStreamingText, setPhase, addContentBlock
  - sendMessage 触发 SSE 连接，将事件转换为状态更新
  - 流式文本增量追加到当前 agent 消息的最后一个 text block
  - card 事件创建新的 card block
  - phase_change 更新 phase
- **预计耗时**：2h
- **依赖**：FE-03, FE-05

## 阶段三：基础 UI 组件

### FE-08: StatusBar 组件
- **描述**：顶部状态栏，显示当前阶段
- **涉及文件**：
  - `web/src/components/StatusBar.tsx`
- **验收标准**：
  - 显示阶段对应的中文描述（COLLECTING→"正在收集您的出行需求"）
  - 脉冲动画状态点
  - 阶段切换时平滑过渡
- **预计耗时**：0.5h
- **依赖**：FE-02, FE-03

### FE-09: InputBar 组件
- **描述**：底部输入框 + 发送按钮
- **涉及文件**：
  - `web/src/components/InputBar.tsx`
- **验收标准**：
  - 文本输入框 + 发送按钮
  - 回车键发送
  - 流式响应中禁用输入（isStreaming 时 disabled）
  - 空内容不可发送
- **预计耗时**：0.5h
- **依赖**：FE-02

### FE-10: MessageBubble 组件
- **描述**：单条消息气泡（用户/Agent）
- **涉及文件**：
  - `web/src/components/MessageBubble.tsx`
- **验收标准**：
  - 用户消息：右对齐，蓝色背景
  - Agent 消息：左对齐，灰色背景
  - 支持流式文本逐字渲染（isStreaming 时内容逐步增长）
  - 渐入动画
- **预计耗时**：0.5h
- **依赖**：FE-02, FE-03

### FE-11: ThinkingIndicator 组件
- **描述**：Agent 思考中的指示器
- **涉及文件**：
  - `web/src/components/ThinkingIndicator.tsx`
- **验收标准**：
  - 可折叠的思考内容块
  - 默认折叠，点击展开查看 thinking 文本
  - 加载动画（三个点）
- **预计耗时**：0.5h
- **依赖**：FE-02

### FE-12: ToolCallIndicator 组件
- **描述**：工具调用状态指示
- **涉及文件**：
  - `web/src/components/ToolCallIndicator.tsx`
- **验收标准**：
  - 显示工具名称（中文映射：flight_search→"搜索航班"）
  - running 状态：旋转图标 + "正在搜索..."
  - done 状态：完成图标 + "搜索完成"
  - error 状态：错误图标 + "搜索失败"
- **预计耗时**：0.5h
- **依赖**：FE-02, FE-03

## 阶段四：业务卡片组件

### FE-13: DestinationCard 组件
- **描述**：目的地推荐卡片
- **涉及文件**：
  - `web/src/components/DestinationCard.tsx`
- **验收标准**：
  - 展示：目的地名称、推荐理由、天气、交通、酒店均价、总费用预估、票务示例
  - "选择这个"按钮，点击触发 onSelect 回调（发送选择消息）
  - is_mock 数据显示"预估"角标
  - 卡片样式与原型一致（圆角 16px、轻阴影）
  - 2x2 网格布局展示信息项
- **预计耗时**：1.5h
- **依赖**：FE-02, FE-03

### FE-14: ItineraryTimeline 组件
- **描述**：行程时间线展示
- **涉及文件**：
  - `web/src/components/ItineraryTimeline.tsx`
  - `web/src/components/TransportCard.tsx`（交通卡片子组件）
  - `web/src/components/HotelCard.tsx`（酒店卡片子组件）
  - `web/src/components/DayPlanCard.tsx`（每日行程子组件）
  - `web/src/components/CostSummary.tsx`（费用汇总子组件）
- **验收标准**：
  - 纵向时间线布局（左侧竖线 + 节点）
  - 交通节点（黄色）、酒店节点（紫色）、日程节点（蓝色）
  - 每日行程：日标题 + 活动列表（时间、名称、费用）
  - 底部费用汇总（分项 + 合计）
  - 底部操作按钮：确认行程 / 我要修改
  - is_mock 数据标注
- **预计耗时**：3h
- **依赖**：FE-02, FE-03

### FE-15: BookingProgress 组件
- **描述**：预订引导进度展示
- **涉及文件**：
  - `web/src/components/BookingProgress.tsx`
  - `web/src/components/BookingItem.tsx`（单个预订项子组件）
- **验收标准**：
  - 顶部进度条（百分比）
  - 预订项列表，按 order 排序
  - 当前项高亮（蓝色边框 + 浅蓝背景）
  - 已完成项：绿色勾选 + "已预订"标签
  - 已跳过项：灰色 + "已跳过"标签
  - 当前项操作按钮：去预订（打开链接）、已支付、跳过
  - "已支付"/"跳过"点击后发送确认消息给 Agent
- **预计耗时**：2h
- **依赖**：FE-02, FE-03

### FE-16: CompleteSummary 组件
- **描述**：预订完成汇总页
- **涉及文件**：
  - `web/src/components/CompleteSummary.tsx`
- **验收标准**：
  - 居中布局
  - 目的地标题 + 完成图标
  - 预订汇总列表（每项显示状态 + 费用）
  - 已预订总花费（大字体突出显示）
  - "完成"按钮
- **预计耗时**：1h
- **依赖**：FE-02, FE-03

## 阶段五：消息列表与主流程

### FE-17: ChatView 消息列表
- **描述**：核心消息列表容器，根据 ContentBlock 类型渲染不同组件
- **涉及文件**：
  - `web/src/components/ChatView.tsx`
- **验收标准**：
  - 遍历 messages，每条消息遍历 blocks
  - block.type='text' → MessageBubble
  - block.type='thinking' → ThinkingIndicator
  - block.type='card' → 根据 cardType 渲染 DestinationCard / ItineraryTimeline / BookingProgress / CompleteSummary
  - block.type='tool_call' → ToolCallIndicator
  - 新消息到达时自动滚动到底部
  - 流式消息实时更新
- **预计耗时**：2h
- **依赖**：FE-10, FE-11, FE-12, FE-13, FE-14, FE-15, FE-16

### FE-18: App 主组件组装
- **描述**：组装所有组件，实现完整的应用主界面
- **涉及文件**：
  - `web/src/App.tsx`
- **验收标准**：
  - 启动时创建 Session
  - 布局：StatusBar（顶） + ChatView（中，flex:1 滚动） + InputBar（底，固定）
  - InputBar 发送消息 → Store.sendMessage → SSE 连接 → 事件更新 UI
  - 卡片上的按钮操作（选择目的地、确认行程、已支付/跳过）发送对应消息给 Agent
  - 全流程可走通
- **预计耗时**：1.5h
- **依赖**：FE-07, FE-08, FE-09, FE-17

## 阶段六：完善与优化

### FE-19: 流式文本渲染优化
- **描述**：优化流式文本的渲染性能和体验
- **涉及文件**：
  - `web/src/components/MessageBubble.tsx`
  - `web/src/stores/chatStore.ts`
- **验收标准**：
  - 使用 requestAnimationFrame 批量更新 DOM
  - 光标闪烁效果（流式文本末尾）
  - 文本内容不闪烁
- **预计耗时**：1h
- **依赖**：FE-18

### FE-20: 错误状态处理
- **描述**：处理各类错误场景的 UI 展示
- **涉及文件**：
  - `web/src/components/ErrorMessage.tsx`
  - `web/src/stores/chatStore.ts`
- **验收标准**：
  - SSE error 事件 → 在对话中显示错误提示气泡
  - 网络断开 → 显示重连提示
  - Session 过期 → 提示重新开始
  - 错误信息用户友好，不暴露技术细节
- **预计耗时**：1h
- **依赖**：FE-18

### FE-21: 响应式适配与移动端优化
- **描述**：确保在不同屏幕尺寸下的良好展示
- **涉及文件**：
  - `web/src/styles/global.css`
  - 各组件样式调整
- **验收标准**：
  - 375px（iPhone SE）正常展示
  - 桌面端（1440px）居中显示，max-width 420px
  - 卡片在窄屏下不溢出
  - 输入框在移动端键盘弹出时正确定位
- **预计耗时**：1h
- **依赖**：FE-18

### FE-22: 端到端手动测试
- **描述**：与后端联调，完成完整流程测试
- **涉及文件**：无新文件
- **验收标准**：
  - 完整流程：创建会话 → 输入需求 → 参数收集 → 目的地推荐 → 选择 → 行程生成 → 确认 → 预订引导 → 完成
  - SSE 事件正确触发 UI 更新
  - 卡片交互（选择、确认、已支付/跳过）正常工作
  - 错误场景（断网、LLM 超时）优雅处理
- **预计耗时**：2h
- **依赖**：FE-18, 后端 BE-26 完成

---

## 任务依赖关系图

```
FE-01 ──┬── FE-02
        ├── FE-03 ──┬── FE-04 ──── FE-06
        │           ├── FE-05 ──── FE-07
        │           ├── FE-07
        │           ├── FE-08
        │           ├── FE-10
        │           ├── FE-12
        │           ├── FE-13
        │           ├── FE-14
        │           ├── FE-15
        │           └── FE-16
        │
        └── FE-09
            FE-11

FE-10 + FE-11 + FE-12 + FE-13 + FE-14 + FE-15 + FE-16 ──── FE-17
FE-07 + FE-08 + FE-09 + FE-17 ──── FE-18
FE-18 ──┬── FE-19
        ├── FE-20
        ├── FE-21
        └── FE-22
```

## 预估总耗时

| 阶段 | 任务数 | 预估耗时 |
|------|--------|---------|
| 阶段一：项目初始化 | 3 | 1.5h |
| 阶段二：API 层与 Hooks | 4 | 5h |
| 阶段三：基础 UI 组件 | 5 | 2.5h |
| 阶段四：业务卡片组件 | 4 | 7.5h |
| 阶段五：消息列表与主流程 | 2 | 3.5h |
| 阶段六：完善与优化 | 4 | 5h |
| **合计** | **22** | **25h** |
