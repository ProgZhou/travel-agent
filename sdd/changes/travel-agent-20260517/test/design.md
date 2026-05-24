# 测试方案 - 旅行规划 Agent

## 1. 测试范围

### 1.1 功能测试
- **需求收集阶段（COLLECTING）**
  - 自然语言意图解析（时间、预算、偏好、出发地提取）
  - 缺失参数追问逻辑
  - 参数完整性判断与阶段流转
- **目的地推荐阶段（RECOMMENDING）**
  - 推荐卡片生成（2-3个目的地）
  - 天气/交通/酒店数据聚合
  - 用户选择目的地响应
- **行程规划阶段（PLANNING）**
  - 完整行程表生成（去程/住宿/每日安排/返程/费用汇总）
  - 行程修改处理（换景点/换酒店/调时间）
  - 预算超限提示
- **预订引导阶段（BOOKING）**
  - 预订链接生成（携程/飞猪/去哪儿）
  - 预订状态追踪（已订/跳过/待处理）
  - 预订汇总输出
- **会话管理**
  - 会话创建/查询/删除
  - 会话状态持久化
  - 断点恢复

### 1.2 接口测试
- **POST /api/session** - 创建会话
- **GET /api/session/:id** - 获取会话状态
- **DELETE /api/session/:id** - 删除会话
- **POST /api/chat** - 发送消息（SSE 流）
  - SSE 事件类型验证（thinking/text/tool_call/card/phase_change/done/error）
  - 流式文本增量推送
  - 卡片数据完整性
  - 错误码覆盖（10001/10002/10003/20001/20002）

### 1.3 安全测试
- **输入校验**
  - 消息长度限制（1-2000字符）
  - Session ID 格式校验（UUID v4）
  - 请求体大小限制（10KB）
  - 特殊字符注入（SQL注入、XSS）
- **API Key 保护**
  - 环境变量注入验证
  - 配置文件不含真实 Key
- **CORS 配置**
  - 仅允许配置的前端域名
  - 拒绝未授权来源
- **错误信息脱敏**
  - 不暴露内部堆栈
  - 不泄露配置细节

### 1.4 性能测试
- **响应时间**
  - 单轮响应时间 < 15秒（含 API 调用和 LLM 推理）
  - 推荐生成时间 < 20秒（并行调用多个 API）
  - 行程生成时间 < 15秒
  - 本地存储读写 < 100ms
- **并发测试**
  - 10 个并发会话同时处理
  - SSE 连接稳定性（长连接保持）
- **资源占用**
  - 单会话内存占用 < 100KB
  - 100 个会话总内存 < 10MB

### 1.5 兼容性测试
- **浏览器兼容性**
  - Chrome 120+
  - Safari 17+
  - Firefox 120+
  - Edge 120+
- **屏幕尺寸**
  - 移动端：375px（iPhone SE）、390px（iPhone 14）、414px（iPhone 14 Pro Max）
  - 桌面端：1280px、1440px、1920px
- **操作系统**
  - macOS 14+
  - Windows 11
  - Linux（Ubuntu 22.04）

### 1.6 稳定性测试
- **降级策略**
  - LLM 超时降级（重试2次 → 返回错误）
  - Tool Provider 超时降级（重试1次 → Mock数据）
  - Mock 数据标记验证（is_mock=true）
- **异常恢复**
  - 客户端断开连接后 context 取消
  - ReAct 循环超限（10轮）强制终止
  - 会话状态自动保存
- **长时间运行**
  - 24小时稳定性测试
  - 内存泄漏检测

## 2. 测试策略

### 2.1 单元测试
**后端 Go**
- **覆盖率目标**：核心模块 ≥ 80%
- **测试框架**：Go testing + testify/assert + testify/mock
- **重点模块**：
  - `engine/planner.go` - 阶段流转逻辑
  - `engine/executor.go` - ReAct 循环
  - `tool/*.go` - 工具调用与降级
  - `provider/*_mock.go` - Mock 数据生成
  - `session/memory.go` - 会话存储并发安全
