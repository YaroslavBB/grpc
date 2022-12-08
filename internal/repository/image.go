package repository

import (
	"database/sql"
	"grpc/app/internal/entity"

	"github.com/jmoiron/sqlx"
)

type ImageRepository struct {
	db *sqlx.DB
}

func NewImageRepository(db *sqlx.DB) ImageRepository {
	return ImageRepository{
		db: db,
	}
}

func (r *ImageRepository) SaveImage(filename string) error {
	sqlQyery := `
	insert into images (filename, created_date, updated_date)
	values ($1, curret_datetime, curret_datetime)
	`
	_, err := r.db.Exec(sqlQyery, filename)
	return err
}

func (r *ImageRepository) LoadImageList() ([]entity.ImageInfo, error) {
	var data []entity.ImageInfo
	sqlQuery := `
	select filename, created_date, updated_date
	from images
	`
	err := r.db.Select(&data, sqlQuery)
	if len(data) == 0 {
		err = sql.ErrNoRows
	}

	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, entity.ErrNoData
	default:
		return nil, err
	}
}
