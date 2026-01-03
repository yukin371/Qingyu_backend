#!/bin/bash

# 青羽写作平台 - 数据库备份脚本
# 用途：备份MongoDB和Redis数据
# 作者：Claude Code
# 日期：2026-01-02

# ============ 配置 ============
# MongoDB配置
MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_DB="qingyu"
MONGO_USER=""
MONGO_PASSWORD=""

# Redis配置
REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_PASSWORD=""

# 备份配置
BACKUP_DIR="/data/backups/qingyu"
RETENTION_DAYS=30  # 保留最近30天的备份
S3_BUCKET=""       # 如果使用S3，配置bucket名称
S3_REGION=""       # S3区域

# 日志配置
LOG_FILE="/var/log/qingyu/backup.log"
LOCK_FILE="/var/run/qingyu/backup.lock"

# 通知配置（可选）
WEBHOOK_URL=""     # 企业微信/钉钉webhook
EMAIL_TO=""        # 邮件接收者

# ============ 函数定义 ============

# 日志函数
log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_FILE"
}

# 锁文件检查
check_lock() {
    if [ -f "$LOCK_FILE" ]; then
        local pid=$(cat "$LOCK_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            log "ERROR" "备份脚本正在运行 (PID: $pid)"
            exit 1
        else
            log "WARN" "删除过期的锁文件"
            rm -f "$LOCK_FILE"
        fi
    fi
    echo $$ > "$LOCK_FILE"
}

# 清理锁文件
cleanup() {
    rm -f "$LOCK_FILE"
    log "INFO" "清理完成"
}

# 发送通知
send_notification() {
    local status=$1
    local message=$2

    # 企业微信/钉钉通知
    if [ -n "$WEBHOOK_URL" ]; then
        curl -s -X POST "$WEBHOOK_URL" \
            -H "Content-Type: application/json" \
            -d "{\"msgtype\": \"text\", \"text\": {\"content\": \"青羽数据库备份 $status\n$message\"}}" \
            > /dev/null 2>&1
    fi

    # 邮件通知（需要配置mail命令）
    if [ -n "$EMAIL_TO" ]; then
        echo "$message" | mail -s "青羽数据库备份 $status" "$EMAIL_TO"
    fi
}

# MongoDB备份
backup_mongodb() {
    log "INFO" "开始备份MongoDB..."

    local date=$(date '+%Y%m%d_%H%M%S')
    local backup_dir="$BACKUP_DIR/mongodb/$date"
    local backup_file="$BACKUP_DIR/mongodb_backup_$date.gz"

    # 创建备份目录
    mkdir -p "$backup_dir"

    # 构建mongodump命令
    local cmd="mongodump --host=$MONGO_HOST --port=$MONGO_PORT --db=$MONGO_DB --out=$backup_dir"

    if [ -n "$MONGO_USER" ]; then
        cmd="$cmd --username=$MONGO_USER --password=$MONGO_PASSWORD --authenticationDatabase=admin"
    fi

    # 执行备份
    log "INFO" "执行: $cmd"
    if eval $cmd >> "$LOG_FILE" 2>&1; then
        log "INFO" "MongoDB备份成功: $backup_dir"

        # 压缩备份
        log "INFO" "压缩备份文件..."
        if tar -czf "$backup_file" -C "$BACKUP_DIR/mongodb" "$date"; then
            log "INFO" "备份压缩成功: $backup_file"

            # 获取文件大小
            local size=$(du -h "$backup_file" | cut -f1)
            log "INFO" "备份文件大小: $size"

            # 删除未压缩的备份
            rm -rf "$backup_dir"

            # 上传到S3（如果配置）
            if [ -n "$S3_BUCKET" ]; then
                log "INFO" "上传到S3..."
                aws s3 cp "$backup_file" "s3://$S3_BUCKET/backups/mongodb/" \
                    --region "$S3_REGION" >> "$LOG_FILE" 2>&1
            fi

            return 0
        else
            log "ERROR" "备份压缩失败"
            return 1
        fi
    else
        log "ERROR" "MongoDB备份失败"
        return 1
    fi
}

# Redis备份
backup_redis() {
    log "INFO" "开始备份Redis..."

    local date=$(date '+%Y%m%d_%H%M%S')
    local backup_file="$BACKUP_DIR/redis/redis_backup_$date.rdb"

    # 创建备份目录
    mkdir -p "$BACKUP_DIR/redis"

    # 构建redis-cli命令
    local cmd="redis-cli -h $REDIS_HOST -p $REDIS_PORT"

    if [ -n "$REDIS_PASSWORD" ]; then
        cmd="$cmd -a $REDIS_PASSWORD"
    fi

    # 触发BGSAVE
    log "INFO" "触发Redis BGSAVE..."
    if $cmd BGSAVE >> "$LOG_FILE" 2>&1; then
        # 等待BGSAVE完成
        local max_wait=300  # 最多等待5分钟
        local waited=0
        while [ $waited -lt $max_wait ]; do
            local lastsave=$($cmd LASTSAVE)
            local current=$($cmd TIME)
            if [ $lastsave -gt $((current - 10)) ]; then
                break
            fi
            sleep 5
            waited=$((waited + 5))
        done

        # 复制RDB文件
        local rdb_dir="/var/lib/redis"  # Redis数据目录
        if [ -f "$rdb_dir/dump.rdb" ]; then
            cp "$rdb_dir/dump.rdb" "$backup_file"
            log "INFO" "Redis备份成功: $backup_file"

            # 压缩备份
            gzip "$backup_file"
            backup_file="$backup_file.gz"

            # 获取文件大小
            local size=$(du -h "$backup_file" | cut -f1)
            log "INFO" "备份文件大小: $size"

            # 上传到S3（如果配置）
            if [ -n "$S3_BUCKET" ]; then
                log "INFO" "上传到S3..."
                aws s3 cp "$backup_file" "s3://$S3_BUCKET/backups/redis/" \
                    --region "$S3_REGION" >> "$LOG_FILE" 2>&1
            fi

            return 0
        else
            log "ERROR" "找不到Redis RDB文件"
            return 1
        fi
    else
        log "ERROR" "Redis BGSAVE失败"
        return 1
    fi
}

# 清理旧备份
cleanup_old_backups() {
    log "INFO" "清理 $RETENTION_DAYS 天前的备份..."

    # 清理MongoDB备份
    log "INFO" "清理MongoDB旧备份..."
    find "$BACKUP_DIR/mongodb" -name "*.gz" -mtime +$RETENTION_DAYS -delete 2>/dev/null

    # 清理Redis备份
    log "INFO" "清理Redis旧备份..."
    find "$BACKUP_DIR/redis" -name "*.gz" -mtime +$RETENTION_DAYS -delete 2>/dev/null

    # 清理S3旧备份（如果配置）
    if [ -n "$S3_BUCKET" ]; then
        log "INFO" "清理S3旧备份..."
        local cutoff_date=$(date -d "$RETENTION_DAYS days ago" '+%Y%m%d')
        aws s3 ls "s3://$S3_BUCKET/backups/" | while read -r line; do
            local file_date=$(echo "$line" | awk '{print $1}' | tr -d '-')
            local file_name=$(echo "$line" | awk '{print $4}')
            if [ "$file_date" -lt "$cutoff_date" ]; then
                aws s3 rm "s3://$S3_BUCKET/$file_name" --region "$S3_REGION" >> "$LOG_FILE" 2>&1
            fi
        done
    fi

    log "INFO" "旧备份清理完成"
}

# 备份统计
backup_stats() {
    log "INFO" "备份统计信息..."

    local mongo_count=$(find "$BACKUP_DIR/mongodb" -name "*.gz" 2>/dev/null | wc -l)
    local redis_count=$(find "$BACKUP_DIR/redis" -name "*.gz" 2>/dev/null | wc -l)
    local mongo_size=$(du -sh "$BACKUP_DIR/mongodb" 2>/dev/null | cut -f1)
    local redis_size=$(du -sh "$BACKUP_DIR/redis" 2>/dev/null | cut -f1)

    log "INFO" "MongoDB备份数量: $mongo_count"
    log "INFO" "MongoDB备份大小: $mongo_size"
    log "INFO" "Redis备份数量: $redis_count"
    log "INFO" "Redis备份大小: $redis_size"
}

# ============ 主流程 ============

main() {
    # 设置trap确保清理
    trap cleanup EXIT INT TERM

    log "INFO" "========== 开始备份 =========="

    # 检查锁文件
    check_lock

    # 创建备份目录
    mkdir -p "$BACKUP_DIR"/{mongodb,redis}

    # 初始化统计
    local success=0
    local failed=0
    local errors=""

    # MongoDB备份
    if backup_mongodb; then
        success=$((success + 1))
    else
        failed=$((failed + 1))
        errors="$errors\nMongoDB备份失败"
    fi

    # Redis备份
    if backup_redis; then
        success=$((success + 1))
    else
        failed=$((failed + 1))
        errors="$errors\nRedis备份失败"
    fi

    # 清理旧备份
    cleanup_old_backups

    # 统计信息
    backup_stats

    # 总结
    log "INFO" "========== 备份完成 =========="
    log "INFO" "成功: $success, 失败: $failed"

    # 发送通知
    if [ $failed -eq 0 ]; then
        send_notification "成功" "所有备份已完成\n时间: $(date '+%Y-%m-%d %H:%M:%S')"
    else
        send_notification "失败" "部分备份失败\n$errors"
        exit 1
    fi
}

# 执行主流程
main "$@"
