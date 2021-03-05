package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shops/pkg"
)

func (h *Handler) AddToCart(c *gin.Context) {
	var cartItem pkg.CartItem

	if err := c.BindJSON(&cartItem); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.serv.AddToCart(userId, &cartItem)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func (h *Handler) GetCarts(c *gin.Context) {
	var carts []pkg.CartJSON

	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	carts, err = h.serv.GetCarts(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": carts,
	})
}

func (h *Handler) DeleteFromCart(c *gin.Context) {
	var item pkg.CartItemsOnDeleteJSON

	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if err := c.BindJSON(&item); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serv.DeleteFromCart(&item, userId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}

func (h *Handler) CreateReceipt(c *gin.Context) {
	var rec struct {
		ShopId int `json:"shop_id"`
	}
	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if _, err := h.serv.CreateReceipt(rec.ShopId, userId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.serv.TrySynchroByUserId(userId); err != nil {
		newErrorResponse(c, http.StatusCreated, err.Error())
		return
	}
}

func (h *Handler) GetReceipts(c *gin.Context) {

}