- **Mock 策略**：
  - LLM Client 使用 mock 接口，预设响应
  - Provider 接口 mock，控制超时/错误场景
  - Time mock 用于测试超时逻辑

**前端 React**
- **覆盖率目标**：组件 ≥ 70%，工具函数 ≥ 90%
- **测试框架**：Vitest + React Testing Library
- **重点模块**：
  - `hooks/useSSE.ts` - SSE 事件解析
  - `hooks/useSession.ts` - 会话管理
  - `stores/chatStore.ts` - 状态管理逻辑
  - `utils/sse.ts` - SSE 工具函数
- **Mock 策略**：
  - fetch API mock（MSW - Mock Service Worker）
  - SSE 流 mock（自定义 ReadableStream）

### 2.2 集成测试
**后端 API 集成**
- **测试框架**：Go testing + httptest
- **测试范围**：
  - HTTP Handler → Engine → Tool → Provider 完整链路
  - SSE 事件流完整性（从请求到 done 事件）
  - 会话状态持久化与恢复
- **测试环境**：
  - 使用 httptest.Server 启动测试服务器
  - 使用 Mock Provider（不依赖真实 API）
  - 使用 Memory SessionStore（不依赖数据库）
- **关键场景**：
  - 完整对话流程（COLLECTING → RECOMMENDING → PLANNING → BOOKING → COMPLETED）
  - 阶段流转触发条件
  - Tool 调用超时与降级
  - 并发会话隔离

**前后端联调**
- **测试框架**：Playwright（E2E）
- **测试范围**：
  - 前端发送消息 → 后端 SSE 推送 → 前端渲染
  - 卡片交互（选择目的地、确认行程、预订操作）
  - 断线重连
- **测试环境**：
  - 后端启动真实服务（使用 Mock Provider）
  - 前端启动 Vite dev server
  - Playwright 自动化浏览器操作

### 2.3 E2E 测试
**工具选型**：Playwright
- **覆盖场景**：
  - 完整旅行规划流程（从输入需求到预订完成）
  - 行程修改流程
  - 错误处理（API 超时、会话不存在）
  - 多浏览器兼容性
- **测试数据**：
  - 预设 5 个典型用户需求场景
  - 覆盖不同预算区间（3000/5000/8000/10000/15000）
  - 覆盖不同目的地类型（海边/山区/城市/古镇）
- **断言策略**：
  - 页面元素可见性
  - 卡片数据完整性
  - 阶段状态正确性
  - 最终预订汇总准确性

### 2.4 性能测试
**工具选型**：
- 后端：Go benchmark + pprof
- 压测：k6（模拟并发用户）
- 前端：Lighthouse + Chrome DevTools

**测试场景**：
- 单用户完整流程响应时间
- 10/50/100 并发用户压测
- SSE 长连接稳定性（保持 5 分钟）
- 内存占用监控（100 个会话）

## 3. 测试环境

### 3.1 后端测试环境
- **操作系统**：macOS 14+ / Linux（CI 环境）
- **Go 版本**：1.22+
- **依赖服务**：
  - 豆包 LLM：使用 Mock Client（单元测试）/ 真实 API（集成测试，需配置 ARK_API_KEY）
  - Provider：统一使用 Mock 实现
- **配置文件**：`config.test.yaml`（测试专用配置，超时时间缩短）

### 3.2 前端测试环境
- **Node 版本**：20+
- **浏览器**：
  - 单元测试：jsdom（Vitest 内置）
  - E2E 测试：Chromium/Firefox/WebKit（Playwright）
- **Mock 服务**：MSW（Mock Service Worker）拦截 API 请求

### 3.3 CI/CD 环境
- **平台**：GitHub Actions
- **测试流水线**：
  1. 后端单元测试（Go test）
  2. 前端单元测试（Vitest）
  3. 后端集成测试（httptest）
  4. E2E 测试（Playwright，仅在 PR 时运行）
