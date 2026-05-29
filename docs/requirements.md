# AssassinIoT 平台功能需求文档

## 1. 项目概述

AssassinIoT 是一个物联网设备管理平台，提供设备接入、数据采集、存储、指令下发等功能。平台采用消息驱动架构，以 MQTT 协议连接设备，通过消息队列解耦数据处理链路，使用时序数据库存储设备数据。

### 1.1 技术栈

| 组件 | 技术选型 | 用途 |
|------|---------|------|
| MQTT Broker | EMQX 5.8 | 设备连接与消息路由 |
| 消息队列 | RabbitMQ 4.2 | 数据流解耦与削峰 |
| 时序数据库 | TimescaleDB 2.27 (PostgreSQL 17) | 设备数据存储 |
| 数据管道 | Bento (WarpStream) 1.17 | 流式数据处理 |
| 后端服务 | Go + Ent ORM | API 与业务逻辑 |
| 监控 | Prometheus 3.5 + Grafana 12.4 | 指标采集与可视化 |

### 1.2 架构概览

```
设备 --MQTT--> EMQX --订阅--> Bento(data-receiver) --AMQP--> RabbitMQ --AMQP--> Bento(data-insert) --SQL--> TimescaleDB
                                                                                          |
                                                                                   Prometheus --> Grafana
```

---

## 2. 已实现功能

### 2.1 设备遥测数据上报

**状态：已实现**

设备通过 MQTT 向 `device/{dev_id}/telemetry` 主题发布 JSON 格式的遥测数据。平台接收后将数据展开为单属性记录存入 TimescaleDB。

**数据格式：**
```json
{
  "ts": 1715700000000,
  "temperature": 25.6,
  "humidity": 68.2,
  "voltage": 3.7,
  "rssi": -45,
  "status": "normal"
}
```

**处理逻辑：**
- 从主题中提取 `device_identifier`
- 将 JSON payload 展开为每属性一条记录（EAV 模型）
- `ts` 字段作为时间戳，缺失时使用接收时间
- 无效消息记录日志后丢弃
- 最终批量写入 TimescaleDB（100 条/批或 50ms 窗口）

### 2.2 数据管道

**状态：已实现**

两条 Bento 管道构成完整的数据处理链路：

| 管道 | 输入 | 输出 | 职责 |
|------|------|------|------|
| device-data-receiver | EMQX (MQTT) | RabbitMQ Exchange `assassin.data_entry` | 消息校验、属性展开 |
| device-data-insert | RabbitMQ Queue `assassin.data_entry.data_raw` | TimescaleDB `device_telemetry` | 批量入库 |

### 2.3 数据库

**状态：已实现**

`device_telemetry` 表（TimescaleDB 超表）：

| 列 | 类型 | 说明 |
|----|------|------|
| time | timestamptz | 数据时间戳 |
| device_identifier | varchar | 设备标识 |
| property_name | varchar | 属性名称 |
| value_type | varchar | 值类型 (number/string/bool/json/null) |
| value_number | double precision | 数值 |
| value_string | varchar | 字符串值 |
| value_bool | boolean | 布尔值 |
| value_json | jsonb | JSON 值 |
| received_at | timestamptz | 平台接收时间 |
| created_at | timestamptz | 数据库写入时间 |

- 主键：`(time, device_identifier, property_name)`
- 分区键：`time`
- 索引：`(device_identifier, property_name, time)`、`(property_name, time)`、`(received_at)`

### 2.4 访问控制

**状态：已实现（基础）**

基于 EMQX 文件 ACL 的访问控制：
- 设备只能发布自己的遥测、事件、状态主题
- 设备只能订阅自己的指令、状态查询主题
- 默认拒绝所有未匹配的发布和订阅

### 2.5 监控基础设施

**状态：已实现（基础）**

- Prometheus 采集 Bento 管道和 RabbitMQ 指标
- Grafana 已部署（仪表盘待配置）

### 2.6 测试工具

**状态：已实现**

- `device_simulator.py`：基于线程的设备模拟器，支持自定义设备数和上报间隔
- `stress_test.py`：基于 asyncio 的压力测试工具，支持高吞吐量测试

---

## 3. 待实现功能

每个模块内的功能按优先级排列（P0 必须实现 / P1 重要 / P2 增强 / P3 远期）。

### 3.1 设备管理

设备从注册到注销的完整生命周期管理。

#### 3.1.1 设备注册表 (P0)

设备需要在平台注册后才能接入。

- **devices 表**：存储设备元信息

| 列 | 类型 | 说明 |
|----|------|------|
| dev_id | varchar | 设备唯一标识（主键） |
| tenant_key | varchar | 所属租户标识 |
| access_key | varchar | 设备接入密钥 |
| firmware_ver | varchar | 固件版本 |
| created_at | timestamptz | 注册时间 |
| last_seen | timestamptz | 最后在线时间 |

- 设备注册 API（创建、查询、删除）
- 设备激活流程
- 设备在线状态跟踪

#### 3.1.2 固件版本跟踪 (P2)

- 记录设备当前固件版本
- 固件升级历史

