# 青羽写作平台 - 数据库备份脚本 (PowerShell版本)
# 用途：备份MongoDB和Redis数据
# 作者：Claude Code
# 日期：2026-01-02

# ============ 配置 ============

# MongoDB配置
$MongoHost = "localhost"
$MongoPort = "27017"
$MongoDb = "qingyu"
$MongoUser = ""
$MongoPassword = ""

# Redis配置
$RedisHost = "localhost"
$RedisPort = "6379"
$RedisPassword = ""

# 备份配置
$BackupDir = "D:\backups\qingyu"
$RetentionDays = 30  # 保留最近30天的备份

# 日志配置
$LogFile = "D:\logs\qingyu\backup.log"

# ============ 函数定义 ============

function Write-Log {
    param(
        [string]$Level,
        [string]$Message
    )

    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logMessage = "[$timestamp] [$Level] $Message"
    Write-Host $logMessage
    Add-Content -Path $LogFile -Value $logMessage
}

function Backup-MongoDB {
    Write-Log "INFO" "开始备份MongoDB..."

    $date = Get-Date -Format "yyyyMMdd_HHmmss"
    $backupDir = Join-Path $BackupDir "mongodb\$date"
    $backupFile = Join-Path $BackupDir "mongodb_backup_$date.zip"

    try {
        # 创建备份目录
        New-Item -ItemType Directory -Path $backupDir -Force | Out-Null

        # 构建mongodump命令
        $mongodumpCmd = "mongodump --host=$MongoHost --port=$MongoPort --db=$MongoDb --out=$backupDir"

        if ($MongoUser -ne "") {
            $mongodumpCmd += " --username=$MongoUser --password=$MongoPassword --authenticationDatabase=admin"
        }

        # 执行备份
        Write-Log "INFO" "执行: $mongodumpCmd"
        $result = Invoke-Expression $mongodumpCmd 2>&1

        if ($LASTEXITCODE -eq 0) {
            Write-Log "INFO" "MongoDB备份成功: $backupDir"

            # 压缩备份
            Write-Log "INFO" "压缩备份文件..."
            Compress-Archive -Path "$backupDir\*" -DestinationPath $backupFile -Force

            # 获取文件大小
            $size = (Get-Item $backupFile).Length / 1MB
            Write-Log "INFO" "备份文件大小: $([math]::Round($size, 2)) MB"

            # 删除未压缩的备份
            Remove-Item -Path $backupDir -Recurse -Force

            return $true
        } else {
            Write-Log "ERROR" "MongoDB备份失败: $result"
            return $false
        }
    } catch {
        Write-Log "ERROR" "MongoDB备份异常: $_"
        return $false
    }
}

function Backup-Redis {
    Write-Log "INFO" "开始备份Redis..."

    $date = Get-Date -Format "yyyyMMdd_HHmmss"
    $backupFile = Join-Path $BackupDir "redis\redis_backup_$date.rdb"

    try {
        # 创建备份目录
        New-Item -ItemType Directory -Path (Split-Path $backupFile) -Force | Out-Null

        # 构建redis-cli命令
        $redisCliCmd = "redis-cli -h $RedisHost -p $RedisPort"

        if ($RedisPassword -ne "") {
            $redisCliCmd += " -a $RedisPassword"
        }

        # 触发BGSAVE
        Write-Log "INFO" "触发Redis BGSAVE..."
        $result = Invoke-Expression "$redisCliCmd BGSAVE" 2>&1

        if ($LASTEXITCODE -eq 0) {
            # 等待BGSAVE完成
            $maxWait = 300  # 最多等待5分钟
            $waited = 0
            while ($waited -lt $maxWait) {
                $lastsave = Invoke-Expression "$redisCliCmd LASTSAVE" 2>&1
                Start-Sleep -Seconds 5
                $waited += 5
            }

            # 复制RDB文件
            $rdbDir = "C:\Program Files\Redis\data"  # Redis数据目录
            $rdbFile = Join-Path $rdbDir "dump.rdb"

            if (Test-Path $rdbFile) {
                Copy-Item -Path $rdbFile -Destination $backupFile -Force
                Write-Log "INFO" "Redis备份成功: $backupFile"

                # 获取文件大小
                $size = (Get-Item $backupFile).Length / 1MB
                Write-Log "INFO" "备份文件大小: $([math]::Round($size, 2)) MB"

                return $true
            } else {
                Write-Log "ERROR" "找不到Redis RDB文件: $rdbFile"
                return $false
            }
        } else {
            Write-Log "ERROR" "Redis BGSAVE失败: $result"
            return $false
        }
    } catch {
        Write-Log "ERROR" "Redis备份异常: $_"
        return $false
    }
}

