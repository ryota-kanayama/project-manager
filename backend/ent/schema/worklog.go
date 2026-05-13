package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Worklog 工数記録
type Worklog struct {
	ent.Schema
}

func (Worklog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(newUUIDv7),
		field.String("user_name").
			NotEmpty().
			MaxLen(255),
		field.Float("hours"),
		field.Time("work_date").
			SchemaType(map[string]string{dialect.Postgres: "date"}),
		field.Text("description").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

func (Worklog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", Task.Type).
			Ref("worklogs").
			Unique().
			Required(),
	}
}
