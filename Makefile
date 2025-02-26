.PHONY: dev build test clean

# 开发环境运行
dev:
	bash scripts/dev.sh

# 构建
build:
	go build -o bin/app cmd/main.go

# 运行测试
test:
	go test ./...

# 清理
clean:
	rm -rf bin/


# 创建数据库
db-create:
	bash scripts/db.sh create

# 删除数据库
db-drop:
	bash scripts/db.sh drop