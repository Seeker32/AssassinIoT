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

// ModelCategory holds the schema definition for the model_categories table.
type ModelCategory struct {
	ent.Schema
}

// Annotations of the ModelCategory.
func (ModelCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "model_categories",
		},
	}
}

// Mixin of the ModelCategory.
func (ModelCategory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SoftDeleteMixin{},
	}
}

// Fields of the ModelCategory.
func (ModelCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_key").
			NotEmpty().
			MaxLen(64).
			Comment("所属租户"),
		field.String("category_key").
			NotEmpty().
			MaxLen(64).
			Comment("分类标识，如 temp_humidity_sensor、smart_lock"),
		field.String("display_name").
			NotEmpty().
			MaxLen(128).
			Comment("分类展示名称，如'温湿度传感器'、'智能门锁'"),
		field.Text("description").
			Default("").
			Comment("描述信息"),
		field.String("icon").
			Default("").
			MaxLen(64).
			Comment("图标标识"),
		field.Int("sort_order").
			Default(0).
			Comment("排序值，数字越小越靠前"),
		field.Int("tenant_id").
			Comment("租户 ID 外键"),
		field.Enum("status").
			Values("active", "disabled").
			Default("active").
			Comment("状态：active=启用, disabled=禁用"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the ModelCategory.
func (ModelCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("model_categories").
			Field("tenant_id").
			Unique().
			Required(),
		edge.To("thing_models", ThingModel.Type),
	}
}

// Indexes of the ModelCategory.
func (ModelCategory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_key", "category_key").Unique(),
		index.Fields("tenant_key"),
		index.Fields("tenant_key", "sort_order"),
	}
}
