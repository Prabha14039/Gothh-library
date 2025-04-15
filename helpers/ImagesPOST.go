package helpers

import (
	"database/sql"
)

func Images_insert(db *sql.DB, name string, url string) error {
	query := "INSERT INTO public.images (name, url) VALUES ($1, $2)"
	_, err := db.Exec(query, name, url)
	return err
}
