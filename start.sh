#!/bin/bash

# 旅行规划 Agent 启动脚本

set -e

echo "=== 旅行规划 Agent 启动脚本 ==="
echo ""

# 检查 API Key
if [ -z "$ARK_API_KEY" ]; then
    echo "错误: 未设置 ARK_API_KEY 环境变量"
    echo "请运行: export ARK_API_KEY=\"your-api-key\""
    exit 1
fi

echo "✓ API Key 已设置"

# 检查 Go 版本
if ! command -v go &> /dev/null; then
    echo "错误: 未安装 Go"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✓ Go 版本: $GO_VERSION"

# 构建服务
echo ""
echo "正在构建服务..."
go build -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✓ 构建成功"
else
    echo "✗ 构建失败"
    exit 1
fi

# 启动服务
echo ""
echo "正在启动服务..."
echo "服务地址: http://localhost:8080"
echo "健康检查: http://localhost:8080/health"
echo ""
echo "按 Ctrl+C 停止服务"
echo ""

./bin/server
