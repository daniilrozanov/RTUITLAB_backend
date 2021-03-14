package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SignInData struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

// @Summary SignIn
// @Tags auth
// @Description get auth token
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body SignInData true "Sign In Data"
// @Success 200 {string} string "token"
// @Failure default {object} Error
// @Router /signin [post]
func (h *Handler) SignIn(c *gin.Context){
	var input SignInData

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.serv.ConfirmUser(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := h.serv.GenerateToken(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

