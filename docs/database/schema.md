# 数据库设计

## 1. 租户表 `tenants`

多租户架构的核心，所有设备和数据按租户隔离。

```sql
CREATE TABLE tenants (
    id              BIGINT          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_key      VARCHAR(64)     NOT NULL UNIQUE,
    name            VARCHAR(128)    NOT NULL,
    description     TEXT            DEFAULT '',
    status          VARCHAR(16)     NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- status 约束
ALTER TABLE tenants ADD CONSTRAINT chk_tenants_status
    CHECK (status IN ('active', 'disabled'));

-- 按 tenant_key 查询（最频繁的查询路径）
CREATE INDEX idx_tenants_tenant_key ON tenants (tenant_key);

COMMENT ON TABLE tenants IS '租户表，平台多租户隔离的核心';
COMMENT ON COLUMN tenants.tenant_key IS '租户业务标识，用于 API、MQTT 认证等场景';
COMMENT ON COLUMN tenants.name IS '租户名称，用于前端展示';
COMMENT ON COLUMN tenants.status IS '租户状态：active=正常, disabled=已禁用';
```

**字段说明：**

| 列 | 类型 | 说明 |
|----|------|------|
| id | bigint (自增) | 内部主键，不对外暴露 |
| tenant_key | varchar(64) UNIQUE | 租户业务标识，贯穿整个系统 |
| name | varchar(128) | 租户名称 |
| description | text | 描述信息 |
| status | varchar(16) | active / disabled |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 最后更新时间 |

---

## 2. 物模型分类表 `model_categories`

物模型分类，用于在新增物模型或筛选时以下拉列表形式展示可选类别。每个租户独立管理自己的分类。

```sql
CREATE TABLE model_categories (
    id              BIGINT          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_key      VARCHAR(64)     NOT NULL,
    category_key    VARCHAR(64)     NOT NULL,
    display_name    VARCHAR(128)    NOT NULL,
    description     TEXT            DEFAULT '',
    icon            VARCHAR(64)     DEFAULT '',
    sort_order      INT             NOT NULL DEFAULT 0,
    status          VARCHAR(16)     NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),

    UNIQUE (tenant_key, category_key)
);

-- 外键约束
ALTER TABLE model_categories ADD CONSTRAINT fk_model_categories_tenant
    FOREIGN KEY (tenant_key) REFERENCES tenants (tenant_key)
    ON DELETE RESTRICT ON UPDATE CASCADE;

-- 状态约束
ALTER TABLE model_categories ADD CONSTRAINT chk_model_categories_status
    CHECK (status IN ('active', 'disabled'));

CREATE INDEX idx_model_categories_tenant ON model_categories (tenant_key);
CREATE INDEX idx_model_categories_sort ON model_categories (tenant_key, sort_order);

COMMENT ON TABLE model_categories IS '物模型分类表，用于前端下拉选择和筛选';
COMMENT ON COLUMN model_categories.category_key IS '分类标识，如 temp_humidity_sensor、smart_lock';
COMMENT ON COLUMN model_categories.display_name IS '分类展示名称，如"温湿度传感器"、"智能门锁"';
COMMENT ON COLUMN model_categories.sort_order IS '排序值，数字越小越靠前';
COMMENT ON COLUMN model_categories.status IS 'active=启用, disabled=禁用';
```

**字段说明：**

| 列 | 类型 | 说明 |
|----|------|------|
| id | bigint (自增) | 内部主键 |
| tenant_key | varchar(64) FK | 所属租户，引用 `tenants.tenant_key` |
| category_key | varchar(64) | 分类业务标识，如 `temp_humidity_sensor` |
| display_name | varchar(128) | 分类展示名称，用于前端下拉列表 |
| description | text | 描述信息 |
| icon | varchar(64) | 图标标识，选填 |
| sort_order | int | 排序值，数字越小越靠前，默认 0 |
| status | varchar(16) | active / disabled |
| created_at | timestamptz | 创建时间 |
| updated_at | timestamptz | 最后更新时间 |

---

## 3. 物模型表 `thing_models`

物模型（Thing Model）定义一类设备的能力：**属性**（上报什么数据）、**服务**（可执行什么指令）、**事件**（可上报什么事件）。设备实例化时必须绑定一个物模型，上报数据时按模型定义校验。

