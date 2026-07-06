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

// ThingModel holds the schema definition for the thing_models table.
type ThingModel struct {
	ent.Schema
}

// Annotations of the ThingModel.
func (ThingModel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "thing_models",
		},
	}
}

// Mixin of the ThingModel.
func (ThingModel) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SoftDeleteMixin{},
	}
}

// Fields of the ThingModel.
func (ThingModel) Fields() []ent.Field {
	return []ent.Field{
		field.String("model_key").
			Unique().
			NotEmpty().
			MaxLen(64).
			Comment("模型标识，如 temp_sensor_v1、smart_lock_v2"),
		field.String("tenant_key").
			NotEmpty().
			MaxLen(64).
			Comment("所属租户"),
		field.String("name").
			NotEmpty().
			MaxLen(128).
			Comment("物模型名称"),
		field.Text("description").
			Default("").
			Comment("描述信息"),
		field.String("category").
			NotEmpty().
			MaxLen(64).
			Comment("设备品类，引用 model_categories.category_key"),
		field.JSON("properties", map[string]any{}).
			Default(map[string]any{}).
			Comment("属性定义（JSON Schema），描述设备上报的遥测数据点"),
		field.JSON("services", map[string]any{}).
			Default(map[string]any{}).
			Comment("服务定义，描述可下发给设备的指令及输入输出参数"),
		field.JSON("events", map[string]any{}).
			Default(map[string]any{}).
			Comment("事件定义，描述设备可上报的事件类型及负载结构"),
		field.String("version").
			Default("1.0").
			MaxLen(16).
			Comment("模型版本号，支持平滑升级"),
		field.Int("tenant_id").
			Comment("租户 ID 外键"),
		field.Int("model_category_id").
			Comment("物模型分类 ID 外键"),
		field.Enum("status").
			Values("active", "deprecated", "disabled").
			Default("active").
			Comment("状态：active=启用, deprecated=已废弃但已有设备仍可用, disabled=禁用"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the ThingModel.
func (ThingModel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("thing_models").
			Field("tenant_id").
			Unique().
			Required(),
		edge.From("model_category", ModelCategory.Type).
			Ref("thing_models").
			Field("model_category_id").
			Unique().
			Required(),
		edge.To("devices", Device.Type),
	}
}

// Indexes of the ThingModel.
func (ThingModel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_key"),
		index.Fields("tenant_key", "category"),
	}
}