function Remove-OldBackups {
    Write-Log "INFO" "清理 $RetentionDays 天前的备份..."

    $cutoffDate = (Get-Date).AddDays(-$RetentionDays)

    # 清理MongoDB备份
    Write-Log "INFO" "清理MongoDB旧备份..."
    Get-ChildItem -Path "$BackupDir\mongodb\*.zip" | Where-Object {
        $_.LastWriteTime -lt $cutoffDate
    } | Remove-Item -Force

    # 清理Redis备份
    Write-Log "INFO" "清理Redis旧备份..."
    Get-ChildItem -Path "$BackupDir\redis\*.rdb" | Where-Object {
        $_.LastWriteTime -lt $cutoffDate
    } | Remove-Item -Force

    Write-Log "INFO" "旧备份清理完成"
}

function Get-BackupStats {
    Write-Log "INFO" "备份统计信息..."

    $mongoBackups = Get-ChildItem -Path "$BackupDir\mongodb\*.zip" -ErrorAction SilentlyContinue
    $redisBackups = Get-ChildItem -Path "$BackupDir\redis\*.rdb" -ErrorAction SilentlyContinue

    $mongoCount = ($mongoBackups | Measure-Object).Count
    $redisCount = ($redisBackups | Measure-Object).Count
    $mongoSize = ($mongoBackups | Measure-Object -Property Length -Sum).Sum / 1MB
    $redisSize = ($redisBackups | Measure-Object -Property Length -Sum).Sum / 1MB

    Write-Log "INFO" "MongoDB备份数量: $mongoCount"
    Write-Log "INFO" "MongoDB备份大小: $([math]::Round($mongoSize, 2)) MB"
    Write-Log "INFO" "Redis备份数量: $redisCount"
    Write-Log "INFO" "Redis备份大小: $([math]::Round($redisSize, 2)) MB"
}

# ============ 主流程 ============

try {
    Write-Log "INFO" "========== 开始备份 =========="

    # 创建备份目录
    New-Item -ItemType Directory -Path "$BackupDir\mongodb" -Force | Out-Null
    New-Item -ItemType Directory -Path "$BackupDir\redis" -Force | Out-Null
    New-Item -ItemType Directory -Path (Split-Path $LogFile) -Force | Out-Null

    # 初始化统计
    $success = 0
    $failed = 0
    $errors = @()

    # MongoDB备份
    if (Backup-MongoDB) {
        $success++
    } else {
        $failed++
        $errors += "MongoDB备份失败"
    }

    # Redis备份
    if (Backup-Redis) {
        $success++
    } else {
        $failed++
        $errors += "Redis备份失败"
    }

    # 清理旧备份
    Remove-OldBackups

    # 统计信息
    Get-BackupStats

    # 总结
    Write-Log "INFO" "========== 备份完成 =========="
    Write-Log "INFO" "成功: $success, 失败: $failed"

    if ($failed -gt 0) {
        Write-Log "ERROR" "备份失败:"
        $errors | ForEach-Object { Write-Log "ERROR" $_ }
        exit 1
    }

    Write-Log "INFO" "所有备份已完成"
    exit 0

} catch {
    Write-Log "ERROR" "备份脚本异常: $_"
    exit 1
}
