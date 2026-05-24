# 测试用例 - 旅行规划 Agent

> ⚠️ **发布门禁**：所有 P0 用例必须 **全部通过**，任一失败即阻塞发布。

## 目录
- [1. 功能测试](#1-功能测试)
  - [1.1 需求收集阶段（COLLECTING）](#11-需求收集阶段collecting)
  - [1.2 目的地推荐阶段（RECOMMENDING）](#12-目的地推荐阶段recommending)
  - [1.3 行程规划阶段（PLANNING）](#13-行程规划阶段planning)
  - [1.4 预订引导阶段（BOOKING）](#14-预订引导阶段booking)
  - [1.5 会话管理](#15-会话管理)
- [2. 接口测试](#2-接口测试)
- [3. 安全测试](#3-安全测试)
- [4. 性能测试](#4-性能测试)
- [5. 边界测试](#5-边界测试)
- [6. 兼容性测试](#6-兼容性测试)

---

## 1. 功能测试

### 1.1 需求收集阶段（COLLECTING）

#### TC-COLLECT-001
- **用例名称**：完整需求一次性提供
- **关联需求**：product.md 3.1 节 - 目的地探索与推荐
- **前置条件**：创建新会话，Phase = COLLECTING
- **测试步骤**：
  1. 发送消息："五一假期想去海边玩，预算5000左右，从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 SSE 事件：thinking → text → phase_change(COLLECTING → RECOMMENDING) → done
  - Session.UserRequest 字段完整填充：
    - departure_city = "北京"
    - travel_start = "2026-05-01"
    - travel_end = "2026-05-05"
    - budget = 5000
    - preferences 包含 "海边"
  - 无追问消息
- **优先级**：P0

#### TC-COLLECT-002
- **用例名称**：缺失出发城市，Agent 追问
- **关联需求**：product.md 3.1 节 - 参数解析失败处理
- **前置条件**：创建新会话，Phase = COLLECTING
- **测试步骤**：
  1. 发送消息："五一想去海边，预算5000"
  2. 等待 Agent 响应
  3. 发送消息："从北京出发"
- **预期结果**：
  - 第一轮：收到 text 事件，内容包含"出发地"或"从哪里出发"
  - Phase 保持 COLLECTING
  - 第二轮：收到 phase_change(COLLECTING → RECOMMENDING)
  - Session.UserRequest.departure_city = "北京"
- **优先级**：P0

#### TC-COLLECT-003
- **用例名称**：缺失预算，Agent 追问
- **关联需求**：product.md 3.1 节 - 参数解析失败处理
- **前置条件**：创建新会话，Phase = COLLECTING
- **测试步骤**：
  1. 发送消息："五一想去三亚，从北京出发"
  2. 等待 Agent 响应
  3. 发送消息："预算5000"
- **预期结果**：
  - 第一轮：收到 text 事件，内容包含"预算"
  - Phase 保持 COLLECTING
  - 第二轮：收到 phase_change(COLLECTING → RECOMMENDING)
  - Session.UserRequest.budget = 5000
- **优先级**：P0

#### TC-COLLECT-004
- **用例名称**：模糊时间表达解析
- **关联需求**：product.md 3.1 节 - 输入处理
- **前置条件**：创建新会话，Phase = COLLECTING
- **测试步骤**：
  1. 发送消息："清明前后想去海边，预算3000，从上海出发"
  2. 等待 Agent 响应
- **预期结果**：
  - Session.UserRequest.travel_start 在 2026-04-03 ~ 2026-04-06 范围内
  - Session.UserRequest.travel_end 晚于 travel_start
  - 收到 phase_change(COLLECTING → RECOMMENDING)
- **优先级**：P1

#### TC-COLLECT-005
- **用例名称**：偏好标签提取
- **关联需求**：product.md 3.1 节 - 输入处理
- **前置条件**：创建新会话，Phase = COLLECTING
- **测试步骤**：
  1. 发送消息："五一想去海边，人少的地方，可以潜水，预算5000，从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - Session.UserRequest.preferences 包含 ["海边", "人少", "潜水"]（顺序不限）
  - 收到 phase_change(COLLECTING → RECOMMENDING)
- **优先级**：P1

### 1.2 目的地推荐阶段（RECOMMENDING）

#### TC-RECOMMEND-001
- **用例名称**：生成 2-3 个推荐目的地
- **关联需求**：product.md 3.1 节 - 输出
- **前置条件**：
  - Phase = RECOMMENDING
  - UserRequest 已完整填充（北京出发，五一，预算5000，海边）
- **测试步骤**：
  1. 触发推荐生成
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 tool_call 事件（flight_search/hotel_search/weather_query）
  - 收到 card 事件，card_type = "destination_recommendation"
  - card.data 包含 2-3 个 Recommendation 对象
  - 每个 Recommendation 包含：
    - destination（非空）
    - reason（非空）
    - weather（temp_range/condition/rain_prob）
    - transport（mode/duration/price）
    - hotel_range（economy/comfort）
    - total_estimate（min/max）
    - sample_tickets（至少1个）
  - 收到 done 事件
- **优先级**：P0

#### TC-RECOMMEND-002
- **用例名称**：推荐卡片数据合理性
- **关联需求**：product.md 3.1 节 - 输出
- **前置条件**：同 TC-RECOMMEND-001
- **测试步骤**：
  1. 触发推荐生成
  2. 验证返回的推荐数据
- **预期结果**：
  - 每个推荐的 total_estimate.min <= total_estimate.max
  - total_estimate.max <= UserRequest.budget * 1.1（允许10%浮动）
  - transport.price > 0
  - hotel_range.economy < hotel_range.comfort
  - weather.rain_prob 在 0.0-1.0 范围内
  - sample_tickets[0].price > 0
- **优先级**：P0

#### TC-RECOMMEND-003
- **用例名称**：用户选择目的地
- **关联需求**：product.md 3.2 节 - 触发入口
- **前置条件**：
  - Phase = RECOMMENDING
  - 已生成推荐卡片（3个目的地）
- **测试步骤**：
  1. 发送消息："选第2个"或"去三亚"
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 phase_change(RECOMMENDING → PLANNING)
  - Session 记录选定的目的地
  - 开始生成行程表
- **优先级**：P0

#### TC-RECOMMEND-004
- **用例名称**：用户要求重新推荐
- **关联需求**：design.md 3.2 节 - 阶段流转
- **前置条件**：
  - Phase = RECOMMENDING
  - 已生成推荐卡片
- **测试步骤**：
  1. 发送消息："这些都不太满意，换几个"
  2. 等待 Agent 响应
- **预期结果**：
  - Phase 保持 RECOMMENDING
  - 收到新的 card 事件（不同的推荐）
  - 收到 done 事件
- **优先级**：P1

#### TC-RECOMMEND-005
- **用例名称**：API 降级处理
- **关联需求**：product.md 3.1 节 - 异常处理
- **前置条件**：
  - Phase = RECOMMENDING
  - Mock Provider 配置为超时
- **测试步骤**：
  1. 触发推荐生成
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 tool_call 事件，status = "error"
  - 仍然收到 card 事件（使用 Mock 数据）
  - card.data[0].is_mock = true
  - card.data[0].weather.is_mock = true
  - 收到 text 事件，包含"预估数据"或"仅供参考"
  - 收到 done 事件
- **优先级**：P0

### 1.3 行程规划阶段（PLANNING）

#### TC-PLAN-001
- **用例名称**：生成完整行程表
- **关联需求**：product.md 3.2 节 - 输出
- **前置条件**：
  - Phase = PLANNING
  - 已选定目的地（三亚）
  - UserRequest 完整（北京出发，2026-05-01 ~ 2026-05-05，预算5000）
- **测试步骤**：
  1. 触发行程生成
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 tool_call 事件（flight_search/hotel_search）
  - 收到 card 事件，card_type = "itinerary"
  - card.data 包含完整 Itinerary 对象：
    - outbound（去程交通，包含 type/number/depart_time/arrive_time/price）
    - hotel（住宿，包含 name/location/check_in/check_out/nights/price_per_night）
    - daily_plans（每日行程，至少4天）
    - return（返程交通）
    - cost_summary（费用汇总，6个分类）
    - total_cost（总费用）
  - total_cost <= UserRequest.budget * 1.1
  - 收到 done 事件
- **优先级**：P0

#### TC-PLAN-002
- **用例名称**：每日行程合理性
- **关联需求**：product.md 3.2 节 - 数据约束
- **前置条件**：同 TC-PLAN-001
- **测试步骤**：
  1. 触发行程生成
  2. 验证每日行程数据
- **预期结果**：
  - 每日行程不超过 3 个主要景点/活动
  - 每个 Activity 包含 time/name/cost
  - Activity.cost >= 0
  - 每日 activities 按时间顺序排列
  - DayPlan.date 连续且在 travel_start ~ travel_end 范围内
- **优先级**：P0

#### TC-PLAN-003
- **用例名称**：费用汇总准确性
- **关联需求**：product.md 3.2 节 - 输出
- **前置条件**：同 TC-PLAN-001
- **测试步骤**：
  1. 触发行程生成
  2. 验证费用汇总
- **预期结果**：
  - cost_summary.outbound_cost = outbound.price
  - cost_summary.return_cost = return.price
  - cost_summary.hotel_cost = hotel.total_price
  - cost_summary.activity_cost = sum(所有 Activity.cost)
  - total_cost = sum(cost_summary 所有字段)
- **优先级**：P0

#### TC-PLAN-004
- **用例名称**：用户确认行程
- **关联需求**：product.md 3.3 节 - 触发入口
- **前置条件**：
  - Phase = PLANNING
  - 已生成完整行程表
- **测试步骤**：
  1. 发送消息："确认行程"或"可以"
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 phase_change(PLANNING → BOOKING)
  - Session.Itinerary 已保存
  - 开始预订引导流程
- **优先级**：P0

#### TC-PLAN-005
- **用例名称**：用户要求修改行程 - 换酒店
- **关联需求**：product.md 3.4 节 - 行程修改与调整
- **前置条件**：
  - Phase = PLANNING
  - 已生成完整行程表
- **测试步骤**：
  1. 发送消息："换个五星级酒店"
  2. 等待 Agent 响应
- **预期结果**：
  - Phase 保持 PLANNING
  - 收到新的 card 事件（itinerary）
  - card.data.hotel.name 不同于原酒店
  - card.data.hotel.price_per_night > 原价格
  - card.data.total_cost 更新
  - 其他部分（去程/返程/每日行程）保持不变
- **优先级**：P1

#### TC-PLAN-006
- **用例名称**：用户要求修改行程 - 换景点
- **关联需求**：product.md 3.4 节 - 行程修改与调整
- **前置条件**：
  - Phase = PLANNING
  - 已生成完整行程表
- **测试步骤**：
  1. 发送消息："把第二天的景点换掉"
  2. 等待 Agent 响应
- **预期结果**：
  - Phase 保持 PLANNING
  - 收到新的 card 事件（itinerary）
  - card.data.daily_plans[1].activities 不同于原行程
  - 其他天的行程保持不变
  - total_cost 更新
- **优先级**：P1

#### TC-PLAN-007
- **用例名称**：超预算提示
- **关联需求**：product.md 3.2 节 - 异常处理
- **前置条件**：
  - Phase = PLANNING
  - UserRequest.budget = 3000（较低预算）
  - 选定目的地为三亚（高消费）
- **测试步骤**：
  1. 触发行程生成
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 card 事件（itinerary）
  - card.data.total_cost > UserRequest.budget
  - 收到 text 事件，包含"超出预算"或"超预算"
  - 提供降级建议（如换经济型酒店）
- **优先级**：P0

### 1.4 预订引导阶段（BOOKING）

#### TC-BOOK-001
- **用例名称**：按顺序输出预订链接
- **关联需求**：product.md 3.3 节 - 预订顺序
- **前置条件**：
  - Phase = BOOKING
  - 已确认完整行程表
- **测试步骤**：
  1. 触发预订引导
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 card 事件，card_type = "booking_item"
  - card.data.order = 1
  - card.data.item_type = "OUTBOUND"（去程交通）
  - card.data 包含：
    - item_name（非空）
    - detail（包含班次号和时间）
    - platform（如"携程"）
    - link（HTTPS 链接）
    - alt_links（至少1个备选平台）
    - price（> 0）
    - status = "PENDING"
  - 收到 text 事件，提示用户点击链接预订
  - 收到 done 事件
- **优先级**：P0

#### TC-BOOK-002
- **用例名称**：用户确认已预订
- **关联需求**：product.md 3.3 节 - 操作流程
- **前置条件**：
  - Phase = BOOKING
  - 已输出第1项预订链接（去程）
- **测试步骤**：
  1. 发送消息："已订"或"搞定"
  2. 等待 Agent 响应
- **预期结果**：
  - Session.BookingItems[0].status = "BOOKED"
  - 收到 card 事件，card_type = "booking_item"
  - card.data.order = 2
  - card.data.item_type = "HOTEL"（住宿）
  - 收到 done 事件
- **优先级**：P0

#### TC-BOOK-003
- **用例名称**：用户跳过预订项
- **关联需求**：product.md 3.3 节 - 操作流程
- **前置条件**：
  - Phase = BOOKING
  - 已输出第2项预订链接（住宿）
- **测试步骤**：
  1. 发送消息："跳过"或"不订了"
  2. 等待 Agent 响应
- **预期结果**：
  - Session.BookingItems[1].status = "SKIPPED"
  - 收到 card 事件，card_type = "booking_item"
  - card.data.order = 3
  - card.data.item_type = "RETURN"（返程）
  - 收到 done 事件
- **优先级**：P0

#### TC-BOOK-004
- **用例名称**：全部预订完成，输出汇总
- **关联需求**：product.md 3.3 节 - 操作流程
- **前置条件**：
  - Phase = BOOKING
  - 已处理所有预订项（4项：去程/住宿/返程/门票）
- **测试步骤**：
  1. 最后一项确认"已订"
  2. 等待 Agent 响应
- **预期结果**：
  - 收到 phase_change(BOOKING → COMPLETED)
  - 收到 card 事件，card_type = "booking_summary"
  - card.data 包含：
    - destination（目的地）
    - items（所有 BookingItem）
    - total_booked（已订项总花费）
    - total_skipped（跳过项数量）
    - booked_count（已订项数量）
  - 收到 text 事件，祝旅途愉快
  - 收到 done 事件
- **优先级**：P0

#### TC-BOOK-005
- **用例名称**：预订链接格式验证
- **关联需求**：product.md 3.3 节 - 数据约束
- **前置条件**：
  - Phase = BOOKING
  - 已输出预订链接
- **测试步骤**：
  1. 验证 BookingItem.link 格式
- **预期结果**：
  - link 以 "https://" 开头
  - link 包含 URL 编码的参数
  - alt_links 中每个 link 也符合 HTTPS 格式
  - platform 为已知平台名称（携程/飞猪/去哪儿）
- **优先级**：P0

### 1.5 会话管理

#### TC-SESSION-001
- **用例名称**：创建会话
- **关联需求**：specs.md 2.3 节 - 创建会话接口
- **前置条件**：无
- **测试步骤**：
  1. 调用 POST /api/session
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 200
  - 响应 JSON：
    - code = 0
    - message = "success"
    - data.session_id 为 UUID v4 格式
    - data.phase = "COLLECTING"
    - data.created_at 为 RFC3339 格式时间戳
- **优先级**：P0

#### TC-SESSION-002
- **用例名称**：获取会话状态
- **关联需求**：specs.md 2.3 节 - 获取会话接口
- **前置条件**：已创建会话（session_id = "test-uuid"）
- **测试步骤**：
  1. 调用 GET /api/session/test-uuid
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 200
  - 响应 JSON：
    - code = 0
    - data 为完整 Session 对象
    - data.id = "test-uuid"
    - data.phase 为有效 Phase 枚举值
    - data.chat_history 为数组
- **优先级**：P0

#### TC-SESSION-003
- **用例名称**：获取不存在的会话
- **关联需求**：specs.md 2.3 节 - 错误响应
- **前置条件**：无
- **测试步骤**：
  1. 调用 GET /api/session/non-existent-uuid
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 404
  - 响应 JSON：
    - code = 10002
    - message = "session not found"
    - data = null
- **优先级**：P0

#### TC-SESSION-004
- **用例名称**：删除会话
- **关联需求**：specs.md 2.3 节 - 删除会话接口
- **前置条件**：已创建会话（session_id = "test-uuid"）
- **测试步骤**：
  1. 调用 DELETE /api/session/test-uuid
  2. 验证响应
  3. 再次调用 GET /api/session/test-uuid
- **预期结果**：
  - DELETE 响应：HTTP 200，code = 0，message = "session deleted"
  - GET 响应：HTTP 404，code = 10002
- **优先级**：P0

#### TC-SESSION-005
- **用例名称**：会话状态持久化
- **关联需求**：design.md 4.2 节 - 状态持久化
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息，触发阶段流转（COLLECTING → RECOMMENDING）
  2. 调用 GET /api/session/:id 获取状态
  3. 验证 Phase 已更新
- **预期结果**：
  - Session.Phase = "RECOMMENDING"
  - Session.UpdatedAt 晚于 CreatedAt
  - Session.UserRequest 已填充
- **优先级**：P0


## 2. 接口测试

### 2.1 SSE 事件流测试

#### TC-API-001
- **用例名称**：SSE 事件类型完整性
- **关联需求**：specs.md 2.3 节 - SSE 事件类型
- **前置条件**：已创建会话
- **测试步骤**：
  1. 调用 POST /api/chat，发送完整需求
  2. 解析 SSE 事件流
- **预期结果**：
  - 收到至少以下事件类型：thinking, text, phase_change, done
  - 每个事件格式符合 SSE 规范（event: xxx\ndata: {...}\n\n）
  - data 字段为有效 JSON
  - 最后一个事件为 done
- **优先级**：P0

#### TC-API-002
- **用例名称**：流式文本增量推送
- **关联需求**：specs.md 2.3 节 - text 事件
- **前置条件**：已创建会话
- **测试步骤**：
  1. 调用 POST /api/chat
  2. 收集所有 text 事件
- **预期结果**：
  - 收到多个 text 事件，delta = true
  - 拼接所有 content 字段，形成完整文本
  - 完整文本语义连贯
- **优先级**：P0

#### TC-API-003
- **用例名称**：tool_call 事件状态流转
- **关联需求**：specs.md 2.3 节 - tool_call 事件
- **前置条件**：Phase = RECOMMENDING
- **测试步骤**：
  1. 调用 POST /api/chat
  2. 监听 tool_call 事件
- **预期结果**：
  - 收到 tool_call 事件，status = "running"，result = null
  - 收到 tool_call 事件（同一 name），status = "done"，result 非空
  - result 包含有效数据（如 tickets 数组）
- **优先级**：P0

#### TC-API-004
- **用例名称**：card 事件数据完整性
- **关联需求**：specs.md 2.3 节 - card 事件
- **前置条件**：Phase = RECOMMENDING
- **测试步骤**：
  1. 调用 POST /api/chat
  2. 等待 card 事件
- **预期结果**：
  - 收到 card 事件，card_type = "destination_recommendation"
  - data 为完整 JSON 对象（非字符串）
  - data 符合 Recommendation 数组结构
  - 每个元素包含所有必填字段
- **优先级**：P0

### 2.2 错误码测试

#### TC-API-005
- **用例名称**：请求参数校验失败 - 缺失 session_id
- **关联需求**：specs.md 2.2 节 - 错误码 10001
- **前置条件**：无
- **测试步骤**：
  1. 调用 POST /api/chat，body = {"message": "test"}（缺失 session_id）
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 400
  - 响应 JSON：
    - code = 10001
    - message 包含 "session_id"
    - data = null
- **优先级**：P0

#### TC-API-006
- **用例名称**：会话不存在
- **关联需求**：specs.md 2.2 节 - 错误码 10002
- **前置条件**：无
- **测试步骤**：
  1. 调用 POST /api/chat，session_id = "non-existent-uuid"
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 404
  - 响应 JSON：
    - code = 10002
    - message = "session not found"
- **优先级**：P0

#### TC-API-007
- **用例名称**：消息内容为空
- **关联需求**：specs.md 2.2 节 - 错误码 10003
- **前置条件**：已创建会话
- **测试步骤**：
  1. 调用 POST /api/chat，message = ""
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 422
  - 响应 JSON：
    - code = 10003
    - message 包含 "empty" 或 "空"
- **优先级**：P0

#### TC-API-008
- **用例名称**：消息内容过长
- **关联需求**：specs.md 2.2 节 - 错误码 10003
- **前置条件**：已创建会话
- **测试步骤**：
  1. 调用 POST /api/chat，message = "a" * 2001（2001字符）
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 422
  - 响应 JSON：
    - code = 10003
    - message 包含 "too long" 或 "过长"
- **优先级**：P0

#### TC-API-009
- **用例名称**：LLM 服务不可用
- **关联需求**：specs.md 2.2 节 - 错误码 20002
- **前置条件**：
  - 已创建会话
  - Mock LLM Client 返回超时错误
- **测试步骤**：
  1. 调用 POST /api/chat
  2. 等待 SSE 事件
- **预期结果**：
  - 收到 error 事件：
    - code = "20002"
    - message 包含 "AI 服务" 或 "不可用"
  - 收到 done 事件
- **优先级**：P0

## 3. 安全测试

### 3.1 输入校验

#### TC-SEC-001
- **用例名称**：SQL 注入防护
- **关联需求**：安全需求
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："' OR '1'='1"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无 SQL 错误
  - 响应内容不包含数据库错误信息
  - Session 状态正常
- **优先级**：P0

#### TC-SEC-002
- **用例名称**：XSS 防护
- **关联需求**：安全需求
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："<script>alert('xss')</script>"
  2. 等待 Agent 响应
  3. 前端渲染响应内容
- **预期结果**：
  - 后端正常处理，无错误
  - 前端渲染时，脚本标签被转义或过滤
  - 浏览器不执行 alert
- **优先级**：P0

#### TC-SEC-003
- **用例名称**：Session ID 格式校验
- **关联需求**：specs.md 2.3 节 - 路径参数校验
- **前置条件**：无
- **测试步骤**：
  1. 调用 GET /api/session/invalid-format
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 400 或 404
  - 响应 JSON：
    - code = 10001 或 10002
    - message 包含 "invalid" 或 "格式"
- **优先级**：P0

#### TC-SEC-004
- **用例名称**：请求体大小限制
- **关联需求**：design.md 8.2 节 - 输入校验
- **前置条件**：已创建会话
- **测试步骤**：
  1. 调用 POST /api/chat，body 大小 > 10KB
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 413（Payload Too Large）或 400
  - 请求被拒绝
- **优先级**：P1

### 3.2 CORS 配置

#### TC-SEC-005
- **用例名称**：允许配置的前端域名
- **关联需求**：design.md 8.3 节 - CORS 限制
- **前置条件**：配置 cors_origins = ["http://localhost:5173"]
- **测试步骤**：
  1. 从 http://localhost:5173 发起 POST /api/chat 请求
  2. 检查响应头
- **预期结果**：
  - 响应头包含：
    - Access-Control-Allow-Origin: http://localhost:5173
    - Access-Control-Allow-Methods: GET, POST, DELETE, OPTIONS
  - 请求成功
- **优先级**：P0

#### TC-SEC-006
- **用例名称**：拒绝未授权来源
- **关联需求**：design.md 8.3 节 - CORS 限制
- **前置条件**：配置 cors_origins = ["http://localhost:5173"]
- **测试步骤**：
  1. 从 http://evil.com 发起 POST /api/chat 请求
  2. 检查响应
- **预期结果**：
  - 响应头不包含 Access-Control-Allow-Origin: http://evil.com
  - 浏览器阻止请求（CORS 错误）
- **优先级**：P0

### 3.3 错误信息脱敏

#### TC-SEC-007
- **用例名称**：不暴露内部堆栈
- **关联需求**：design.md 8.3 节 - 错误信息脱敏
- **前置条件**：触发服务器内部错误（如 panic）
- **测试步骤**：
  1. 调用 API 触发错误
  2. 检查响应内容
- **预期结果**：
  - 响应 JSON 不包含：
    - Go 堆栈信息（goroutine、panic）
    - 文件路径（/Users/xxx/project/...）
    - 内部变量名
  - 仅返回用户友好的错误信息
- **优先级**：P0

#### TC-SEC-008
- **用例名称**：不泄露配置细节
- **关联需求**：design.md 8.3 节 - 错误信息脱敏
- **前置条件**：LLM API Key 配置错误
- **测试步骤**：
  1. 调用 POST /api/chat
  2. 检查响应内容
- **预期结果**：
  - 响应不包含：
    - API Key 值
    - API 端点 URL
    - 配置文件路径
  - 仅返回"服务配置错误"类提示
- **优先级**：P0

## 4. 性能测试

### 4.1 响应时间

#### TC-PERF-001
- **用例名称**：单轮响应时间 < 15秒
- **关联需求**：product.md 7.1 节 - 性能指标
- **前置条件**：
  - 已创建会话
  - 使用 Mock Provider（避免真实 API 延迟）
- **测试步骤**：
  1. 发送消息："五一想去海边，预算5000，从北京出发"
  2. 记录从请求发送到收到 done 事件的时间
- **预期结果**：
  - 总耗时 < 15秒
  - 首个 text 事件在 3秒内到达
- **优先级**：P0

#### TC-PERF-002
- **用例名称**：推荐生成时间 < 20秒
- **关联需求**：product.md 7.1 节 - 性能指标
- **前置条件**：
  - Phase = RECOMMENDING
  - 使用 Mock Provider
- **测试步骤**：
  1. 触发推荐生成
  2. 记录从请求到收到 card 事件的时间
- **预期结果**：
  - 总耗时 < 20秒
  - tool_call 事件在 5秒内开始
- **优先级**：P0

#### TC-PERF-003
- **用例名称**：行程生成时间 < 15秒
- **关联需求**：product.md 7.1 节 - 性能指标
- **前置条件**：
  - Phase = PLANNING
  - 使用 Mock Provider
- **测试步骤**：
  1. 触发行程生成
  2. 记录从请求到收到 itinerary card 的时间
- **预期结果**：
  - 总耗时 < 15秒
- **优先级**：P0

#### TC-PERF-004
- **用例名称**：本地存储读写 < 100ms
- **关联需求**：product.md 7.1 节 - 性能指标
- **前置条件**：已创建会话
- **测试步骤**：
  1. 调用 GET /api/session/:id
  2. 记录响应时间
- **预期结果**：
  - 响应时间 < 100ms
- **优先级**：P1

### 4.2 并发测试

#### TC-PERF-005
- **用例名称**：10 并发会话处理
- **关联需求**：design.md 11 节 - 性能考量
- **前置条件**：无
- **测试步骤**：
  1. 创建 10 个会话
  2. 并发发送消息到 10 个会话
  3. 等待所有响应完成
- **预期结果**：
  - 所有会话正常响应
  - 无 race condition 错误
  - 平均响应时间 < 20秒
  - 无会话数据混淆
- **优先级**：P0

#### TC-PERF-006
- **用例名称**：SSE 长连接稳定性
- **关联需求**：design.md 11 节 - 性能考量
- **前置条件**：已创建会话
- **测试步骤**：
  1. 建立 SSE 连接
  2. 保持连接 5 分钟
  3. 期间发送 3 条消息
- **预期结果**：
  - 连接保持稳定，无断开
  - 所有消息正常响应
  - 无内存泄漏
- **优先级**：P1

### 4.3 资源占用

#### TC-PERF-007
- **用例名称**：单会话内存占用 < 100KB
- **关联需求**：design.md 11 节 - 性能考量
- **前置条件**：无
- **测试步骤**：
  1. 创建会话
  2. 完成完整流程（COLLECTING → COMPLETED）
  3. 测量会话对象内存占用
- **预期结果**：
  - 单个 Session 对象 < 100KB
  - 包含完整对话历史和行程数据
- **优先级**：P1

#### TC-PERF-008
- **用例名称**：100 会话总内存 < 10MB
- **关联需求**：design.md 11 节 - 性能考量
- **前置条件**：无
- **测试步骤**：
  1. 创建 100 个会话
  2. 每个会话完成部分流程
  3. 测量总内存占用
- **预期结果**：
  - 总内存占用 < 10MB
  - 无内存泄漏
- **优先级**：P1

## 5. 边界测试

### 5.1 极端输入

#### TC-EDGE-001
- **用例名称**：极低预算（500元）
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："五一想去海边，预算500，从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无错误
  - 收到 text 事件，提示预算不足或建议周边游
  - 或推荐低成本目的地（如周边海滩）
- **优先级**：P1

#### TC-EDGE-002
- **用例名称**：超高预算（50000元）
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："五一想去海边，预算50000，从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无错误
  - 推荐高端目的地（如马尔代夫、巴厘岛）
  - 行程包含豪华酒店和高端活动
- **优先级**：P2

#### TC-EDGE-003
- **用例名称**：当天出发
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："今天想去海边，预算5000，从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无错误
  - 提示当天出发票务紧张
  - 推荐周边目的地或次日出发
- **优先级**：P2

#### TC-EDGE-004
- **用例名称**：1年后出发
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："2027年5月想去海边，预算5000，从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无错误
  - 提示时间较远，数据为预估
  - 正常生成推荐
- **优先级**：P2

#### TC-EDGE-005
- **用例名称**：跨年行程
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："12月30日到1月3日想去海边，预算5000，从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无错误
  - travel_start = 2026-12-30
  - travel_end = 2027-01-03
  - 每日行程跨年正确
- **优先级**：P2

### 5.2 特殊字符

#### TC-EDGE-006
- **用例名称**：消息包含 emoji
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："五一想去海边🏖️，预算5000💰，从北京出发✈️"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无错误
  - 正确提取参数
  - 响应内容正常
- **优先级**：P2

#### TC-EDGE-007
- **用例名称**：消息包含换行符
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："五一想去海边\n预算5000\n从北京出发"
  2. 等待 Agent 响应
- **预期结果**：
  - 正常处理，无错误
  - 正确提取参数
- **优先级**：P2

#### TC-EDGE-008
- **用例名称**：消息全为空格
- **关联需求**：边界场景
- **前置条件**：已创建会话
- **测试步骤**：
  1. 发送消息："     "（5个空格）
  2. 验证响应
- **预期结果**：
  - HTTP 状态码 422
  - code = 10003
  - message 包含 "empty" 或 "空"
- **优先级**：P1

## 6. 兼容性测试

### 6.1 浏览器兼容性

#### TC-COMPAT-001
- **用例名称**：Chrome 120+ 完整流程
- **关联需求**：design.md 1.5 节 - 兼容性
- **前置条件**：Chrome 120+
- **测试步骤**：
  1. 打开前端应用
  2. 完成完整旅行规划流程
- **预期结果**：
  - 所有功能正常
  - SSE 连接稳定
  - 卡片渲染正确
  - 无控制台错误
- **优先级**：P0

#### TC-COMPAT-002
- **用例名称**：Safari 17+ 完整流程
- **关联需求**：design.md 1.5 节 - 兼容性
- **前置条件**：Safari 17+
- **测试步骤**：
  1. 打开前端应用
  2. 完成完整旅行规划流程
- **预期结果**：
  - 所有功能正常
  - SSE 连接稳定
  - 卡片渲染正确
  - 无控制台错误
- **优先级**：P0

#### TC-COMPAT-003
- **用例名称**：Firefox 120+ 完整流程
- **关联需求**：design.md 1.5 节 - 兼容性
- **前置条件**：Firefox 120+
- **测试步骤**：
  1. 打开前端应用
  2. 完成完整旅行规划流程
- **预期结果**：
  - 所有功能正常
  - SSE 连接稳定
  - 卡片渲染正确
  - 无控制台错误
- **优先级**：P1

### 6.2 屏幕尺寸

#### TC-COMPAT-004
- **用例名称**：移动端 375px（iPhone SE）
- **关联需求**：product.md 11.2 节 - 设计规范
- **前置条件**：浏览器窗口宽度 = 375px
- **测试步骤**：
  1. 打开前端应用
  2. 完成完整流程
- **预期结果**：
  - 布局正常，无横向滚动条
  - 卡片自适应宽度
  - 按钮可点击
  - 文字可读
- **优先级**：P0

#### TC-COMPAT-005
- **用例名称**：桌面端 1920px
- **关联需求**：product.md 11.2 节 - 设计规范
- **前置条件**：浏览器窗口宽度 = 1920px
- **测试步骤**：
  1. 打开前端应用
  2. 完成完整流程
- **预期结果**：
  - 内容居中，max-width 420px
  - 布局美观
  - 无拉伸变形
- **优先级**：P1

---

**测试用例总数**：60+  
**P0 用例数**：45  
**P1 用例数**：12  
**P2 用例数**：3+  

**覆盖率统计**：
- 功能模块覆盖：5/5（100%）
- 接口覆盖：4/4（100%）
- 错误码覆盖：5/5（100%）
- 安全场景覆盖：8/8（100%）
- 性能指标覆盖：4/4（100%）

**下一步**：
1. 审阅测试用例，确认覆盖完整性
2. 实现测试代码（单元测试 + 集成测试 + E2E 测试）
3. 配置 CI/CD 自动化测试流水线
4. 执行测试并生成覆盖率报告

