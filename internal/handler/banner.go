package handler

import (
	"avito/internal/banner"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetBanner struct {
	Tag     int    `json:"tag_id" binding:"required"`
	Feature int    `json:"feature_id" binding:"required"`
	Uselast bool   `json:"use_last_revision"`
	Token   string `json:"token"`
}

func (h *Handler) BannerForUser(c *gin.Context) {
  var input GetBanner
	var err error
tag , ok :=c.GetQuery("tag_id")
		if !ok {
			ErrorResponse(c,400,"Некорректные данные", tag)
			return
		} 
input.Tag, err =strconv.Atoi(tag)
		if err != nil {
			ErrorResponse(c,400,"Некорректные данные", tag)
			return
		}
feature , ok :=c.GetQuery("feature_id")
		if !ok {
			ErrorResponse(c,400,"Некорректные данные", feature)
			return
		} 
input.Feature, err =strconv.Atoi(feature)
	  if err != nil {
		  ErrorResponse(c,400,"Некорректные данные", feature)
			return
	  }
uselast,ok:=c.GetQuery("use_last_revision")
if ok{
input.Uselast,err =strconv.ParseBool(uselast)
	  if err != nil {
	  	ErrorResponse(c,400,"Некорректные данные", uselast)
			return
	  }
	}
input.Token=c.GetHeader("token")
		if input.Token != "user_token" && input.Token != "admin_token"{
			newErrorResponse(c,401,"Пользователь не авторизован")
			return
		}

output,err:=h.services.Banner.GetBanner(input.Tag,input.Feature,input.Token, input.Uselast)
    if fmt.Sprint(err)=="sql: no rows in result set"{
			newErrorResponse(c,404,"Баннер не найден")
			return
		}
    if err!=nil{
			newErrorResponse(c,500,"Внутренняя ошибка сервера")
			return
		}

		if !output.IsActive && input.Token!="admin_token"{
			newErrorResponse(c,403,"Пользователь не имеет доступа")
			return
		}

	newErrorResponse(c,200,"Баннер пользователя")
	c.JSON(http.StatusOK, map[string]interface{}{
		"title": output.Title,
		"text": output.Text,
		"url": output.Url,
	})

}



type ContentInput struct{
 Title    string `json:"title"`
 Text     string `json:"text"`
 Url      string `json:"url"`
}

type CreateInput struct{
 TagIds    []int      `json:"tag_ids"`
 FeatureId int        `json:"feature_id"`
 Content ContentInput `json:"content"`
 IsActive bool        `json:"is_active"`
}


func(h *Handler) CreateNewBanner (c *gin.Context){
	var input CreateInput
	var token string
 if err:= c.BindJSON(&input); err!=nil{
		ErrorResponse(c,400, "Некорректные данные", err.Error())
		return 
	}
	token=c.GetHeader("token")
		if token != "user_token" && token != "admin_token"{
			newErrorResponse(c,401,"Пользователь не авторизован")
			return
		}
	if token!="admin_token"{
		newErrorResponse(c,403,"Пользователь не имеет доступа")
		return
	}
	var Input banner.Banner=banner.Banner{Token: token,TagIds: input.TagIds,FeatureId: input.FeatureId}
  Input.Content=banner.Content{Title: input.Content.Title, Text: input.Content.Text,Url:input.Content.Url, IsActive: input.IsActive}
	
	bannerID,err:=h.services.CreateBanner(Input)
	if err!= nil{
		ErrorResponse(c,500,"Внутренняя ошибка сервера", err.Error())
			return
	}
	c.JSON(201, map[string]interface{}{
	"banner_id":bannerID,
	})
}

type GetBannerByFilter struct {
	Token     string `json:"token"`
	FeatureId int    `json:"feature_id"`
	TagId     int    `json:"tag_id"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

func(h *Handler) GetBannerByTagOrFeature(c *gin.Context){
	var input GetBannerByFilter
	var err error

	input.Token=c.GetHeader("token")
		if input.Token != "user_token" && input.Token != "admin_token"{
			newErrorResponse(c,401,"Пользователь не авторизован")
			return
		} 

		if input.Token!="admin_token"{
			newErrorResponse(c,403,"Пользователь не имеет доступа")
			return
		}

feature , ok :=c.GetQuery("feature_id")
		if ok {
input.FeatureId, err =strconv.Atoi(feature)
	  if err != nil {
		  newErrorResponse(c,400,"Некорректные данные")
			return
	  }
	}

	tag , ok :=c.GetQuery("tag_id")
		if ok {
input.TagId, err =strconv.Atoi(tag)
		if err != nil {
			newErrorResponse(c,400,"Некорректные данные")
			return
		}
	}

limit , ok :=c.GetQuery("limit")
		if ok {
input.Limit, err =strconv.Atoi(limit)
	  if err != nil {
		  newErrorResponse(c,400,"Некорректные данные")
			return
	  }
	}

offset , ok :=c.GetQuery("offset")
		if ok {
input.Offset, err =strconv.Atoi(offset)
	  if err != nil {
		  newErrorResponse(c,400,"Некорректные данные")
			return
	  }
	}
fmt.Println(input)
ban,err:=h.services.GetAllBannerWithFilter(input.FeatureId,input.TagId,input.Limit,input.Offset)
if err != nil {
	ErrorResponse(c,500,"Внутренняя ошибка сервера", err.Error())
	return
}
c.JSON(201, map[string]interface{}{
	"items":ban,
	}) 
}


type PatchInput struct{
	Id    string `json:"id" binding:"required"`
	Token string `json:"token"`
}

func(h *Handler) UpdateBanner(c *gin.Context){
	var head PatchInput
	var input banner.UpdateBody

	if err:= c.BindJSON(&input); err!=nil{
		newErrorResponse(c,400, "Некорректные данные")
		return 
	}


	 head.Id=c.Param("id")
   id,err:=strconv.Atoi(head.Id)
	 if err!=nil{
		newErrorResponse(c,400,"Некорректные данные")
		return
	 }

	 head.Token=c.GetHeader("token")
	 if head.Token != "user_token" && head.Token != "admin_token"{
		 newErrorResponse(c,401,"Пользователь не авторизован")
		 return
	 } 

	 if head.Token!="admin_token"{
		 newErrorResponse(c,403,"Пользователь не имеет доступа")
		 return
	 }

  err=h.services.UpdateBanner(id,input)
if err!=nil{
	ErrorResponse(c,500,"Внутренняя ошибка сервера", err.Error())
		 return
}

newErrorResponse(c,200,"OK")
}

type DeleteInput struct{
	Id    string `json:"id" binding:"required"`
	Token string `json:"token"`
}

func(h *Handler) DeleteBannerById(c *gin.Context){
 var input DeleteInput
 
 input.Id=c.Param("id")
 id,err:=strconv.Atoi(input.Id)
 if err!=nil{
	ErrorResponse(c,400,"Некорректные данные", input.Id)
	return
 }

 input.Token=c.GetHeader("token")
 if input.Token != "user_token" && input.Token != "admin_token"{
	 newErrorResponse(c,401,"Пользователь не авторизован")
	 return
 } 

 if input.Token!="admin_token"{
	 newErrorResponse(c,403,"Пользователь не имеет доступа")
	 return
 }

err=h.services.DeleteBanner(id)

if err!=nil{
	ErrorResponse(c,500,"Внутренняя ошибка сервера", err.Error())
		 return
}

 newErrorResponse(c,200,"OK")
}

