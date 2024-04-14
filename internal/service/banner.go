package service

import (
	"avito/internal/banner"
	"avito/internal/cache"
	"avito/internal/repository"
)

type BannerService struct {
	repo repository.TodoBanner
	cach cache.Cache
}

func NewBannerService(repo repository.TodoBanner, cach cache.Cache) *BannerService{
	return &BannerService{repo : repo, cach:cach}
}

func (b *BannerService) GetBanner(tagId int, featureId int, token string, uselast bool)  (ban banner.Content, err error){
	if !uselast || token=="user_token"{
		ban,err=b.cach.GetBannerFromCache(tagId,featureId)
		if err!=nil{
			ban, err =b.repo.FindOne(tagId, featureId)
			b.cach.AddBannertoCache(tagId,featureId,ban)
			return
		}
	return
	}
 ban, err =b.repo.FindOne(tagId, featureId)
 b.cach.AddBannertoCache(tagId,featureId,ban)
 return
}

func(b *BannerService) CreateBanner(banner banner.Banner) (bannerId int, err error){
	return b.repo.Create(banner)
}

func(b *BannerService) GetAllBannerWithFilter(feature int, tag int, limit int, offset int) (ban []banner.FullBanner, err error){
	return b.repo.FindAllByFilter(feature,tag,limit,offset)
}

func(b *BannerService) UpdateBanner(id int,update banner.UpdateBody) (err error){
	return b.repo.Update(id ,update)
}

func(b *BannerService) DeleteBanner(bannerId int) error{
	return b.repo.Delete(bannerId)
}