- **覆盖率报告**：上传至 Codecov

## 4. 测试数据

### 4.1 基础数据集
**用户需求场景**
| 场景 ID | 描述 | 出发地 | 时间 | 预算 | 偏好 | 预期目的地 |
|---------|------|--------|------|------|------|-----------|
| REQ-01 | 五一海边游 | 北京 | 2026-05-01 ~ 2026-05-05 | 5000 | 海边、人少 | 三亚/厦门 |
| REQ-02 | 清明古镇游 | 上海 | 2026-04-04 ~ 2026-04-06 | 3000 | 古镇、安静 | 乌镇/周庄 |
| REQ-03 | 国庆山区游 | 广州 | 2026-10-01 ~ 2026-10-07 | 8000 | 登山、摄影 | 张家界/黄山 |
| REQ-04 | 暑假亲子游 | 成都 | 2026-07-15 ~ 2026-07-20 | 10000 | 亲子、乐园 | 上海迪士尼/珠海长隆 |
| REQ-05 | 周末短途游 | 北京 | 2026-06-14 ~ 2026-06-15 | 2000 | 放松、温泉 | 古北水镇/南戴河 |

**Mock Provider 数据覆盖**
- **航班/火车**：北京-三亚、北京-厦门、上海-三亚、上海-杭州、广州-三亚、成都-大理（各3个班次）
- **酒店**：三亚/厦门/大理/杭州（各档次3家酒店）
- **天气**：三亚/厦门/大理/杭州（5-10月每月典型天气）

### 4.2 边界数据集
- **极端预算**：500元（低于合理范围）、50000元（超高预算）
- **极端时间**：当天出发、1年后出发、跨年行程
- **极端输入**：
  - 空消息
  - 2000字符消息（上限）
  - 2001字符消息（超限）
  - 特殊字符：`<script>alert('xss')</script>`、`' OR '1'='1`
- **无效 Session ID**：非 UUID 格式、不存在的 UUID

### 4.3 数据清理策略
- **单元测试**：每个测试用例独立创建 Session，测试结束后销毁
- **集成测试**：使用独立的 Memory SessionStore 实例，测试结束后清空
- **E2E 测试**：每个测试场景开始前调用 DELETE /api/session/:id 清理旧会话

## 5. 风险分析

### 5.1 高风险模块
| 模块 | 风险描述 | 缓解措施 |
|------|---------|---------|
| LLM 调用 | 豆包 API 不稳定、返回格式异常、幻觉 | Mock Client 覆盖异常场景；集成测试使用真实 API 验证；Prompt 明确约束输出格式 |
| SSE 流式推送 | 网络中断、客户端断开、事件丢失 | Context 取消传播；前端重连机制；事件序列化测试 |
| 阶段流转逻辑 | 状态机错误流转、死循环 | 单元测试覆盖所有流转路径；集成测试验证完整流程；最大迭代次数保护 |
| Tool 降级 | 降级数据不合理、标记缺失 | Mock Provider 返回合理数据；单元测试验证 is_mock 标记；集成测试验证降级触发 |
| 并发安全 | SessionStore 并发读写冲突 | 单元测试使用 race detector；压测验证并发场景 |

### 5.2 测试限制
- **LLM 不确定性**：大模型输出不完全可控，测试断言需关注结构而非具体文本
- **真实 API 依赖**：集成测试依赖豆包 API 可用性，CI 环境需配置 API Key
- **Mock 数据局限**：Mock Provider 仅覆盖常见城市对，边缘场景可能缺失
- **E2E 测试成本**：Playwright 测试耗时较长，仅在 PR 时运行，不阻塞本地开发

## 6. 测试工具与框架

### 6.1 后端测试工具
| 工具 | 用途 | 版本 |
|------|------|------|
| Go testing | 单元测试框架 | 标准库 |
| testify/assert | 断言库 | latest |
| testify/mock | Mock 框架 | latest |
| httptest | HTTP 测试工具 | 标准库 |
| go-cmp | 深度比较 | latest |
| race detector | 并发安全检测 | go test -race |

