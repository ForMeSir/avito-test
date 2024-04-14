package cache

import (
	"avito/internal/banner"
	"fmt"
	"sync"
	"time"
)

type Cache interface{
	GetBannerFromCache(tag int,feature int) (banner.Content,error)
	AddBannertoCache(tag int,feature int, cont banner.Content)
}

type Cach struct{
	Cache
}

type Simple struct {
	Data map[banner.BannerCache]banner.ContentCache
	ResetTime int64
	m sync.RWMutex
}

func NewSimple() *Simple{
	data:=make(map[banner.BannerCache]banner.ContentCache)
	//der:=banner.BannerCache{TagId: 123, FeatureId: 324}
	// ger:= banner.ContentCache{Time:time.Now().Unix(),Title:"fge", Text:"dhd" ,Url:"fgw" }
	// data[der]=ger
	reset:=time.Now().Unix()+300
	return &Simple{Data:data, ResetTime:reset}
}
func NewCache(simp *Simple) *Cach{
	return &Cach{
		Cache: NewSimple(),
	}
}
func(s *Simple) GetBannerFromCache(tag int,feature int) (ban banner.Content, err error){
	var mapKey banner.BannerCache
	var mapValue banner.ContentCache
	mapKey   = banner.BannerCache{TagId: tag,FeatureId:feature}
	s.m.RLock()
  mapValue,ok:=s.Data[mapKey]
	s.m.RUnlock()
	if !ok || time.Now().Unix()-mapValue.Time>300{
		return ban, fmt.Errorf("")
	}
if s.ResetTime<=time.Now().Unix(){
	s.m.Lock()
	s.ResetTime=time.Now().Unix()+300
	for key,value :=range s.Data{
		if time.Now().Unix()-value.Time>300{
			delete(s.Data,key)
		}
	}
	s.m.Unlock()
}
ban = banner.Content{Title: mapValue.Title,Text: mapValue.Text, Url: mapValue.Url, IsActive: mapValue.IsActive}
	return ban,err
}

func(s *Simple) AddBannertoCache(tag int,feature int, cont banner.Content){
	var mapKey banner.BannerCache
	var mapValue banner.ContentCache
	mapKey   = banner.BannerCache{TagId: tag,FeatureId:feature}
	mapValue = banner.ContentCache{Time: time.Now().Unix(),Title: cont.Title,Text: cont.Text, Url: cont.Url}
	s.m.Lock()
	s.Data[mapKey]=mapValue
	s.m.Unlock()
}