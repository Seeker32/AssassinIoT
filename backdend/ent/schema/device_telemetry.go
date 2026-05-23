package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// DeviceTelemetry holds one device property value at one point in time.
type DeviceTelemetry struct {
	ent.View
}

// Annotations of the DeviceTelemetry.
func (DeviceTelemetry) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "device_telemetry",
		},
	}
}

// Fields of the DeviceTelemetry.
func (DeviceTelemetry) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time"),
		field.String("device_identifier"),
		field.String("property_name"),
		field.Enum("value_type").
			Values("int", "float", "string", "bool", "json", "null"),
		field.Int32("value_int").
			Optional().
			Nillable(),
		field.Float("value_float").
			Optional().
			Nillable(),
		field.String("value_string").
			Optional().
			Nillable(),
		field.Bool("value_bool").
			Optional().
			Nillable(),
		field.JSON("value_json", map[string]any{}).
			Optional(),
		field.Time("received_at"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the DeviceTelemetry.
func (DeviceTelemetry) Edges() []ent.Edge {
	return nil
}
