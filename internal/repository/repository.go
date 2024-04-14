package repository

import (
	"avito/internal/banner"

	"github.com/jmoiron/sqlx"
)

type TodoBanner interface {
	Create(banner banner.Banner)(id int,err error)
	FindOne(tag int, feature int) (banner banner.Content, err error)
	FindAllByFilter(feature int, tag int, limit int, offset int) (ban []banner.FullBanner, err error)
	Update(id int,update banner.UpdateBody)(err error)
	Delete(id int) error
}

type Repository struct {
	TodoBanner
}

func NewRepository(db *sqlx.DB) *Repository{
	return &Repository{
		TodoBanner: NewBannerPostgres(db),
	}
}