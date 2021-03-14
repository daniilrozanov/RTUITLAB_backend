package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"shops/pkg/service"
)

type Handler struct {
	serv *service.Service
}

func InitNewHandler(serv *service.Service) *Handler {
	return &Handler{serv: serv}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/signin", h.SignIn)
	router.GET("/shops", h.GetShops)
	router.GET("/products", h.GetProducts)
	//router.POST("/receive", h.CreateProduct)
	api := router.Group("", h.IdentifyUser)
	{
		api.POST("/products", h.AddToCart)
		api.GET("/carts", h.GetCarts)
		api.DELETE("/carts", h.DeleteFromCart)
		api.POST("/carts", h.CreateReceipt)
		api.GET("/receipts", h.GetReceipts)
	}
	return router
}
