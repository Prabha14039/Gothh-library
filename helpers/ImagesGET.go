package helpers

import (
	"database/sql"
)

type Images struct {
	Id int
	Name string
	Url string
}

func Images_fetch(db *sql.DB)([]Images, error) {

	rows , err := db.Query("SELECT id,name,url FROM public.images ORDER BY id DESC")
	if err != nil {
		return nil,err;
	}

 	defer rows.Close()
	var images []Images

	for rows.Next() {
		var image Images;
		if err := rows.Scan(&image.Id, &image.Name ,&image.Url); err !=nil {
			return nil, err;
		}

		images = append(images,image)
	}

	return images , err
}
