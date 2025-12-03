#!/usr/bin/env zsh

# Etcd 服务管理脚本
# 用于查看和管理服务注册信息
# 注意: 此脚本使用 zsh，macOS 默认支持

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查 etcdctl 是否安装
check_etcdctl() {
    if ! command -v etcdctl &> /dev/null; then
        echo -e "${RED}错误: etcdctl 未安装${NC}"
        echo "安装方法:"
        echo "  macOS: brew install etcd"
        echo "  Linux: sudo apt-get install etcd-client"
        exit 1
    fi
}

# 列出所有已注册的服务
list_services() {
    echo -e "${BLUE}=== 已注册的服务 ===${NC}\n"

    # 获取所有服务
    services=$(etcdctl get "" --prefix --keys-only | grep -E "\.rpc/")

    if [ -z "$services" ]; then
        echo -e "${YELLOW}没有已注册的服务${NC}"
        return
    fi

    # 提取所有唯一的服务名
    service_names=$(echo "$services" | cut -d'/' -f1 | sort -u)

    # 显示每个服务的信息
    for service_name in $(echo "$service_names"); do
        echo -e "${GREEN}服务名称:${NC} $service_name"

        # 获取该服务的所有实例
        instances=$(etcdctl get "$service_name" --prefix --keys-only)
        instance_count=$(echo "$instances" | wc -l | tr -d ' ')

        echo -e "${GREEN}实例数量:${NC} $instance_count"
        echo -e "${GREEN}服务地址:${NC}"

        # 显示每个实例的地址
        i=1
        for key in $(echo "$instances"); do
            address=$(etcdctl get "$key" --print-value-only)
            echo -e "  [$i] $address"
            i=$((i+1))
        done

        echo ""
    done
}

# 查看指定服务的详细信息
show_service() {
    local service_name=$1

    if [ -z "$service_name" ]; then
        echo -e "${RED}错误: 请指定服务名称${NC}"
        echo "用法: $0 show <service_name>"
        echo "示例: $0 show sms.rpc"
        exit 1
    fi

    echo -e "${BLUE}=== $service_name 详细信息 ===${NC}\n"

    # 获取服务信息
    keys=$(etcdctl get "$service_name" --prefix --keys-only)

    if [ -z "$keys" ]; then
        echo -e "${YELLOW}服务 $service_name 未注册${NC}"
        return
    fi

    # 显示每个实例的详细信息
    instance_num=1
    while IFS= read -r key; do
        echo -e "${GREEN}实例 #$instance_num${NC}"
        echo -e "${GREEN}Key:${NC} $key"

        address=$(etcdctl get "$key" --print-value-only)
        echo -e "${GREEN}地址:${NC} $address"

        # 提取 host 和 port
        host=$(echo "$address" | cut -d':' -f1)
        port=$(echo "$address" | cut -d':' -f2)

        echo -e "${GREEN}Host:${NC} $host"
        echo -e "${GREEN}Port:${NC} $port"

        # 检查服务是否在线
        if nc -z "$host" "$port" 2>/dev/null; then
            echo -e "${GREEN}状态:${NC} ✓ 在线"
        else
            echo -e "${RED}状态:${NC} ✗ 离线"
        fi

        echo ""
        instance_num=$((instance_num+1))
    done <<< "$keys"
}

# 监控服务注册变化
watch_services() {
    echo -e "${BLUE}=== 监控服务注册变化 ===${NC}"
    echo -e "${YELLOW}按 Ctrl+C 退出${NC}\n"

    etcdctl watch "" --prefix
}

# 删除指定服务的注册信息
delete_service() {
    local service_name=$1

    if [ -z "$service_name" ]; then
        echo -e "${RED}错误: 请指定服务名称${NC}"
        echo "用法: $0 delete <service_name>"
        echo "示例: $0 delete sms.rpc"
        exit 1
    fi

    # 确认删除
    echo -e "${YELLOW}警告: 即将删除服务 $service_name 的所有注册信息${NC}"
    read -p "确认删除? (y/N): " confirm

    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        echo "已取消"
        exit 0
    fi

    # 删除服务
    deleted=$(etcdctl del "$service_name" --prefix)

    if [ "$deleted" -gt 0 ]; then
        echo -e "${GREEN}成功删除 $deleted 条记录${NC}"
    else
        echo -e "${YELLOW}没有找到服务 $service_name${NC}"
    fi
}

# 清理所有离线的服务实例
cleanup_offline() {
    echo -e "${BLUE}=== 清理离线服务 ===${NC}\n"

    services=$(etcdctl get "" --prefix --keys-only | grep -E "\.rpc/")

    if [ -z "$services" ]; then
        echo -e "${YELLOW}没有已注册的服务${NC}"
        return
    fi

    offline_count=0

    while IFS= read -r key; do
        address=$(etcdctl get "$key" --print-value-only)
        host=$(echo "$address" | cut -d':' -f1)
        port=$(echo "$address" | cut -d':' -f2)

        # 检查服务是否在线
        if ! nc -z "$host" "$port" 2>/dev/null; then
            echo -e "${YELLOW}删除离线服务:${NC} $key ($address)"
            etcdctl del "$key" > /dev/null
            offline_count=$((offline_count+1))
        fi
    done <<< "$services"

    if [ $offline_count -eq 0 ]; then
        echo -e "${GREEN}没有离线的服务${NC}"
    else
        echo -e "\n${GREEN}成功清理 $offline_count 个离线服务${NC}"
    fi
}

# 显示帮助信息
show_help() {
    local script_name="etcd-manager.sh"
    echo "Etcd 服务管理脚本"
    echo ""
    echo "用法:"
    echo "  $script_name list                  列出所有已注册的服务"
    echo "  $script_name show <service>        查看指定服务的详细信息"
    echo "  $script_name watch                 监控服务注册变化"
    echo "  $script_name delete <service>      删除指定服务的注册信息"
    echo "  $script_name cleanup               清理所有离线的服务实例"
    echo "  $script_name help                  显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $script_name list"
    echo "  $script_name show sms.rpc"
    echo "  $script_name watch"
    echo "  $script_name delete sms.rpc"
    echo "  $script_name cleanup"
}

# 主程序
main() {
    check_etcdctl

    case "${1:-list}" in
        list)
            list_services
            ;;
        show)
            show_service "$2"
            ;;
        watch)
            watch_services
            ;;
        delete)
            delete_service "$2"
            ;;
        cleanup)
            cleanup_offline
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            echo -e "${RED}错误: 未知命令 '$1'${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"
