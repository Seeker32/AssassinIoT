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

type Account struct {
	ent.Schema
}

func (Account) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "accounts",
		},
	}
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SoftDeleteMixin{},
	}
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Comment("用户名称"),
		field.String("password").Comment("用户密码"),
		field.String("email").Comment("用户邮箱"),
		field.String("avatar_url").Comment("头像链接"),
		field.Int("tenant_id").Comment("关联租户 ID"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("accounts").
			Required(),
	}
}

func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("email"),
		index.Fields("username"),
		index.Fields("tenant_id", "email").Unique(),
	}
}
