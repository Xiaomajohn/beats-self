#!/bin/bash
# AuditMetricbeat 构建脚本

set -e

echo "========================================="
echo "  AuditMetricbeat Build Script"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到 Go 编译器${NC}"
    echo "请先安装 Go: https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓${NC} Go 版本: $(go version)"

# 获取项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 构建函数
build() {
    echo ""
    echo -e "${YELLOW}开始构建 AuditMetricbeat...${NC}"
    echo ""
    
    # 设置输出文件名
    OUTPUT="auditmetricbeat"
    if [ "$1" != "" ]; then
        OUTPUT="$1"
    fi
    
    # Go build
    echo "执行 go build..."
    go build -o "$OUTPUT" -ldflags="-s -w" ./main.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} 构建成功: $OUTPUT"
        ls -lh "$OUTPUT"
    else
        echo -e "${RED}✗${NC} 构建失败"
        exit 1
    fi
}

# 清理函数
clean() {
    echo ""
    echo -e "${YELLOW}清理构建文件...${NC}"
    
    rm -f auditmetricbeat
    rm -f auditmetricbeat.exe
    
    echo -e "${GREEN}✓${NC} 清理完成"
}

# 测试函数
test_build() {
    echo ""
    echo -e "${YELLOW}运行测试...${NC}"
    
    go test -v ./...
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} 测试通过"
    else
        echo -e "${RED}✗${NC} 测试失败"
        exit 1
    fi
}

# 安装函数
install() {
    echo ""
    echo -e "${YELLOW}安装 AuditMetricbeat...${NC}"
    
    # 检查权限
    if [ "$EUID" -ne 0 ]; then 
        echo -e "${RED}请使用 root 权限运行安装${NC}"
        exit 1
    fi
    
    # 创建目录
    mkdir -p /usr/share/auditmetricbeat
    mkdir -p /etc/auditmetricbeat
    mkdir -p /var/log/auditmetricbeat
    
    # 复制文件
    cp auditmetricbeat /usr/share/auditmetricbeat/
    cp auditmetricbeat.yml /etc/auditmetricbeat/
    cp -r audit.rules.d /etc/auditmetricbeat/
    cp auditmetricbeat.service /etc/systemd/system/
    
    # 设置权限
    chmod 755 /usr/share/auditmetricbeat/auditmetricbeat
    chmod 644 /etc/auditmetricbeat/auditmetricbeat.yml
    
    # 重新加载 systemd
    systemctl daemon-reload
    
    echo -e "${GREEN}✓${NC} 安装完成"
    echo ""
    echo "启动服务: systemctl start auditmetricbeat"
    echo "启用服务: systemctl enable auditmetricbeat"
    echo "查看状态: systemctl status auditmetricbeat"
}

# 显示帮助
show_help() {
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  build [name]    - 构建二进制文件 (可选: 指定输出文件名)"
    echo "  clean           - 清理构建文件"
    echo "  test            - 运行测试"
    echo "  install         - 安装到系统 (需要 root)"
    echo "  all             - 清理 + 构建 + 测试"
    echo "  help            - 显示此帮助信息"
    echo ""
}

# 主逻辑
case "${1:-build}" in
    build)
        build "$2"
        ;;
    clean)
        clean
        ;;
    test)
        test_build
        ;;
    install)
        build
        install
        ;;
    all)
        clean
        build
        test_build
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo -e "${RED}未知命令: $1${NC}"
        show_help
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  操作完成${NC}"
echo -e "${GREEN}=========================================${NC}"
