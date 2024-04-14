package handler

import (
	"avito/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router:= gin.New()
  

  router.GET("/user_banner",h.BannerForUser)

banner:=router.Group("/banner")
{
   banner.POST("",h.CreateNewBanner)
	 banner.GET("",h.GetBannerByTagOrFeature)
	 banner.PATCH("/:id",h.UpdateBanner)
	  banner.DELETE("/:id",h.DeleteBannerById)
	}
	return router
}