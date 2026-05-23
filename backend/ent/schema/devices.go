package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Device holds the schema definition for the devices table.
type Device struct {
	ent.Schema
}

// Annotations of the Device.
func (Device) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "devices",
		},
	}
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			StorageKey("dev_id").
			NotEmpty().
			Immutable().
			Comment("设备唯一标识，与 MQTT 客户端用户名一致"),
		field.String("tenant_key").
			NotEmpty().
			MaxLen(64).
			Comment("所属租户"),
		field.String("model_key").
			NotEmpty().
			MaxLen(64).
			Comment("绑定的物模型标识，决定设备的数据结构和能力"),
		field.String("device_name").
			Default("").
			MaxLen(128).
			Comment("设备名称，用于前端展示"),
		field.String("access_key").
			Unique().
			NotEmpty().
			MaxLen(128).
			Comment("设备接入密钥，用于 MQTT 认证"),
		field.String("firmware_ver").
			Default("").
			MaxLen(32).
			Comment("当前固件版本"),
		field.JSON("properties_cfg", map[string]any{}).
			Default(map[string]any{}).
			Comment("属性个性化配置，可覆盖物模型中属性的默认值"),
		field.Int("tenant_id").
			Comment("租户 ID 外键"),
		field.Int("thing_model_id").
			Comment("物模型 ID 外键"),
		field.Enum("status").
			Values("active", "inactive", "disabled").
			Default("active").
			Comment("设备状态：active=正常, inactive=未激活, disabled=已禁用"),
		field.Bool("online").
			Default(false).
			Comment("是否在线（运行时更新）"),
		field.JSON("metadata", map[string]any{}).
			Default(map[string]any{}).
			Comment("扩展元数据，如安装位置、厂商信息等"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("last_seen").
			Optional().
			Nillable().
			Comment("最后在线时间"),
	}
}

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("devices").
			Field("tenant_id").
			Unique().
			Required(),
		edge.From("thing_model", ThingModel.Type).
			Ref("devices").
			Field("thing_model_id").
			Unique().
			Required(),
	}
}

// Indexes of the Device.
func (Device) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_key"),
		index.Fields("tenant_key", "model_key"),
		index.Fields("status"),
		index.Fields("tenant_key", "status"),
	}
}
