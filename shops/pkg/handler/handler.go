package handler

import (
	"shops/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	serv *service.Service
}

func InitNewHandler(serv *service.Service) *Handler {
	return &Handler{serv: serv}
}

func (h *Handler) InitRoutes () *gin.Engine{
	router := gin.New()
	//router.POST("/", h.SignIn)
	router.POST("/signin", h.SignIn)
	router.GET("/shop", h.GetShops)
	router.GET("/shop/:id", h.GetShop)
	api := router.Group("", h.IdentifyUser)
	{
		api.POST("/products", h.AddToCart)
		api.GET("/cart", h.GetCart)
		api.POST("/cart", h.CreateReceip)
	}
	return router
}