### 6.2 前端测试工具
| 工具 | 用途 | 版本 |
|------|------|------|
| Vitest | 单元测试框架 | 2.x |
| React Testing Library | 组件测试 | 16.x |
| @testing-library/user-event | 用户交互模拟 | 14.x |
| MSW | API Mock | 2.x |
| @vitest/coverage-v8 | 覆盖率报告 | 2.x |

### 6.3 E2E 测试工具
| 工具 | 用途 | 版本 |
|------|------|------|
| Playwright | E2E 自动化 | 1.48+ |
| @playwright/test | 测试运行器 | 1.48+ |

### 6.4 性能测试工具
| 工具 | 用途 | 版本 |
|------|------|------|
| k6 | 压测工具 | latest |
| pprof | Go 性能分析 | 标准库 |
| Lighthouse | 前端性能分析 | Chrome 内置 |

## 7. 测试覆盖率目标

### 7.1 代码覆盖率
| 模块 | 行覆盖率目标 | 分支覆盖率目标 |
|------|------------|--------------|
| 后端核心模块（engine/tool/session） | ≥ 80% | ≥ 75% |
| 后端 Handler 层 | ≥ 70% | ≥ 65% |
| 后端 Provider 层 | ≥ 60% | ≥ 55% |
| 前端 Hooks/Stores | ≥ 80% | ≥ 75% |
| 前端组件 | ≥ 70% | ≥ 65% |
| 前端工具函数 | ≥ 90% | ≥ 85% |

### 7.2 功能覆盖率
- **P0 功能**：100% 覆盖（单元 + 集成 + E2E）
- **P1 功能**：≥ 80% 覆盖（单元 + 集成）
- **P2 功能**：≥ 60% 覆盖（单元测试为主）

### 7.3 场景覆盖率
- **正常流程**：100% 覆盖（5个典型场景）
- **异常流程**：≥ 80% 覆盖（API 超时、输入错误、会话不存在等）
- **边界场景**：≥ 60% 覆盖（极端预算、极端时间、特殊字符等）

## 8. 测试执行计划

### 8.1 本地开发测试
```bash
# 后端单元测试
make test-backend

# 前端单元测试
make test-frontend

# 后端集成测试
make test-integration

# E2E 测试（需启动服务）
make test-e2e

# 覆盖率报告
make coverage
```

### 8.2 CI 自动化测试
**触发条件**：
- Push 到 feature 分支：运行单元测试 + 集成测试
- 创建 PR：运行全量测试（含 E2E）
- Merge 到 main：运行全量测试 + 生成覆盖率报告

**流水线步骤**：
1. 代码检查（golangci-lint、eslint）
2. 后端单元测试（go test -race -cover）
3. 前端单元测试（vitest run --coverage）
4. 后端集成测试（go test -tags=integration）
5. E2E 测试（playwright test）
6. 覆盖率上传（codecov）

### 8.3 测试报告
- **单元测试报告**：JUnit XML 格式，上传至 CI
- **覆盖率报告**：HTML 格式，本地查看；Codecov 在线查看
- **E2E 测试报告**：Playwright HTML Report，包含截图和视频
- **性能测试报告**：k6 HTML Report，包含响应时间分布和错误率

## 9. 测试维护策略

### 9.1 测试代码规范
- 测试文件命名：`*_test.go`（Go）、`*.test.ts(x)`（TypeScript）
- 测试函数命名：`Test<功能>_<场景>_<预期结果>`（如 `TestPlanner_CollectingToRecommending_Success`）
- 每个测试用例独立（不依赖其他测试的执行顺序）
- 使用 Table-Driven Tests 减少重复代码

### 9.2 Mock 数据维护
- Mock Provider 数据集中管理（`internal/provider/testdata/`）
- 新增城市对时同步更新 Mock 数据
- Mock 数据版本化（与产品需求文档对齐）

