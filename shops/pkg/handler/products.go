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


// @Summary Get Products
// @Description get products list
// @ID get-products
// @Produce  json
// @Success 200 {object} getAllProductsResponse
// @Failure default {object} Error
// @Router /products [get]
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

// @Summary Get Shops
// @Description get shops list
// @ID get-shops
// @Produce  json
// @Success 200 {object} getAllShopsResponse
// @Failure 400,500 {object} Error
// @Failure default {object} Error
// @Router /shops [get]
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