#### 3.1.3 设备影子 Device Shadow (P3)

- 设备离线时缓存期望状态
- 上线后自动同步差异

#### 3.1.4 固件 OTA 升级 (P3)

- 通过指令通道下发固件更新包
- 支持断点续传
- 升级结果回执

---

### 3.2 数据采集

设备数据的接收、处理与存储。**遥测数据上报已实现**，以下为待扩展的部分。

#### 3.2.1 设备事件处理 (P0)

接收和处理设备主动上报的事件消息（告警、状态变更等）。

- 订阅 `device/{dev_id}/event` 主题，建立事件处理管道
- 事件数据存储（`device_events` 超表）
- 事件分类：

| 类型 | 说明 |
|------|------|
| alarm | 告警事件（温度超阈值等） |
| state_change | 状态变更（上线/下线） |
| error | 设备错误 |
| info | 一般信息 |

**device_events 表结构：**

| 列 | 类型 | 说明 |
|----|------|------|
| time | timestamptz | 事件时间 |
| dev_id | varchar | 设备标识 |
| tenant_key | varchar | 租户标识 |
| event_type | varchar | 事件类型 |
| payload | jsonb | 事件内容 |
| received_at | timestamptz | 接收时间 |

#### 3.2.2 数据管道可靠性增强 (P2)

- 死信队列：格式错误的消息路由到专用死信主题，便于排查
- 写入重试：数据库写入失败自动重试（最多 3 次，指数退避）
- 时间戳异常检测：设备时间与接收时间偏差超过阈值时告警
- QoS 策略调整：遥测 QoS 1、事件 QoS 1、指令 QoS 2

#### 3.2.3 Protobuf 数据格式 (P2)

在 JSON 基础上支持 Protobuf 编码，减少数据传输量。

- 定义 `.proto` 文件描述遥测、事件、状态数据结构
- 通过主题后缀区分格式（如 `device/{dev_id}/telemetry/protobuf`）
- 向下兼容 JSON 格式

#### 3.2.4 数据生命周期管理 (P2)

TimescaleDB 自动化数据管理策略：

- 自动分区：每 7 天一个分区
- 自动压缩：30 天以上的数据启用压缩
- 自动删除：180 天以上的数据自动清理

---

### 3.3 状态管理

设备状态的实时跟踪与查询。

#### 3.3.1 设备状态快照 (P0)

维护设备的最新状态，支持上线/离线检测和健康监控。

- 订阅 `device/{dev_id}/state` 主题，接收设备主动上报的状态
- 状态缓存：`device_state` 表，每个设备一行，更新时覆盖
- 平台主动查询：向 `device/{dev_id}/state/query` 发布查询请求
- 离线检测：基于 EMQX Last Will 机制和心跳超时

**device_state 表结构：**

| 列 | 类型 | 说明 |
|----|------|------|
| dev_id | varchar | 设备标识（主键） |
| tenant_key | varchar | 租户标识 |
| online | boolean | 在线状态 |
| state | jsonb | 完整状态信息 |
| updated_at | timestamptz | 最后更新时间 |

**状态 JSON 格式：**
```json
{
  "dev_id": "sensor-001",
  "ts": 1715700000000,
  "online": true,
  "config": {
    "report_interval": 30,
    "temperature_threshold": 80.0
  },
  "health": {
    "battery": 85,
    "signal_strength": -42,
    "uptime": 86400,
    "firmware_version": "v1.2.3",
    "memory_usage": 65
  },
  "last_telemetry": 1715699990000
}
```

#### 3.3.2 设备配置下发 (P1)

平台下发配置更新到设备。

- 平台向 `device/{dev_id}/state/update` 发布配置
- 设备应用配置后通过 `device/{dev_id}/state` 上报最新状态确认
- 配置变更历史记录

---

### 3.4 指令控制

平台到设备的下行控制能力。

#### 3.4.1 下行指令通道 (P1)

平台向设备下发控制指令，设备执行后回执。

- 平台向 `device/{dev_id}/command` 发布指令
- 设备向 `device/{dev_id}/command/reply` 回复执行结果
- 指令状态追踪（已发送、已送达、已执行、失败、超时）
- 指令超时与重试机制

**指令 JSON 格式：**
```json
{
  "cmd_id": "uuid",
  "cmd_type": "restart|config_update|ota|ping",
  "params": {},
  "ts": 1715700000000
}
```

---

### 3.5 认证授权

设备接入安全与权限控制。

#### 3.5.1 多租户隔离 (P1)

支持多个租户共用平台，实现数据与设备的隔离。

- 租户标识（TenantKey）贯穿所有数据表
- 设备只能访问所属租户的资源
- API 层面按租户过滤数据
- EMQX 层面按租户分组管理设备

#### 3.5.2 TLS 加密与证书认证 (P1)

增强设备连接的安全性。

- MQTT over TLS（端口 8883）
- X.509 客户端证书认证
- 证书 CN → 设备 AccessKey，证书 O/OU → 租户 TenantKey

#### 3.5.3 EMQX HTTP 认证回调 (P1)

