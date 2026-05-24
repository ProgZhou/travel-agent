#!/bin/bash

# 旅行规划 Agent 启动脚本

set -e

echo "=== 旅行规划 Agent 启动脚本 ==="
echo ""

# 检查 Go 版本
if ! command -v go &> /dev/null; then
    echo "错误: 未安装 Go"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✓ Go 版本: $GO_VERSION"

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "错误: 未找到 config.yaml 文件"
    echo "请复制 config.example.yaml 为 config.yaml 并配置 API Key"
    exit 1
fi

echo "✓ 配置文件已找到"

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
echo "前端地址: http://localhost:5173"
echo ""
echo "按 Ctrl+C 停止服务"
echo ""

./bin/server

