package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"purchases/pkg/buisness"
)

type Handler struct {
	buis *buisness.Buisness
}

func InitHandlersLayer(bl *buisness.Buisness) *Handler {
	return &Handler{buis: bl}
}
func (h *Handler) InitRouting() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("")
	{
		auth.POST("/signin", h.SingIn)
		auth.POST("/signup", h.SignUp)
		auth.POST("/confirm", h.ConfirmUser)
	}
	api := router.Group("", h.identifyUser)
	{
		api.GET("/cheques", h.GetReceipts)
		products := api.Group("/products")
		{
			products.GET("/", h.GetItems)
			products.POST("/", h.CreateItem)
			products.GET("/:id", h.GetItem)
			products.PUT("/:id", h.UpdateItem)
			products.DELETE("/:id", h.DeleteItem)
		}
	}
	return router
}