替代文件 ACL，实现动态认证。

- EMQX 配置 HTTP 认证后端
- Go 服务提供 `/api/v1/mqtt/auth` 回调端点
- 验证 AccessKey 与 TenantKey 的匹配关系

---

### 3.6 API 服务

对外提供 HTTP REST API。

#### 3.6.1 REST API (P0)

供前端和外部系统调用的 HTTP 接口。

**设备管理：**
- `GET /api/v1/devices` — 设备列表
- `GET /api/v1/devices/{dev_id}` — 设备详情
- `POST /api/v1/devices` — 注册设备
- `DELETE /api/v1/devices/{dev_id}` — 删除设备

**数据查询：**
- `GET /api/v1/devices/{dev_id}/telemetry` — 查询遥测数据（支持时间范围、属性筛选）
- `GET /api/v1/devices/{dev_id}/events` — 查询事件数据
- `GET /api/v1/devices/{dev_id}/state` — 查询设备当前状态

**指令下发：**
- `POST /api/v1/devices/{dev_id}/command` — 向设备下发指令
- `GET /api/v1/devices/{dev_id}/command/{cmd_id}` — 查询指令执行状态

#### 3.6.2 Go 后端服务 (P2)

构建模块化 Go 服务承载 API 和消息处理。

**服务模块：**
- `cmd/server` — 入口，启动 MQTT 客户端 + HTTP 服务
- `internal/mqtt` — MQTT 连接管理与消息路由
- `internal/handler` — 消息处理（遥测、事件、状态、指令）
- `internal/repository` — 数据库操作（基于 Ent ORM）
- `internal/auth` — EMQX HTTP 认证回调
- `internal/api` — REST API 路由

---

### 3.7 监控运维

平台可观测性与运维保障。

#### 3.7.1 Grafana 仪表盘与告警 (P2)

完善可观测性体系（Prometheus + Grafana 基础已部署）。

- 设备数据实时看板（遥测趋势、设备地图）
- 设备在线率统计
- 数据吞吐量监控（MQTT 消息速率、数据库写入速率）
- 告警规则：设备离线、数据中断、管道异常、队列积压

#### 3.7.2 EMQX 集群部署 (P3)

- 多节点分布式桥接，支持更大规模的设备连接
- 节点间消息路由与负载均衡

#### 3.7.3 TimescaleDB 分布式集群 (P3)

- 多节点读写分离，提升存储和查询能力

---

## 4. MQTT 主题设计

### 4.1 完整主题树

| 主题 | 方向 | QoS | 用途 | 状态 |
|------|------|-----|------|------|
| `device/{dev_id}/telemetry` | 设备 → 平台 | 1 | 遥测数据上报 | 已实现 |
| `device/{dev_id}/event` | 设备 → 平台 | 1 | 事件/告警上报 | 待实现 |
| `device/{dev_id}/state` | 设备 → 平台 | 1 | 设备状态上报 | 待实现 |
| `device/{dev_id}/command` | 平台 → 设备 | 2 | 下行控制指令 | 待实现 |
| `device/{dev_id}/command/reply` | 设备 → 平台 | 1 | 指令执行回执 | 待实现 |
| `device/{dev_id}/state/query` | 平台 → 设备 | 1 | 查询设备状态 | 待实现 |
| `device/{dev_id}/state/update` | 平台 → 设备 | 1 | 配置更新下发 | 待实现 |

### 4.2 通配符订阅规则

- 平台订阅 `device/+/telemetry` 接收所有设备遥测数据
- 平台订阅 `device/+/event` 接收所有设备事件
- 平台订阅 `device/+/state` 接收所有设备状态
- 设备订阅 `device/{dev_id}/command`、`device/{dev_id}/state/query`、`device/{dev_id}/state/update`

---

## 5. 非功能需求

### 5.1 性能

| 指标 | 目标值 |
|------|--------|
| 设备连接数 | 单节点支持 10,000+ 并发连接 |
| 消息吞吐量 | 1,000+ msg/s 稳定处理 |
| 数据写入延迟 | P99 < 2 秒 |
| 数据查询延迟 | P99 < 500ms（7 天内数据） |

### 5.2 安全

- MQTT 连接启用 TLS 1.3 加密
- 设备使用 X.509 证书认证
- API 接口使用 JWT 认证
- 设备间数据隔离（基于主题 ACL）
- 租户间数据隔离（基于 TenantKey）
- 敏感信息（密钥、密码）不落盘

### 5.3 可靠性

- 消息队列持久化，防止数据丢失
- 数据库写入失败自动重试
- 服务异常自动重启（Docker restart policy）
- 数据定期备份
- 关键主题使用 QoS 1/2 保证消息送达

### 5.4 可维护性

- Docker Compose 一键部署
- 结构化日志（JSON 格式）
- 健康检查端点
- 配置与代码分离（环境变量注入）

---

## 6. 版本记录

| 版本 | 日期 | 说明 |
|------|------|------|
| v0.1 | 2026-05-21 | 初稿，梳理已实现功能与待实现需求 |
