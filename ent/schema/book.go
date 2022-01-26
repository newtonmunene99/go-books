package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{

		field.String("title").NotEmpty(),

		field.String("author").NotEmpty(),

		field.Int("category_id"),

		field.Int("year").Positive(),

		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).Ref("books"),
	}
}

// func (Book) Indexes() []ent.Index {
// 	return []ent.Index{
// 		index.Fields("title", "author").Unique(),
// 	}
// }
