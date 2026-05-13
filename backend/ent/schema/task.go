package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Task WBS タスク
type Task struct {
	ent.Schema
}

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(newUUIDv7),
		field.String("wbs_code").
			Optional().
			Nillable().
			MaxLen(50),
		field.String("name").
			NotEmpty().
			MaxLen(255),
		field.Text("description").
			Optional().
			Nillable(),
		field.Enum("status").
			Values("not_started", "in_progress", "completed", "blocked").
			Default("not_started"),
		field.Enum("priority").
			Values("low", "medium", "high").
			Default("medium"),
		field.String("assignee").
			Optional().
			Nillable().
			MaxLen(255),
		field.Float("estimated_hours").
			Optional().
			Nillable(),
		field.Float("actual_hours").
			Optional().
			Nillable(),
		field.Time("start_date").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "date"}),
		field.Time("end_date").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "date"}),
		field.Int("sort_order").
			Default(0),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("tasks").
			Unique().
			Required(),
		edge.From("milestone", Milestone.Type).
			Ref("tasks").
			Unique(),
		edge.To("children", Task.Type).
			From("parent").
			Unique(),
		edge.To("worklogs", Worklog.Type),
	}
}
