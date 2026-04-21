# AuditMetricbeat

AuditMetricbeat 是一个合并了 Auditbeat 和 Metricbeat 功能的统一数据采集器。它在一个可执行文件中同时提供：

- **系统审计监控** (来自 Auditbeat)
  - Linux 内核审计事件
  - 文件完整性监控
  - 用户行为追踪
  - 系统调用监控

- **系统指标采集** (来自 Metricbeat)
  - CPU、内存、磁盘使用率
  - 网络流量统计
  - 进程监控
  - 系统负载

## 项目结构

```
auditmetricbeat/
├── main.go                      # 入口文件
├── cmd/
│   └── root.go                 # 命令配置
├── include/
│   └── imports.go              # 模块注册
├── audit.rules.d/
│   └── audit.rules             # 审计规则配置
├── auditmetricbeat.yml          # 主配置文件
├── auditmetricbeat.service      # systemd 服务文件
├── magefile.go                  # 构建配置
└── go.mod                       # Go 模块依赖
```

## 快速开始

### 1. 构建

```bash
# 进入项目目录
cd auditmetricbeat

# 使用 mage 构建
mage build

# 或直接使用 go build
go build -o auditmetricbeat ./main.go
```

### 2. 配置

编辑 `auditmetricbeat.yml` 配置文件：

```yaml
# Elasticsearch 输出配置
output.elasticsearch:
  hosts: ["localhost:9200"]
  # username: "elastic"
  # password: "changeme"
```

### 3. 运行

```bash
# 测试运行
sudo ./auditmetricbeat -e -c auditmetricbeat.yml --strict.perms=false

# 后台运行
sudo ./auditmetricbeat -c auditmetricbeat.yml &
```

### 4. 安装为系统服务

```bash
# 复制二进制文件
sudo cp auditmetricbeat /usr/share/auditmetricbeat/

# 复制配置文件
sudo mkdir -p /etc/auditmetricbeat
sudo cp auditmetricbeat.yml /etc/auditmetricbeat/
sudo cp -r audit.rules.d /etc/auditmetricbeat/

# 安装服务
sudo cp auditmetricbeat.service /etc/systemd/system/

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable auditmetricbeat
sudo systemctl start auditmetricbeat

# 查看状态
sudo systemctl status auditmetricbeat
```

## 配置说明

### Auditbeat 模块配置

```yaml
# 审计模块
- module: auditd
  enabled: true
  audit_rule_files:
    - '${path.config}/audit.rules.d/*.rules'
  resolve_ids: true
  failure_mode: silent

# 文件完整性监控
- module: file_integrity
  enabled: true
  paths:
    - /bin
    - /usr/bin
    - /etc
```

### Metricbeat 模块配置

```yaml
# 系统指标
- module: system
  enabled: true
  metricsets:
    - cpu
    - memory
    - network
    - diskio
    - filesystem
    - process
  period: 10s
```

## 索引策略

AuditMetricbeat 会自动将不同类型的数据发送到不同的索引：

- `auditmetricbeat-audit-YYYY.MM.DD` - 审计事件
- `auditmetricbeat-fileintegrity-YYYY.MM.DD` - 文件完整性事件
- `auditmetricbeat-metrics-YYYY.MM.DD` - 系统指标

## 字段说明

### 审计字段
- `auditd.message_type` - 审计消息类型
- `auditd.sequence` - 序列号
- `auditd.result` - 结果 (success/fail)
- `user.*` - 用户信息
- `process.*` - 进程信息
- `file.*` - 文件信息

### 指标字段
- `system.cpu.*` - CPU 指标
- `system.memory.*` - 内存指标
- `system.network.*` - 网络指标
- `system.process.*` - 进程指标

## 审计规则

审计规则位于 `audit.rules.d/` 目录，支持以下配置：

- 系统调用监控 (execve, setuid, etc.)
- 文件访问监控 (/etc, /var/log, etc.)
- 特权操作监控
- 网络配置变更监控
- 用户登录监控

## 性能调优

### 队列配置

```yaml
queue.mem:
  events: 4096
  flush.min_events: 512
  flush.timeout: 10s
```

### 批量配置

```yaml
output.elasticsearch:
  bulk_max_size: 2048
  compression_level: 5
```

## 日志

日志文件位于 `/var/log/auditmetricbeat/auditmetricbeat.log`

```bash
# 查看日志
sudo tail -f /var/log/auditmetricbeat/auditmetricbeat.log
```

## 故障排查

### 检查配置

```bash
# 测试配置
sudo ./auditmetricbeat test config -c auditmetricbeat.yml

# 测试输出
sudo ./auditmetricbeat test output -c auditmetricbeat.yml
```

### 调试模式

```bash
# 启用调试日志
sudo ./auditmetricbeat -e -d "*" -c auditmetricbeat.yml
```

## 许可证

Apache License 2.0

## 贡献

欢迎提交 Issue 和 Pull Request！
