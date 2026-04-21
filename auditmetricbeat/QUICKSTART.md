# AuditMetricbeat 快速开始指南

## 项目概述

AuditMetricbeat 成功将 Auditbeat 和 Metricbeat 合并为一个统一的可执行文件，实现了：

✅ **单一进程** - 同时运行审计和指标采集  
✅ **统一配置** - 一个配置文件管理所有功能  
✅ **资源优化** - 共享运行时，减少内存占用  
✅ **简化管理** - 一个服务，统一监控  

## 已创建的文件

```
auditmetricbeat/
├── main.go                          # 主入口文件
├── cmd/
│   └── root.go                     # 命令配置和初始化
├── include/
│   └── imports.go                  # 模块注册（auditd, file_integrity, system.*）
├── audit.rules.d/
│   └── audit.rules                 # Linux审计规则配置
├── auditmetricbeat.yml              # 主配置文件
├── auditmetricbeat.service          # systemd服务文件
├── magefile.go                      # Mage构建配置
├── go.mod                           # Go模块依赖
├── build.sh                         # Linux/Mac构建脚本
├── build.bat                        # Windows构建脚本
└── README.md                        # 项目文档
```

## 快速构建

### Windows 环境

```powershell
# 进入项目目录
cd f:\Code\beats-my\beats-self\auditmetricbeat

# 运行构建脚本
.\build.bat

# 或使用 go 命令
go build -o auditmetricbeat.exe .\main.go
```

### Linux 环境

```bash
# 进入项目目录
cd /path/to/beats/auditmetricbeat

# 运行构建脚本
chmod +x build.sh
./build.sh build

# 或使用 go 命令
go build -o auditmetricbeat ./main.go
```

### 使用 Mage 构建

```bash
# 安装 mage
go install github.com/magefile/mage@latest

# 构建
mage build

# 清理
mage clean

# 测试
mage test
```

## 测试运行

### Windows

```powershell
# 测试配置（需要管理员权限）
.\auditmetricbeat.exe test config -c auditmetricbeat.yml

# 运行（调试模式）
.\auditmetricbeat.exe -e -c auditmetricbeat.yml --strict.perms=false
```

### Linux

```bash
# 测试配置
sudo ./auditmetricbeat test config -c auditmetricbeat.yml

# 运行（调试模式）
sudo ./auditmetricbeat -e -c auditmetricbeat.yml --strict.perms=false
```

## 配置说明

### 1. 修改 Elasticsearch 连接

编辑 `auditmetricbeat.yml`：

```yaml
output.elasticsearch:
  hosts: ["your-elasticsearch-host:9200"]
  username: "elastic"
  password: "your-password"
```

### 2. 配置审计规则

编辑 `audit.rules.d/audit.rules`，添加你需要的审计规则。

### 3. 启用/禁用模块

在 `auditmetricbeat.yml` 中：

```yaml
# 禁用文件完整性监控
- module: file_integrity
  enabled: false
  
# 启用系统指标
- module: system
  enabled: true
  period: 10s
```

## 安装为系统服务（Linux）

```bash
# 1. 构建
./build.sh build

# 2. 安装（需要root权限）
sudo ./build.sh install

# 3. 启动服务
sudo systemctl start auditmetricbeat

# 4. 设置开机启动
sudo systemctl enable auditmetricbeat

# 5. 查看状态
sudo systemctl status auditmetricbeat

# 6. 查看日志
sudo journalctl -u auditmetricbeat -f
```

## 数据采集验证

### 1. 检查审计事件

```bash
# 运行后，在Kibana中查询
GET /auditmetricbeat-audit-*/_search
{
  "query": {
    "match_all": {}
  },
  "size": 10
}
```

### 2. 检查系统指标

```bash
# 在Kibana中查询
GET /auditmetricbeat-metrics-*/_search
{
  "query": {
    "match": {
      "event.module": "system"
    }
  },
  "size": 10
}
```

## 包含的模块

### Auditbeat 模块

✅ **auditd** - Linux内核审计
- 系统调用监控
- 用户行为追踪
- 文件访问审计

✅ **file_integrity** - 文件完整性监控
- 文件变化检测
- 哈希计算
- 权限监控

### Metricbeat 模块

✅ **system.cpu** - CPU使用率  
✅ **system.memory** - 内存使用  
✅ **system.network** - 网络流量  
✅ **system.diskio** - 磁盘I/O  
✅ **system.filesystem** - 文件系统  
✅ **system.load** - 系统负载  
✅ **system.process** - 进程信息  
✅ **system.process_summary** - 进程摘要  
✅ **system.uptime** - 运行时间  
✅ **system.core** - CPU核心信息  

## 性能优化建议

### 1. 调整采集频率

```yaml
# 审计事件 - 实时处理
- module: auditd
  # 无需配置period

# 系统指标 - 10秒一次
- module: system
  period: 10s
```

### 2. 优化队列

```yaml
queue.mem:
  events: 4096          # 队列大小
  flush.min_events: 512 # 最小刷新事件数
  flush.timeout: 10s    # 超时时间
```

### 3. 批量发送

```yaml
output.elasticsearch:
  bulk_max_size: 2048      # 批量大小
  compression_level: 5     # 压缩级别 (0-9)
```

## 故障排查

### 问题1: 构建失败

```bash
# 检查Go版本
go version

# 需要 Go 1.21+
# 更新依赖
go mod tidy
```

### 问题2: 权限错误

```bash
# Auditbeat需要root权限
sudo ./auditmetricbeat -e -c auditmetricbeat.yml

# 或使用 --strict.perms=false
./auditmetricbeat -e -c auditmetricbeat.yml --strict.perms=false
```

### 问题3: 审计规则加载失败

```bash
# 检查审计规则文件
sudo auditctl -l

# 手动加载规则
sudo auditctl -R audit.rules.d/audit.rules
```

### 问题4: 无法连接Elasticsearch

```bash
# 测试连接
curl -u elastic:password http://localhost:9200

# 检查配置
./auditmetricbeat test output -c auditmetricbeat.yml
```

## 下一步

1. **自定义审计规则** - 根据安全需求调整 `audit.rules.d/audit.rules`
2. **配置告警** - 在Kibana中设置Watchers
3. **创建仪表板** - 可视化审计和指标数据
4. **性能调优** - 根据系统负载调整采集频率
5. **扩展模块** - 添加更多Metricbeat模块（如docker, nginx等）

## 技术支持

- 项目文档: README.md
- 配置示例: auditmetricbeat.yml
- 审计规则: audit.rules.d/audit.rules
- 日志位置: /var/log/auditmetricbeat/auditmetricbeat.log

## 总结

你现在拥有一个功能完整的合并版Beat！它同时具备：

🔒 **安全审计能力** (来自Auditbeat)  
📊 **系统监控能力** (来自Metricbeat)  
⚡ **高性能** - 单一进程，资源优化  
🎯 **易管理** - 统一配置，统一部署  

开始使用吧！🚀
