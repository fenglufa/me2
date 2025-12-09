#!/bin/bash
# 笔记服务 MVP 快速设置脚本

echo "正在设置 Note Service MVP..."

# 1. 初始化数据库
echo "1. 请手动执行以下命令初始化数据库:"
echo "   mysql -u root -p123456 me2 < rpc/note.sql"

# 2. 安装依赖
echo "2. 安装依赖..."
go mod edit -require=github.com/me2/ai@v0.0.0
go mod edit -replace=github.com/me2/ai=../ai
go get github.com/zeromicro/go-zero@latest
go get google.golang.org/protobuf@latest
go get google.golang.org/grpc@latest
go mod tidy

echo "设置完成！"
echo "下一步:"
echo "1. 手动初始化数据库"
echo "2. 补充实现业务逻辑文件"
echo "3. 运行: make run-dev"