### 9.3 测试更新触发
- **需求变更**：更新对应的测试用例和断言
- **Bug 修复**：先写失败的测试用例，再修复代码
- **重构**：确保测试通过后再提交
- **新增功能**：TDD 模式，先写测试再实现

## 10. 特殊测试场景

### 10.1 LLM 不确定性测试
**策略**：关注结构而非具体文本
- 断言返回的 JSON 结构完整性
- 断言必填字段非空
- 断言数值范围合理（如预算不为负数）
- 使用正则表达式匹配关键词（如推荐理由包含"海边"）

**示例**：
```go
// 不推荐：断言具体文本
assert.Equal(t, "三亚是国内顶级海滨度假胜地", recommendation.Reason)

// 推荐：断言结构和关键词
assert.NotEmpty(t, recommendation.Reason)
assert.Contains(t, recommendation.Reason, "海")
assert.Len(t, recommendations, 3) // 推荐数量
```

### 10.2 SSE 流式响应测试
**策略**：模拟 SSE 事件序列
- 使用 httptest.ResponseRecorder 捕获 SSE 输出
- 解析 SSE 事件流，验证事件顺序
- 验证 delta 文本拼接正确性
- 验证 done 事件最后发送

**示例**：
```go
// 捕获 SSE 流
recorder := httptest.NewRecorder()
handler.ServeHTTP(recorder, req)

// 解析事件
events := parseSSEEvents(recorder.Body.String())
assert.Equal(t, "thinking", events[0].Type)
assert.Equal(t, "text", events[1].Type)
assert.Equal(t, "done", events[len(events)-1].Type)
```

### 10.3 并发安全测试
**策略**：使用 race detector + 并发压测
```bash
# 启用 race detector
go test -race ./internal/session/...

# 并发压测
go test -run TestSessionStore_Concurrent -count=100
```

**测试场景**：
- 多个 goroutine 同时读写同一 Session
- 多个 goroutine 同时创建不同 Session
- 并发删除和查询 Session

### 10.4 超时与降级测试
**策略**：Mock Provider 模拟超时
```go
// Mock Provider 返回超时错误
mockProvider.On("Search", mock.Anything, mock.Anything).
    Return(nil, context.DeadlineExceeded).
    Once()

// 验证降级触发
result := tool.Execute(ctx, args)
assert.True(t, result.IsMock)
assert.Equal(t, "fallback", result.Source)
```

## 11. 验收标准

### 11.1 测试通过标准
- [ ] 所有单元测试通过（后端 + 前端）
- [ ] 所有集成测试通过
- [ ] 所有 E2E 测试通过（至少在 Chrome 上）
- [ ] 代码覆盖率达到目标（后端核心 ≥ 80%，前端 Hooks ≥ 80%）
- [ ] 无 race condition 警告（go test -race）
- [ ] 性能测试达标（单轮响应 < 15s，推荐生成 < 20s）

### 11.2 发布门禁
⚠️ **以下条件必须全部满足，任一失败即阻塞发布**：
- [ ] 所有 P0 功能测试用例 100% 通过
- [ ] 集成测试通过率 ≥ 95%
- [ ] E2E 测试通过率 ≥ 90%（允许偶发性网络问题）
- [ ] 无 P0/P1 级别的已知 Bug
- [ ] 性能测试无回归（响应时间不超过基线 20%）
- [ ] 安全测试通过（输入校验、CORS、错误脱敏）

### 11.3 测试文档完整性
- [ ] 所有测试用例有清晰的描述和断言
- [ ] Mock 数据有文档说明（覆盖范围、更新策略）
- [ ] E2E 测试有截图或视频记录
- [ ] 性能测试有基线数据和趋势图
- [ ] 已知问题有 Issue 跟踪（标注影响范围和优先级）

---

**文档版本**：v1.0  
**最后更新**：2026-05-24  
**负责人**：测试设计师 Agent  
**审阅状态**：待审阅

**下一步**：请审阅本测试方案，确认后可由测试执行 agent 进行自动化执行与报告。