```sql
CREATE TABLE thing_models (
    id              BIGINT          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    model_key       VARCHAR(64)     NOT NULL UNIQUE,
    tenant_key      VARCHAR(64)     NOT NULL,
    name            VARCHAR(128)    NOT NULL,
    description     TEXT            DEFAULT '',
    category        VARCHAR(64)     NOT NULL,
    properties      JSONB           NOT NULL DEFAULT '{}',
    services        JSONB           NOT NULL DEFAULT '{}',
    events          JSONB           NOT NULL DEFAULT '{}',
    version         VARCHAR(16)     NOT NULL DEFAULT '1.0',
    status          VARCHAR(16)     NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- 外键约束
ALTER TABLE thing_models ADD CONSTRAINT fk_thing_models_tenant
    FOREIGN KEY (tenant_key) REFERENCES tenants (tenant_key)
    ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE thing_models ADD CONSTRAINT fk_thing_models_category
    FOREIGN KEY (tenant_key, category) REFERENCES model_categories (tenant_key, category_key)
    ON DELETE RESTRICT ON UPDATE CASCADE;

-- 状态约束
ALTER TABLE thing_models ADD CONSTRAINT chk_thing_models_status
    CHECK (status IN ('active', 'deprecated', 'disabled'));

CREATE INDEX idx_thing_models_tenant ON thing_models (tenant_key);
CREATE INDEX idx_thing_models_category ON thing_models (tenant_key, category);

COMMENT ON TABLE thing_models IS '物模型表，定义一类设备的属性/服务/事件规范';
COMMENT ON COLUMN thing_models.model_key IS '模型标识，如 temp_sensor_v1、smart_lock_v2';
COMMENT ON COLUMN thing_models.category IS '设备品类，引用 model_categories.category_key';
COMMENT ON COLUMN thing_models.properties IS '属性定义（JSON Schema），描述设备上报的遥测数据点';
COMMENT ON COLUMN thing_models.services IS '服务定义，描述可下发给设备的指令及输入输出参数';
COMMENT ON COLUMN thing_models.events IS '事件定义，描述设备可上报的事件类型及负载结构';
COMMENT ON COLUMN thing_models.version IS '模型版本号，支持平滑升级';
COMMENT ON COLUMN thing_models.status IS 'active=启用, deprecated=已废弃但已有设备仍可用, disabled=禁用';
```

### 2.1 properties 结构

定义设备上报的遥测属性，每个属性包含类型、单位、取值范围等元信息：

```json
{
  "temperature": {
    "name": "温度",
    "type": "float",
    "unit": "°C",
    "range": { "min": -40, "max": 85 },
    "step": 0.1,
    "required": true
  },
  "humidity": {
    "name": "湿度",
    "type": "int",
    "unit": "%",
    "range": { "min": 0, "max": 100 },
    "required": true
  },
  "battery": {
    "name": "电池电量",
    "type": "int",
    "unit": "%",
    "range": { "min": 0, "max": 100 }
  },
  "rssi": {
    "name": "信号强度",
    "type": "int",
    "unit": "dBm",
    "range": { "min": -100, "max": 0 }
  },
  "status": {
    "name": "工作状态",
    "type": "string",
    "enum": ["normal", "warning", "error"]
  }
}
```

字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| name | string | 属性显示名称 |
| type | string | 数据类型：int / float / string / bool / json |
| unit | string | 单位，选填 |
| range | object | 数值范围 { min, max }，选填 |
| step | number | 精度/步长，选填 |
| enum | array | 枚举值列表（仅 string 类型），选填 |
| required | bool | 是否必须上报，默认 false |

### 2.2 services 结构

定义可下发给设备的指令：

```json
{
  "restart": {
    "name": "重启设备",
    "call_type": "async",
    "input": {
      "delay": { "type": "int", "description": "延迟重启秒数", "default": 0 }
    },
    "output": {
      "result": { "type": "string", "description": "执行结果", "enum": ["ok", "fail"] }
    }
  },
  "set_report_interval": {
    "name": "设置上报间隔",
    "call_type": "sync",
    "input": {
      "interval": { "type": "int", "description": "上报间隔（秒）", "range": { "min": 5, "max": 3600 } }
    },
    "output": {
      "accepted": { "type": "bool", "description": "是否接受" }
    }
  }
}
```

字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| name | string | 服务显示名称 |
| call_type | string | sync=同步（设备须立即返回结果）, async=异步（设备先确认收到，执行完再上报结果） |
| input | object | 输入参数定义，key 为参数名，value 含 type、description、default、range、enum |
| output | object | 输出参数定义，同上结构 |

### 2.3 events 结构

定义设备可上报的事件：

```json
{
  "over_temp_alarm": {
    "name": "超温告警",
    "type": "alarm",
    "level": "warning",
    "payload": {
      "temperature": { "type": "float", "description": "当前温度" },
      "threshold": { "type": "float", "description": "告警阈值" }
    }
  },
  "low_battery": {
    "name": "低电量告警",
    "type": "alarm",
    "level": "warning",
    "payload": {
      "battery": { "type": "int", "description": "当前电量" }
    }
  },
  "boot": {
    "name": "设备启动",
    "type": "info",
    "level": "info",
    "payload": {
      "firmware_ver": { "type": "string", "description": "固件版本" },
      "boot_reason": { "type": "string", "description": "启动原因", "enum": ["power_on", "reset", "watchdog", "ota"] }
    }
  }
}
```

