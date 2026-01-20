package vtype

type Type struct {
	ID   int    `db:"type_id"`
	Name string `db:"type_name"`
}
