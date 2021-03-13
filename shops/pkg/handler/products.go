package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shops/pkg"
)


type getAllProductsResponse struct {
	Data []pkg.Product `json:"data"`
}

type getAllShopsResponse struct {
	Data []pkg.Shop `json:"data"`
}

func (h *Handler) GetProducts(c *gin.Context) {
	prods, err := h.serv.GetAllProducts()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllProductsResponse{
		Data: prods,
	})
}

func (h *Handler) GetShops(c *gin.Context) {
	shops, err := h.serv.GetAllShops()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllShopsResponse{
		Data: shops,
	})
}
