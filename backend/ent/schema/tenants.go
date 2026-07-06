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

// Tenant holds the schema definition for the tenants table.
type Tenant struct {
	ent.Schema
}

// Annotations of the Tenant.
func (Tenant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "tenants",
		},
	}
}

// Mixin of the Tenant.
func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SoftDeleteMixin{},
	}
}

// Fields of the Tenant.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_key").
			Unique().
			NotEmpty().
			MaxLen(64).
			Comment("租户业务标识，用于 API、MQTT 认证等场景"),
		field.String("name").
			NotEmpty().
			MaxLen(128).
			Comment("租户名称，用于前端展示"),
		field.Text("description").
			Default("").
			Comment("描述信息"),
		field.Enum("status").
			Values("active", "disabled").
			Default("active").
			Comment("租户状态：active=正常, disabled=已禁用"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Tenant.
func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("accounts", Account.Type),
		edge.To("model_categories", ModelCategory.Type),
		edge.To("thing_models", ThingModel.Type),
		edge.To("devices", Device.Type),
	}
}

// Indexes of the Tenant.
func (Tenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_key").Unique(),
	}
}
