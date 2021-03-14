package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary GetReceipts
// @Tags logged
// @Security ApiKeyAuth
// @Description get receipts
// @ID get-receipts
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]templates.ReceiptJSON "response"
// @Failure default {object} Error
// @Router /cheques [get]
func (h *Handler) GetReceipts(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	recs, err := h.buis.GetReceipts(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": recs,
	})
}
