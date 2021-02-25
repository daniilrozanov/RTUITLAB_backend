package handlers

import (
	templates "purchases/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type getAllProductsResponse struct {
	Data []templates.Product `json:"data"`
}

type getProductResponse struct {
	Data templates.Product `json:"data"`
}

func (h *Handler) CreateItem (c *gin.Context){
	userId, err := h.getUserId(c)
	if  err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	var input templates.Product
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.buis.ProductLogging.CreateProduct(userId, &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetItem (c *gin.Context){
	userId, err := h.getUserId(c)
	if  err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	prod, err := h.buis.GetProductById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getProductResponse{
		Data: prod,
	})
}

func (h *Handler) GetItems (c *gin.Context){
	userId, err := h.getUserId(c)
	if  err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	prods, err := h.buis.ProductLogging.GetAllProducts(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllProductsResponse{
		Data: prods,
	})
}

func (h *Handler) UpdateItem (c *gin.Context){
	userId, err := h.getUserId(c)
	if  err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	var input templates.UpdateProductInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	err = h.buis.ProductLogging.UpdateProduct(userId, id, &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) DeleteItem (c *gin.Context){
	userId, err := h.getUserId(c)
	if  err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.buis.DeleteProduct(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}