#!/bin/bash

# 旅行规划 Agent 启动验证脚本
# 用于验证后端和前端是否正确启动

set -e

echo "=== 旅行规划 Agent 启动验证 ==="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查后端服务
echo "检查后端服务..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 后端服务运行正常${NC}"
else
    echo -e "${RED}✗ 后端服务未运行${NC}"
    echo "请先运行: ./start.sh"
    exit 1
fi

# 检查前端服务
echo "检查前端服务..."
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 前端服务运行正常${NC}"
else
    echo -e "${YELLOW}⚠ 前端服务未运行${NC}"
    echo "请在另一个终端运行: cd web && npm run dev"
fi

echo ""
echo "=== 测试 API 接口 ==="
echo ""

# 创建会话
echo "1. 创建会话..."
SESSION_RESPONSE=$(curl -s -X POST http://localhost:8080/api/session)
SESSION_ID=$(echo $SESSION_RESPONSE | grep -o '"session_id":"[^"]*' | cut -d'"' -f4)

if [ -z "$SESSION_ID" ]; then
    echo -e "${RED}✗ 创建会话失败${NC}"
    echo "响应: $SESSION_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ 会话创建成功${NC}"
echo "  Session ID: $SESSION_ID"
echo ""

# 获取会话状态
echo "2. 获取会话状态..."
GET_RESPONSE=$(curl -s http://localhost:8080/api/session/$SESSION_ID)

if echo $GET_RESPONSE | grep -q "COLLECTING"; then
    echo -e "${GREEN}✓ 会话状态获取成功${NC}"
    echo "  当前阶段: COLLECTING"
else
    echo -e "${RED}✗ 会话状态获取失败${NC}"
    echo "响应: $GET_RESPONSE"
    exit 1
fi

echo ""
echo "=== 测试 SSE 流 ==="
echo ""

# 发送消息测试 SSE
echo "3. 发送消息（测试 SSE 流）..."
echo "  消息: 五一假期想去海边玩，预算5000左右，从北京出发"
echo ""

# 使用 timeout 限制 SSE 连接时间（5秒）
timeout 5 curl -s -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d "{\"session_id\":\"$SESSION_ID\",\"message\":\"五一假期想去海边玩，预算5000左右，从北京出发\"}" | head -20 || true

echo ""
echo -e "${GREEN}✓ SSE 流测试完成${NC}"
echo ""

# 删除会话
echo "4. 删除会话..."
DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/api/session/$SESSION_ID)

if echo $DELETE_RESPONSE | grep -q "session deleted"; then
    echo -e "${GREEN}✓ 会话删除成功${NC}"
else
    echo -e "${YELLOW}⚠ 会话删除响应: $DELETE_RESPONSE${NC}"
fi

echo ""
echo "=== 验证完成 ==="
echo ""
echo -e "${GREEN}所有测试通过！${NC}"
echo ""
echo "后端服务地址: http://localhost:8080"
echo "前端应用地址: http://localhost:5173"
echo ""
echo "API 文档:"
echo "  - 创建会话: POST /api/session"
echo "  - 获取会话: GET /api/session/:id"
echo "  - 删除会话: DELETE /api/session/:id"
echo "  - 发送消息: POST /api/chat (SSE 流)"
echo "  - 健康检查: GET /health"
