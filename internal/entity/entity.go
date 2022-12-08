package entity

import "time"

type ImageInfo struct {
	FileName    string    `db:"filename"`
	CreatedDate time.Time `db:"created_date"`
	UpdatedDate time.Time `db:"updated_date"`
}
