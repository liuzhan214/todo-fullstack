#!/bin/bash
# backend/build.sh — 编译 backend 并将产物输出到 dist/ 目录
#
# 用法: cd backend && ./build.sh
# 产物: dist/main        (可执行文件)
#       dist/etc/        (配置文件目录，运行时需要)
set -euo pipefail

# 切换到脚本所在目录（backend/），确保路径不依赖执行位置
cd "$(dirname "$0")"

echo "Building backend..."

# 编译当前 package（todo.go 是入口文件），输出到 dist/main
go build -o dist/backend/main .

# 复制配置文件目录到 dist/backend/etc/
# 运行时: ./dist/backend/main -f dist/backend/etc/todo.yaml
mkdir -p dist/backend/etc
cp -r etc/* dist/backend/etc/

echo "Done. Binary: dist/backend/main"
ls -lh dist/backend/main
