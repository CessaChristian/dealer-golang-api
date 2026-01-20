package brand

type Brand struct {
	ID   int    `db:"brand_id"`
	Name string `db:"brand_name"`
}
