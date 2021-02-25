package handlers

import (
	"purchases/pkg/buisness"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	buis *buisness.Buisness
}

func InitHandlersLayer (bl *buisness.Buisness) *Handler {
	return &Handler{buis: bl}
}
func (h *Handler) InitRouting() *gin.Engine{
	router := gin.New()

	auth := router.Group("")
	{
		auth.POST("/signin", h.SingIn)
		auth.POST("/signup", h.SignUp)
		auth.POST("/confirm", h.ConfirmUser)
	}
	api := router.Group("/products", h.identifyUser)
	{

		api.GET("/", h.GetItems)
		api.POST("/", h.CreateItem)
		api.GET("/:id", h.GetItem)
		api.PUT("/:id", h.UpdateItem)
		api.DELETE("/:id", h.DeleteItem)
	}

	return router
}