字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| name | string | 事件显示名称 |
| type | string | 事件分类：alarm / state_change / error / info |
| level | string | 严重级别：info / warning / error / critical |
| payload | object | 事件负载定义，key 为字段名，value 含 type、description、enum |

---

## 3. 设备表 `devices`

设备是物模型的**实例**，注册时必须指定 `model_key`。

```sql
CREATE TABLE devices (
    dev_id          VARCHAR(64)     PRIMARY KEY,
    tenant_key      VARCHAR(64)     NOT NULL,
    model_key       VARCHAR(64)     NOT NULL,
    device_name     VARCHAR(128)    NOT NULL DEFAULT '',
    access_key      VARCHAR(128)    NOT NULL,
    firmware_ver    VARCHAR(32)     DEFAULT '',
    properties_cfg  JSONB           DEFAULT '{}',
    status          VARCHAR(16)     NOT NULL DEFAULT 'active',
    online          BOOLEAN         NOT NULL DEFAULT FALSE,
    metadata        JSONB           DEFAULT '{}',
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    last_seen       TIMESTAMPTZ
);

-- 外键约束
ALTER TABLE devices ADD CONSTRAINT fk_devices_tenant
    FOREIGN KEY (tenant_key) REFERENCES tenants (tenant_key)
    ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE devices ADD CONSTRAINT fk_devices_model
    FOREIGN KEY (model_key) REFERENCES thing_models (model_key)
    ON DELETE RESTRICT ON UPDATE CASCADE;

-- 状态约束
ALTER TABLE devices ADD CONSTRAINT chk_devices_status
    CHECK (status IN ('active', 'inactive', 'disabled'));

-- 常用查询索引
CREATE INDEX idx_devices_tenant_key ON devices (tenant_key);
CREATE INDEX idx_devices_model_key ON devices (tenant_key, model_key);
CREATE INDEX idx_devices_status ON devices (status);
CREATE INDEX idx_devices_tenant_status ON devices (tenant_key, status);
CREATE INDEX idx_devices_last_seen ON devices (last_seen DESC);

CREATE UNIQUE INDEX idx_devices_access_key ON devices (access_key);

COMMENT ON TABLE devices IS '设备注册表，每个设备是某个物模型的实例';
COMMENT ON COLUMN devices.dev_id IS '设备唯一标识，与 MQTT 客户端用户名一致';
COMMENT ON COLUMN devices.model_key IS '绑定的物模型标识，决定设备的数据结构和能力';
COMMENT ON COLUMN devices.access_key IS '设备接入密钥，用于 MQTT 认证';
COMMENT ON COLUMN devices.properties_cfg IS '属性个性化配置，可覆盖物模型中属性的默认值（如上报间隔）';
COMMENT ON COLUMN devices.status IS '设备状态：active=正常, inactive=未激活, disabled=已禁用';
COMMENT ON COLUMN devices.online IS '是否在线（运行时更新）';
COMMENT ON COLUMN devices.metadata IS '扩展元数据，如安装位置、厂商信息等';
COMMENT ON COLUMN devices.last_seen IS '最后在线时间';
```

**字段说明：**

| 列 | 类型 | 说明 |
|----|------|------|
| dev_id | varchar(64) PK | 设备唯一标识，同时也是 MQTT 客户端用户名 |
| tenant_key | varchar(64) FK | 所属租户，引用 `tenants.tenant_key` |
| model_key | varchar(64) FK | 物模型标识，引用 `thing_models.model_key` |
| device_name | varchar(128) | 设备名称，用于前端展示 |
| access_key | varchar(128) UNIQUE | 设备接入密钥，用于 MQTT 连接认证 |
| firmware_ver | varchar(32) | 当前固件版本 |
| properties_cfg | jsonb | 属性个性化配置（如自定义上报间隔、阈值），可覆盖物模型默认值 |
| status | varchar(16) | active=正常, inactive=未激活, disabled=已禁用 |
| online | boolean | 运行时在线状态，由状态管理模块更新 |
| metadata | jsonb | 扩展元数据，如安装位置、厂商信息等 |
| created_at | timestamptz | 注册时间 |
| updated_at | timestamptz | 最后更新时间 |
| last_seen | timestamptz | 最后在线时间 |

---

## 4. ER 关系

