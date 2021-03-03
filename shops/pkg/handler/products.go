package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shops/pkg"
)

type CreateProductData struct {
	Prod pkg.Product `json:"product" binding:"required"`
	ShopsCount []pkg.ShopsProducts `json:"map" binding:"required"`
}

type getAllProductsResponse struct {
	Data []pkg.Product `json:"data"`
}

type getAllShopsResponse struct {
	Data []pkg.Shop `json:"data"`
}

func (h *Handler) CreateProduct(c *gin.Context){
	var data CreateProductData

	if err := c.BindJSON(&data); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Printf("%d %d\n", data.ShopsCount[0].ShopID, data.ShopsCount[0].Quantity)
	id, err := h.serv.ReceiveProduct(data.Prod, data.ShopsCount)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})
}

func (h *Handler) GetProducts(c *gin.Context){
	prods, err := h.serv.GetAllProducts()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllProductsResponse{
		Data: prods,
	})
}

func (h *Handler) GetShops(c *gin.Context){
	shops, err := h.serv.GetAllShops()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllShopsResponse{
		Data: shops,
	})
}


