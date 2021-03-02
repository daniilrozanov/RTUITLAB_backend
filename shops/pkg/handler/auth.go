package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SingInData struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

func (h *Handler) SignIn(c *gin.Context){
	var input SingInData

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


