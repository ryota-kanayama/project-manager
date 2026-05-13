package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Milestone マイルストーン
type Milestone struct {
	ent.Schema
}

func (Milestone) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(newUUIDv7),
		field.String("name").
			NotEmpty().
			MaxLen(255),
		field.Time("due_date").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "date"}),
		field.Enum("status").
			Values("pending", "completed").
			Default("pending"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Milestone) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("milestones").
			Unique().
			Required(),
		edge.To("tasks", Task.Type),
	}
}
