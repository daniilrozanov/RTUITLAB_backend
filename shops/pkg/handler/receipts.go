package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shops/pkg"
)

type Recstruct struct {
	ShopId      int `json:"shop_id"`
	PayOptionId int `json:"payoption"`
}

type AddToCartInput struct {
	Category string `json:"category"`
	Quantity int `json:"quantity"`
	ProductId int `json:"product_id"`
	ShopId int `json:"shop_id"`
}

// @Summary AddToCart
// @Tags logged
// @Security ApiKeyAuth
// @Description add product to cart
// @ID add-to-cart
// @Accept json
// @Produce json
// @Param input body AddToCartInput true "Cart Item JSON"
// @Success 200 {object} map[string]string "response"
// @Failure default {object} Error
// @Router /products [post]
func (h *Handler) AddToCart(c *gin.Context) {
	var cartItem pkg.CartItemJSON
	var atci AddToCartInput

	if err := c.BindJSON(&atci); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	cartItem.Quantity = atci.Quantity
	cartItem.ShopID = atci.ShopId
	cartItem.ProductID = atci.ProductId
	cartItem.Category = atci.Category
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

// @Summary GetCarts
// @Tags logged
// @Security ApiKeyAuth
// @Description get carts
// @ID get-carts
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]pkg.CartJSON "response"
// @Failure default {object} Error
// @Router /carts [get]
func (h *Handler) GetCarts(c *gin.Context) {
	var carts *[]pkg.CartJSON

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
		"data": *carts,
	})
}

// @Summary DeleteFromCart
// @Tags logged
// @Security ApiKeyAuth
// @Description delete from cart
// @ID delete-carts
// @Accept json
// @Produce json
// @Param input body pkg.CartItemsOnDeleteJSON true "Cart Items On Delete JSON"
// @Success 200 {object} map[string]string "response"
// @Failure default {object} Error
// @Router /carts [delete]
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

// @Summary CreateReceipt
// @Tags logged
// @Security ApiKeyAuth
// @Description create receipt
// @ID create-receipt
// @Accept json
// @Produce json
// @Param input body Recstruct true "create receipt data"
// @Success 200 {object} map[string]int "response"
// @Failure default {object} Error
// @Router /carts [post]
func (h *Handler) CreateReceipt(c *gin.Context) {
	var rec Recstruct
	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if err := c.BindJSON(&rec); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	recId, err := h.serv.CreateReceipt(rec.ShopId, userId, rec.PayOptionId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.serv.SendReceiptToRabbit(recId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recIds := []int{recId}
	if err := h.serv.SetReceiptsSynchro(&recIds); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": recId,
	})
	//if err := h.serv.TrySynchroByUserId(userId); err != nil {
	//	newErrorResponse(c, http.StatusCreated, err.Error())
	//	return
	//}

}

// @Summary GetReceipts
// @Tags logged
// @Security ApiKeyAuth
// @Description get receipts
// @ID get-receipts
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]pkg.ReceiptJSON "response"
// @Failure default {object} Error
// @Router /receipts [get]
func (h *Handler) GetReceipts(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	recs, err := h.serv.GetReceipts(userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": *recs,
	})
}
