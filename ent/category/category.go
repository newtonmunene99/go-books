// Code generated by entc, DO NOT EDIT.

package category

import (
	"time"
)

const (
	// Label holds the string label denoting the category type in the database.
	Label = "category"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeBooks holds the string denoting the books edge name in mutations.
	EdgeBooks = "books"
	// Table holds the table name of the category in the database.
	Table = "categories"
	// BooksTable is the table that holds the books relation/edge. The primary key declared below.
	BooksTable = "category_books"
	// BooksInverseTable is the table name for the Book entity.
	// It exists in this package in order to avoid circular dependency with the "book" package.
	BooksInverseTable = "books"
)

// Columns holds all SQL columns for category fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldCreatedAt,
}

var (
	// BooksPrimaryKey and BooksColumn2 are the table columns denoting the
	// primary key for the books relation (M2M).
	BooksPrimaryKey = []string{"category_id", "book_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
)