```
tenants (1) ────< (N) model_categories ────< (N) thing_models ────< (N) devices
    │                   │                          │                      │
    │ tenant_key ───────┘ (FK)                     │                      │
    │                   category_key ──────────────┘ (FK)                 │
    │                                              model_key ────────────┘ (FK)
    │
    └─── device_telemetry (tenant_key)
    └─── device_events     (tenant_key)
    └─── device_state      (tenant_key)
```

- 一个租户可以定义多个分类
- 一个分类可以有多个物模型
- 一个物模型可以有多个设备实例
- 设备通过 `model_key` 绑定物模型，继承其属性/服务/事件定义
- `tenant_key` 贯穿所有表，实现租户数据隔离

---

## 5. 典型查询

### 5.1 创建物模型

```sql
INSERT INTO thing_models (model_key, tenant_key, name, category, properties, services, events)
VALUES (
    'temp_sensor_v1',
    'tenant_01',
    '温湿度传感器 v1',
    'environment_sensor',
    '{
        "temperature": {"name": "温度", "type": "float", "unit": "°C", "range": {"min": -40, "max": 85}, "step": 0.1, "required": true},
        "humidity": {"name": "湿度", "type": "int", "unit": "%", "range": {"min": 0, "max": 100}, "step": 1, "required": true},
        "battery": {"name": "电池电量", "type": "int", "unit": "%", "range": {"min": 0, "max": 100}},
        "rssi": {"name": "信号强度", "type": "int", "unit": "dBm"}
    }',
    '{
        "restart": {"name": "重启", "call_type": "async", "input": {}, "output": {"result": {"type": "string"}}},
        "set_report_interval": {"name": "设置上报间隔", "call_type": "sync", "input": {"interval": {"type": "int", "range": {"min": 5, "max": 3600}}}, "output": {"accepted": {"type": "bool"}}}
    }',
    '{
        "over_temp_alarm": {"name": "超温告警", "type": "alarm", "level": "warning", "payload": {"temperature": {"type": "float"}, "threshold": {"type": "float"}}},
        "low_battery": {"name": "低电量", "type": "alarm", "level": "warning", "payload": {"battery": {"type": "int"}}}
    }'
);
```

### 5.2 设备注册

```sql
-- 创建租户
INSERT INTO tenants (tenant_key, name)
VALUES ('tenant_01', '演示租户');

-- 注册设备，绑定物模型
INSERT INTO devices (dev_id, tenant_key, model_key, device_name, access_key, properties_cfg)
VALUES (
    'sensor-001',
    'tenant_01',
    'temp_sensor_v1',
    '温湿度传感器A',
    'ak_abc123def456',
    '{"temperature": {"report_interval": 30}, "humidity": {"report_interval": 60}}'
);
```

### 5.3 设备鉴权

EMQX HTTP 认证回调时，验证设备凭据：

```sql
-- 验证 access_key 有效且设备状态正常
SELECT d.dev_id, d.tenant_key, t.status AS tenant_status
FROM devices d
JOIN tenants t ON d.tenant_key = t.tenant_key
WHERE d.access_key = $1
  AND d.status = 'active'
  AND t.status = 'active';
```

### 5.4 设备列表（关联物模型）

```sql
SELECT d.dev_id, d.device_name, d.online, d.firmware_ver, d.last_seen,
       tm.name AS model_name, tm.category, tm.version AS model_version
FROM devices d
JOIN thing_models tm ON d.model_key = tm.model_key
WHERE d.tenant_key = $1
  AND d.status = 'active'
ORDER BY d.last_seen DESC NULLS LAST
LIMIT 20 OFFSET 0;
```

### 5.5 查询设备的数据定义（用于数据校验和前端渲染）

```sql
SELECT d.dev_id, d.properties_cfg,
       tm.properties, tm.services, tm.events
FROM devices d
JOIN thing_models tm ON d.model_key = tm.model_key
WHERE d.dev_id = $1
  AND d.tenant_key = $2;
```

### 5.6 在线状态更新

```sql
-- 设备上线
UPDATE devices
SET online = TRUE, last_seen = now(), updated_at = now()
WHERE dev_id = $1;

-- 设备离线（由 Last Will 或超时检测触发）
UPDATE devices
SET online = FALSE, updated_at = now()
WHERE dev_id = $1;
```

### 5.7 按物模型统计设备

```sql
SELECT tm.model_key, tm.name, tm.category, tm.version,
       COUNT(d.dev_id) AS device_count,
       COUNT(d.dev_id) FILTER (WHERE d.online) AS online_count
FROM thing_models tm
LEFT JOIN devices d ON d.model_key = tm.model_key AND d.status = 'active'
WHERE tm.tenant_key = $1
  AND tm.status = 'active'
GROUP BY tm.model_key, tm.name, tm.category, tm.version
ORDER BY device_count DESC;
```
