package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Project プロジェクト
type Project struct {
	ent.Schema
}

func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(newUUIDv7),
		field.String("name").
			NotEmpty().
			MaxLen(255),
		field.Text("description").
			Optional().
			Nillable(),
		field.Enum("status").
			Values("planning", "in_progress", "completed", "on_hold").
			Default("planning"),
		field.Time("start_date").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "date"}),
		field.Time("end_date").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "date"}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("milestones", Milestone.Type),
		edge.To("tasks", Task.Type),
	}
}
