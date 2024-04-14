package service

import (
	"avito/internal/banner"
	"avito/internal/cache"
	"avito/internal/repository"
)

type Banner interface {
	CreateBanner(banner banner.Banner) (bannerId int, err error)
	GetBanner(tagId int, featureId int, token string, uselast bool) (banner.Content, error)
	GetAllBannerWithFilter(feature int, tag int, limit int, offset int) (ban []banner.FullBanner, err error) //Получение всех баннеров c фильтрацией по фиче и/или тегу 
	UpdateBanner(id int,update banner.UpdateBody) (err error)
	DeleteBanner(bannerId int) error
}

type Service struct {
	Banner
}

func NewService(repos *repository.Repository, cache *cache.Cach) *Service{
	return &Service{
		Banner: NewBannerService(repos.TodoBanner, cache.Cache),
	}
}