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
	router.POST("/signin", h.SignIn)
	router.GET("/shops", h.GetShops)
	router.GET("/products", h.GetProducts)
	router.POST("/receive", h.CreateProduct)
	api := router.Group("", h.IdentifyUser)
	{
		api.POST("/products", h.AddToCart)
		api.GET("/carts", h.GetCarts)
		api.DELETE("/carts", h.DeleteFromCart)
		api.GET("/receips", h.GetReceipts)
	}
	return router
}



