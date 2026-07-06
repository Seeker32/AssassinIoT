package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MqttUser holds the schema definition for the mqtt_user table.
type MqttUser struct {
	ent.Schema
}

// Annotations of the MqttUser.
func (MqttUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "mqtt_user",
		},
	}
}

// Mixin of the MqttUser.
func (MqttUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SoftDeleteMixin{},
	}
}

// Fields of the MqttUser.
func (MqttUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique().
			NotEmpty().
			MaxLen(128).
			Comment("MQTT 客户端用户名"),
		field.String("password_hash").
			NotEmpty().
			MaxLen(256).
			Comment("密码哈希值 (sha256(password + salt))"),
		field.String("salt").
			NotEmpty().
			MaxLen(64).
			Comment("密码盐值，每个用户随机生成"),
		field.Bool("is_superuser").
			Default(false).
			Comment("是否为超级用户"),
		field.Time("created").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
	}
}

// Edges of the MqttUser.
func (MqttUser) Edges() []ent.Edge {
	return nil
}

// Indexes of the MqttUser.
func (MqttUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
	}
}
