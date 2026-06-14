#!/bin/bash
# gateway/build.sh — 编译 gateway 并将产物输出到 dist/ 目录
#
# 用法: cd gateway && ./build.sh
# 产物: dist/main        (可执行文件)
#       dist/etc/        (配置文件目录，运行时需要)
set -euo pipefail

# 切换到脚本所在目录（gateway/），确保路径不依赖执行位置
cd "$(dirname "$0")"

echo "Building gateway..."

# 编译当前 package（todo.go 是入口文件），输出到 dist/main
go build -o dist/gateway/main .

# 复制配置文件目录到 dist/gateway/etc/
# 运行时: ./dist/gateway/main -f dist/gateway/etc/todo-api.yaml
mkdir -p dist/gateway/etc
cp -r etc/* dist/gateway/etc/

echo "Done. Binary: dist/gateway/main"
ls -lh dist/gateway/main
