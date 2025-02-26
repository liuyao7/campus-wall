#!/bin/bash

# 数据库配置
DB_USER="campus_wall_user"
DB_PASS="your_password"
DB_NAME="campus_wall"

# 函数：创建数据库和用户
create_db() {
    mysql -u root -p <<EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASS';
GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
EOF
    echo "Database and user created successfully"
}

# 函数：删除数据库和用户
drop_db() {
    mysql -u root -p <<EOF
DROP DATABASE IF EXISTS $DB_NAME;
DROP USER IF EXISTS '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
EOF
    echo "Database and user dropped successfully"
}

# 主命令处理
case "$1" in
    "create")
        create_db
        ;;
    "drop")
        drop_db
        ;;
    *)
        echo "Usage: $0 {create|drop}"
        exit 1
        ;;
esac